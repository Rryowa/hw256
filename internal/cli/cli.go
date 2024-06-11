package cli

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"homework-1/internal/entities"
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
	"time"
)

type CLI struct {
	validationService service.ValidationService
	orderService      service.OrderService
	fileService       service.FileService
	commandList       []command

	maxGoroutines    uint64
	activeGoroutines uint64
	semaphore        chan struct{}
}

func NewCLI(v service.ValidationService, o service.OrderService, f service.FileService) *CLI {
	return &CLI{
		validationService: v,
		orderService:      o,
		fileService:       f,
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
				description: "Список возвратов: list_returns -page=1 -size=10",
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
	//if err := c.updateCache(); err != nil {
	//	return err
	//}

	c.semaphore = make(chan struct{}, 1)
	commandChannel := make(chan string)
	done := make(chan struct{})
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	if err := c.setMaxGoroutines(fmt.Sprintf("set_mg -n=%s", strconv.Itoa(runtime.GOMAXPROCS(0))), c.semaphore); err != nil {
		return err
	}

	var wg sync.WaitGroup

	//Reader
	scanner := bufio.NewScanner(os.Stdin)
	go reader(scanner, commandChannel, done)

	//Handler
	go c.handler(signalChannel, commandChannel, c.semaphore, done, &wg)

	<-done

	wg.Wait()

	//close where created
	close(c.semaphore)
	fmt.Println("All goroutines finished. Exiting...")

	return nil
}

func (c *CLI) handler(signalChannel chan os.Signal, commandChannel chan string, semaphore chan struct{}, done chan struct{}, wg *sync.WaitGroup) {
	for {
		select {
		case <-signalChannel:
			fmt.Println("\nReceived shutdown signal")
			done <- struct{}{}
			return
		case cmd, ok := <-commandChannel:
			if !ok {
				return
			}
			if strings.HasPrefix(cmd, exit) {
				done <- struct{}{}
			} else if strings.HasPrefix(cmd, setMaxGoroutines) {
				if err := c.setMaxGoroutines(cmd, semaphore); err != nil {
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
}

func (c *CLI) worker(cmd string, id uint64, semaphore chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Worker %d: Waiting to acquire semaphore\n", id)
	semaphore <- struct{}{}
	log.Printf("Worker %d: Working\n", id)
	c.processCommand(cmd)
	<-semaphore
	log.Printf("Worker %d: Semaphore released\n\n", id)
}

func reader(scanner *bufio.Scanner, commandChannel chan string, done chan struct{}) {
	for {
		if scanner.Scan() {
			input := scanner.Text()
			if len(input) > 0 {
				select {
				case commandChannel <- input:
				case <-done:
					return
				}
			}
		}
	}
}

// TODO: add set active goroutines to 0 after set of max
func (c *CLI) setMaxGoroutines(input string, semaphore chan struct{}) error {
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
	semaphore = make(chan struct{}, n)
	fmt.Printf("Number of goroutines set to %d\n", n)
	return nil
}

func (c *CLI) processCommand(input string) {
	args := strings.Split(input, " ")
	commandName := args[0]

	switch commandName {
	case help:
		c.help()
	case acceptOrder:
		if err := c.acceptOrder(args[1:]); err != nil {
			log.Println(err)
		}
	case returnOrderToCourier:
		if err := c.returnOrderToCourier(args[1:]); err != nil {
			log.Println(err)
		}
	case issueOrders:
		if err := c.issueOrders(args[1:]); err != nil {
			log.Println(err)
		}
	case acceptReturn:
		if err := c.acceptReturn(args[1:]); err != nil {
			log.Println(err)
		}
	case listReturns:
		if err := c.listReturns(args[1:]); err != nil {
			log.Println(err)
		}
	case listOrders:
		if err := c.listOrders(args[1:]); err != nil {
			log.Println(err)
		}
	default:
		fmt.Println("Unknown command. Type 'help' for a list of commands.")
	}
}

func (c *CLI) updateCache() error {
	if err := c.fileService.CheckFile(); err != nil {
		return err
	}

	return c.orderService.UpdateCache()
}

func (c *CLI) acceptOrder(args []string) error {
	var id, userId, dateStr string
	fs := flag.NewFlagSet(acceptOrder, flag.ContinueOnError)
	fs.StringVar(&id, "id", "0", "use -id=12345")
	fs.StringVar(&userId, "u_id", "0", "use -u_id=54321")
	fs.StringVar(&dateStr, "date", "0", "use -date=2024-06-06")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := c.validationService.AcceptValidation(id, userId, dateStr); err != nil {
		return err
	}

	orders := c.orderService.AcceptOrder(id, userId, dateStr)

	if err := c.fileService.Write(orders); err != nil {
		return err
	}

	return nil
}

func (c *CLI) returnOrderToCourier(args []string) error {
	var id string
	fs := flag.NewFlagSet(returnOrderToCourier, flag.ContinueOnError)
	fs.StringVar(&id, "id", "0", "use -id=12345")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := c.validationService.ReturnToCourierValidation(id); err != nil {
		return err
	}

	orders := c.orderService.ReturnOrderToCourier(id)

	if err := c.fileService.Write(orders); err != nil {
		return err
	}

	return nil
}

func (c *CLI) issueOrders(args []string) error {
	var idString string
	fs := flag.NewFlagSet(issueOrders, flag.ContinueOnError)
	fs.StringVar(&idString, "ids", "", "use -ids=1,2,3")
	if err := fs.Parse(args); err != nil {
		return err
	}
	ids := strings.Split(idString, ",")

	if err := c.validationService.IssueValidation(ids); err != nil {
		return err
	}

	orders := c.orderService.IssueOrders(ids)

	if err := c.fileService.Write(orders); err != nil {
		return err
	}

	return nil
}

func (c *CLI) acceptReturn(args []string) error {
	var id, userId string
	fs := flag.NewFlagSet(acceptReturn, flag.ContinueOnError)
	fs.StringVar(&id, "id", "0", "use -id=12345")
	fs.StringVar(&userId, "u_id", "0", "use -u_id=54321")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := c.validationService.ReturnValidation(id, userId); err != nil {
		return err
	}

	orders := c.orderService.AcceptReturn(id)

	if err := c.fileService.Write(orders); err != nil {
		return err
	}

	return nil
}

func (c *CLI) listReturns(args []string) error {
	var page, size string
	fs := flag.NewFlagSet(listReturns, flag.ContinueOnError)
	fs.StringVar(&page, "page", "0", "use -page=1")
	fs.StringVar(&size, "size", "0", "use -size=10")

	if err := fs.Parse(args); err != nil {
		return err
	}
	p, err := strconv.Atoi(page)
	if err != nil {
		return err
	}
	ps, err := strconv.Atoi(size)
	if err != nil {
		return err
	}

	printList(c.orderService.ListReturns(p, ps))
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

	l, err := strconv.Atoi(limit)
	if err != nil {
		return err
	}

	printList(c.orderService.ListOrders(userId, l))
	return nil
}

func printList(Orders []entities.Order) {
	if len(Orders) == 0 {
		fmt.Println("There are no Orders or they all issued!")
		return
	}
	//To prettify output
	time.Sleep(250 * time.Millisecond)
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
	fmt.Printf("%-15s | %-25s | %s\n", "Command", "Description", "Example")
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
		fmt.Printf("%-15s | %-25s | %s\n", cmd.name, description, example)
	}
}
