package view

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"homework/internal/models"
	"homework/internal/service"
	"homework/pkg/kafka"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
)

type CLI struct {
	orderService  service.OrderService
	loggerService kafka.LoggerService
	zapLogger     *zap.SugaredLogger
	commandList   []command

	maxGoroutines    uint64
	activeGoroutines uint64
}

func NewCLI(os service.OrderService, log kafka.LoggerService, zap *zap.SugaredLogger) *CLI {
	return &CLI{
		orderService:  os,
		loggerService: log,
		zapLogger:     zap,
		commandList: []command{
			{
				name:        help,
				description: "Справка",
			},
			{
				name:        acceptOrder,
				description: "Принять заказ: accept -id=12345 -u_id=54321 -date=2077-06-06",
			},
			{
				name:        returnOrderToCourier,
				description: "Вернуть заказ курьеру: return_courier -id=12345",
			},
			{
				name:        issueOrders,
				description: "Выдать заказ клиенту: issue -ids=1,2,3",
			},
			{
				name:        acceptReturn,
				description: "Принять возврат: accept_return -id=1 -u_id=2",
			},
			{
				name:        listReturns,
				description: "Список возвратов: list_returns -lmt=10 -ofs=0",
			},
			{
				name:        listOrders,
				description: "Список заказов: list_orders -u_id=1 -lmt=10 -ofs=0",
			},
			{
				name:        setMaxGoroutines,
				description: "Максимальное кол-во горутин: set_mg -n=1",
			},
			{
				name:        exit,
				description: "Выход",
			},
		},
	}
}

func (c *CLI) Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	commandChannel := make(chan string)
	semaphore := make(chan struct{}, 1)
	var wg sync.WaitGroup

	if err := c.setMaxGoroutines(fmt.Sprintf(
		"set_mg -n=%s", strconv.Itoa(runtime.GOMAXPROCS(0))), &semaphore); err != nil {
		return err
	}

	go c.signalHandler(ctx, cancel)

	loggerClose := c.loggerService.Start(ctx, &wg)
	defer loggerClose()

	//Reader
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			commandChannel <- scanner.Text()
		}
	}()

	go c.commandHandler(ctx, &wg, cancel, commandChannel, semaphore)

	<-ctx.Done()

	wg.Wait()

	//Close where created
	close(semaphore)
	c.zapLogger.Infoln("All goroutines finished. Bye!")

	return nil
}

func (c *CLI) signalHandler(ctx context.Context, cancel context.CancelFunc) {
	<-ctx.Done()
	c.zapLogger.Debugln("Received shutdown signal...")
	cancel()
}

func (c *CLI) commandHandler(ctx context.Context, wg *sync.WaitGroup, cancel context.CancelFunc, commandChannel chan string, semaphore chan struct{}) {
	for {
		cmd := <-commandChannel

		if strings.HasPrefix(cmd, exit) {
			cancel()
		} else if strings.HasPrefix(cmd, setMaxGoroutines) {
			if err := c.setMaxGoroutines(cmd, &semaphore); err != nil {
				c.zapLogger.Fatalf("setMaxGoroutines error: %v", err)
			}
		} else {
			atomic.AddUint64(&c.activeGoroutines, 1)
			id := atomic.LoadUint64(&c.activeGoroutines)

			go c.worker(ctx, wg, semaphore, cmd, id)
		}
	}
}

func (c *CLI) worker(ctx context.Context, wg *sync.WaitGroup, semaphore chan struct{}, cmd string, id uint64) {
	wg.Add(1)
	defer wg.Done()
	c.zapLogger.Infof("Worker %d: Waiting to acquire semaphore\n", id)
	semaphore <- struct{}{}

	c.zapLogger.Infof("Worker %d: Working\n", id)
	c.processCommand(ctx, cmd)

	c.zapLogger.Infof("Worker %d: Semaphore released\n", id)
	<-semaphore
}

func (c *CLI) setMaxGoroutines(input string, semaphore *chan struct{}) error {
	args := strings.Split(input, " ")
	args = args[1:]
	var ns string
	fs := flag.NewFlagSet(setMaxGoroutines, flag.ContinueOnError)
	fs.StringVar(&ns, "n", "0", "use -n=1")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if len(ns) == 0 {
		return errors.New("number of goroutines is required")
	}
	n, err := strconv.Atoi(ns)
	if err != nil {
		return errors.Join(err, errors.New("invalid argument"))
	}
	if n < 1 {
		return errors.New("number of goroutines must be > 0")
	}

	atomic.StoreUint64(&c.maxGoroutines, uint64(n))
	*semaphore = make(chan struct{}, n)

	c.zapLogger.Infof("Number of goroutines set to %d\n", n)
	return nil
}

func (c *CLI) processCommand(ctx context.Context, input string) {
	args := strings.Split(input, " ")
	commandName := args[0]

	event, err := c.loggerService.CreateEvent(ctx, input)
	if err != nil {
		c.zapLogger.Errorf("error - logger Create event: %v\n", err)
		return
	}

	switch commandName {
	case acceptOrder:
		if err = c.acceptOrder(ctx, args[1:]); err == nil {
			c.zapLogger.Infoln("Order accepted")
		}
	case issueOrders:
		if err = c.issueOrders(ctx, args[1:]); err == nil {
			c.zapLogger.Infoln("Orders issued")
		}
	case acceptReturn:
		if err = c.acceptReturn(ctx, args[1:]); err == nil {
			c.zapLogger.Infoln("Return accepted")
		}
	case returnOrderToCourier:
		if err = c.returnOrderToCourier(ctx, args[1:]); err == nil {
			c.zapLogger.Infoln("Order returned.")
		}
	case listReturns:
		err = c.listReturns(ctx, args[1:])
	case listOrders:
		err = c.listOrders(ctx, args[1:])
	case help:
		c.help()
		return
	default:
		c.zapLogger.Warnln("Unknown command. Type 'help' for a list of commands.")
		return
	}

	if err != nil {
		c.zapLogger.Errorln(err)
	} else if err = c.loggerService.ProcessEvent(ctx, event); err != nil {
		c.zapLogger.Errorf("error - logger Process event: %v\n", err)
	}
}

func (c *CLI) acceptOrder(ctx context.Context, args []string) error {
	dto := models.Dto{}
	fs := flag.NewFlagSet(acceptOrder, flag.ContinueOnError)
	fs.StringVar(&dto.ID, "id", "", "use -id=12345")
	fs.StringVar(&dto.UserID, "u_id", "", "use -u_id=54321")
	fs.StringVar(&dto.StorageUntil, "date", "", "use -date=2024-06-06")
	fs.StringVar(&dto.OrderPrice, "price", "", "use -price=999.99")
	fs.StringVar(&dto.Weight, "w", "", "use -w=10.0")
	fs.StringVar(&dto.PackageType, "p", "", "use -p=box")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := ValidateAcceptArgs(dto); err != nil {
		return err
	}

	return c.orderService.Accept(ctx, dto, dto.PackageType)
}

func (c *CLI) issueOrders(ctx context.Context, args []string) error {
	var idsStr string
	fs := flag.NewFlagSet(issueOrders, flag.ContinueOnError)
	fs.StringVar(&idsStr, "ids", "", "use -ids=1,2,3")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := ValidateIssueArgs(idsStr); err != nil {
		return err
	}
	return c.orderService.Issue(ctx, idsStr)
}

func (c *CLI) acceptReturn(ctx context.Context, args []string) error {
	var id, userId string
	fs := flag.NewFlagSet(acceptReturn, flag.ContinueOnError)
	fs.StringVar(&id, "id", "0", "use -id=12345")
	fs.StringVar(&userId, "u_id", "0", "use -u_id=54321")
	if err := fs.Parse(args); err != nil {
		return err
	}

	err := ValidateAcceptReturnArgs(id, userId)
	if err != nil {
		return err
	}

	return c.orderService.Return(ctx, id, userId)
}

func (c *CLI) returnOrderToCourier(ctx context.Context, args []string) error {
	var id string
	fs := flag.NewFlagSet(returnOrderToCourier, flag.ContinueOnError)
	fs.StringVar(&id, "id", "0", "use -id=12345")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := ValidateReturnToCourierArgs(id); err != nil {
		return err
	}

	return c.orderService.ReturnToCourier(ctx, id)
}

func (c *CLI) listReturns(ctx context.Context, args []string) error {
	var offsetStr, limitStr string
	fs := flag.NewFlagSet(listReturns, flag.ContinueOnError)
	fs.StringVar(&offsetStr, "ofs", "0", "use -ofs=0")
	fs.StringVar(&limitStr, "lmt", "0", "use -lmt=10")

	if err := fs.Parse(args); err != nil {
		return err
	}

	err := ValidateListArgs(offsetStr, limitStr)
	if err != nil {
		return err
	}

	orderIDs, err := c.orderService.ListReturns(ctx, offsetStr, limitStr)
	if err != nil {
		return err
	}

	printList(orderIDs)

	return nil
}

func (c *CLI) listOrders(ctx context.Context, args []string) error {
	var userId, offsetStr, limitStr string
	fs := flag.NewFlagSet(listOrders, flag.ContinueOnError)
	fs.StringVar(&userId, "u_id", "0", "use -u_id=1")
	fs.StringVar(&offsetStr, "ofs", "0", "use -ofs=0")
	fs.StringVar(&limitStr, "lmt", "0", "use -lmt=10")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := ValidateListArgs(offsetStr, limitStr); err != nil {
		return err
	}

	orders, err := c.orderService.ListOrders(ctx, userId, offsetStr, limitStr)
	if err != nil {
		return err
	}

	printList(orders)

	return nil
}

func printList(orders []models.Order) {
	if len(orders) == 0 {
		defer fmt.Printf("\n\n")
	}
	fmt.Printf("%-5s%-10s%-15s%-15v%-10v%-13v%-10v%-13s%-13v\n", "id", "user_id", "storage_until", "issued_at", "returned", "order_price", "weight", "package_type", "package_price")
	fmt.Println(strings.Repeat("-", 100))
	for _, order := range orders {
		fmt.Printf("%-5s%-10s%-15s%-15v%-10v%-13v%-10v%-13s%-13v\n",
			order.ID,
			order.UserID,
			order.StorageUntil.Format("2006-01-02"),
			order.IssuedAt.Format("2006-01-02"),
			order.Returned,
			order.OrderPrice,
			order.Weight,
			order.PackageType,
			order.PackagePrice)
	}
	fmt.Printf("\n")
}

func (c *CLI) help() {
	fmt.Println("Command list:")
	fmt.Printf("%-15s | %-30s | %s\n", "Command", "Description", "Example")
	fmt.Println("---------------------------------------------------------------------------------------------------")
	for _, cmd := range c.commandList {
		parts := strings.SplitN(cmd.description, ":", 2)
		description := ""
		example := ""
		if len(parts) > 0 {
			description = strings.TrimSpace(parts[0])
		}
		if len(parts) > 1 {
			example = strings.TrimSpace(parts[1])
		}
		fmt.Printf("%-15s | %-30s | %s\n", cmd.name, description, example)
	}
}