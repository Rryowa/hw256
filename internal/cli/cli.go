package cli

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"homework-1/internal/models"
	"homework-1/internal/service"
	"log"
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
	orderService service.OrderService
	commandList  []command

	maxGoroutines    uint64
	activeGoroutines uint64
}

func NewCLI(os service.OrderService) *CLI {
	return &CLI{
		orderService: os,
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
				description: "Список возвратов: list_returns -limit=1 -offset=0",
			},
			{
				name:        listOrders,
				description: "Список заказов: list_orders -u_id=1 -limit=3",
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

func (c *CLI) Run() error {
	semaphore := make(chan struct{}, 1)
	commandChannel := make(chan string)
	done := make(chan struct{})

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	if err := c.setMaxGoroutines(fmt.Sprintf(
		"set_mg -n=%s", strconv.Itoa(runtime.GOMAXPROCS(0))),
		&semaphore); err != nil {
		return err
	}

	var wg sync.WaitGroup

	go signalListener(signalChannel, done)

	//Reader
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			commandChannel <- scanner.Text()
		}
	}()

	//Handler
	go c.commandHandler(commandChannel, semaphore, done, &wg)

	<-done

	wg.Wait()

	//Close where created
	close(semaphore)
	fmt.Println("All goroutines finished. Exiting...")

	return nil
}

func signalListener(signalChannel chan os.Signal, done chan struct{}) {
	for {
		<-signalChannel
		fmt.Println("\nReceived shutdown signal")
		done <- struct{}{}
	}
}

func (c *CLI) commandHandler(commandChannel chan string, semaphore chan struct{}, done chan struct{}, wg *sync.WaitGroup) {
	for {
		cmd := <-commandChannel

		if strings.HasPrefix(cmd, exit) {
			done <- struct{}{}
		} else if strings.HasPrefix(cmd, setMaxGoroutines) {
			if err := c.setMaxGoroutines(cmd, &semaphore); err != nil {
				log.Fatal(err)
			}
		} else {
			wg.Add(1)
			atomic.AddUint64(&c.activeGoroutines, 1)
			id := atomic.LoadUint64(&c.activeGoroutines)

			go c.worker(cmd, id, semaphore, wg)
		}
	}
}

func (c *CLI) worker(cmd string, id uint64, semaphore chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Worker %d: Waiting to acquire semaphore\n", id)
	semaphore <- struct{}{}

	log.Printf("Worker %d: Working\n", id)
	c.processCommand(cmd)

	log.Printf("Worker %d: Semaphore released\n\n", id)
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

	fmt.Printf("Number of goroutines set to %d\n", n)
	return nil
}

func (c *CLI) processCommand(input string) {
	args := strings.Split(input, " ")
	commandName := args[0]

	switch commandName {
	case acceptOrder:
		if err := c.acceptOrder(args[1:]); err != nil {
			log.Println(err)
		} else {
			log.Println("Order accepted.")
		}
	case issueOrders:
		if err := c.issueOrders(args[1:]); err != nil {
			log.Println(err)
		}
	case acceptReturn:
		if err := c.acceptReturn(args[1:]); err != nil {
			log.Println(err)
		} else {
			log.Println("Return accepted.")
		}
	case returnOrderToCourier:
		if err := c.returnOrderToCourier(args[1:]); err != nil {
			log.Println(err)
		} else {
			log.Println("Order returned.")
		}
	case listReturns:
		if err := c.listReturns(args[1:]); err != nil {
			log.Println(err)
		}
	case listOrders:
		if err := c.listOrders(args[1:]); err != nil {
			log.Println(err)
		}
	case help:
		c.help()
	default:
		fmt.Println("Unknown command. Type 'help' for a list of commands.")
	}
}

func (c *CLI) acceptOrder(args []string) error {
	var idStr, userId, dateStr string
	fs := flag.NewFlagSet(acceptOrder, flag.ContinueOnError)
	fs.StringVar(&idStr, "id", "0", "use -id=12345")
	fs.StringVar(&userId, "u_id", "0", "use -u_id=54321")
	fs.StringVar(&dateStr, "date", "0", "use -date=2024-06-06")

	if err := fs.Parse(args); err != nil {
		return err
	}

	return c.orderService.Accept(idStr, userId, dateStr)
}

func (c *CLI) issueOrders(args []string) error {
	var idString string
	fs := flag.NewFlagSet(issueOrders, flag.ContinueOnError)
	fs.StringVar(&idString, "ids", "", "use -ids=1,2,3")
	if err := fs.Parse(args); err != nil {
		return err
	}
	ids := strings.Split(idString, ",")

	return c.orderService.Issue(ids)
}

func (c *CLI) acceptReturn(args []string) error {
	var id, userId string
	fs := flag.NewFlagSet(acceptReturn, flag.ContinueOnError)
	fs.StringVar(&id, "id", "0", "use -id=12345")
	fs.StringVar(&userId, "u_id", "0", "use -u_id=54321")
	if err := fs.Parse(args); err != nil {
		return err
	}

	return c.orderService.Return(id, userId)
}

func (c *CLI) returnOrderToCourier(args []string) error {
	var id string
	fs := flag.NewFlagSet(returnOrderToCourier, flag.ContinueOnError)
	fs.StringVar(&id, "id", "0", "use -id=12345")

	if err := fs.Parse(args); err != nil {
		return err
	}

	return c.orderService.ReturnToCourier(id)
}

func (c *CLI) listReturns(args []string) error {
	var limit, offset string
	fs := flag.NewFlagSet(listReturns, flag.ContinueOnError)
	fs.StringVar(&limit, "limit", "0", "use -limit=1")
	fs.StringVar(&offset, "offset", "0", "use -offset=0")

	if err := fs.Parse(args); err != nil {
		return err
	}

	orderIDs, err := c.orderService.ListReturns(limit, offset)
	if err != nil {
		return err
	}
	printList(orderIDs)
	return nil
}

func (c *CLI) listOrders(args []string) error {
	var userId, limit string
	fs := flag.NewFlagSet(listOrders, flag.ContinueOnError)
	fs.StringVar(&userId, "u_id", "0", "use -u_id=1")
	fs.StringVar(&limit, "limit", "0", "use -limit=3")

	if err := fs.Parse(args); err != nil {
		return err
	}

	orderIDs, err := c.orderService.ListOrders(userId, limit)
	if err != nil {
		return err
	}
	printList(orderIDs)
	return nil
}

func printList(Orders []models.Order) {
	if len(Orders) == 0 {
		defer fmt.Printf("\n\n")
	}
	fmt.Printf("%-20s %-20s %-20s %-10s %-20s %-10s\n", "ID", "userId", "StorageUntil", "Issued", "IssuedAt", "Returned")
	fmt.Println(strings.Repeat("-", 100))
	for _, order := range Orders {
		fmt.Printf("%-20s %-20s %-20s %-10v %-20s %-10v\n",
			order.ID,
			order.UserID,
			order.StorageUntil.Format("2006-01-02 15:04:05"),
			order.Issued,
			order.IssuedAt.Format("2006-01-02 15:04:05"),
			order.Returned)
	}
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
