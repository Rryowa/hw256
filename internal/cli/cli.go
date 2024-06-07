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
	mu               sync.Mutex
	cond             *sync.Cond
}

func NewCLI(v service.ValidationService, o service.OrderService, f service.FileService) *CLI {
	cli := &CLI{
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
				name:        setWorkers,
				description: "Количество работающих гоферов: set_workers -n=1",
			},
			{
				name:        exit,
				description: "Выход",
			},
		},
	}
	cli.cond = sync.NewCond(&cli.mu)
	return cli
}

// TODO: START OF REFACTOR
func (c *CLI) Run() error {
	if err := c.updateCache(); err != nil {
		return err
	}
	scanner := bufio.NewScanner(os.Stdin)

	commandChannel := make(chan string)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	err := c.setWorkers([]string{"-n", strconv.Itoa(runtime.GOMAXPROCS(0))})
	if err != nil {
		return err
	}

	done := make(chan struct{})
	var wg sync.WaitGroup

	go func() {
		for {
			select {
			case <-signalChannel:
				fmt.Println("\nReceived shutdown signal")
				close(done)
				close(commandChannel)
				return
			case cmd, ok := <-commandChannel:
				if !ok {
					return
				}
				for atomic.LoadUint64(&c.activeGoroutines) >= atomic.LoadUint64(&c.maxGoroutines) {
					c.mu.Lock()
					c.cond.Wait()
					c.mu.Unlock()
				}
				atomic.AddUint64(&c.activeGoroutines, 1)

				wg.Add(1)
				go func(cmd string) {
					defer wg.Done()
					c.processCommand(cmd, atomic.LoadUint64(&c.activeGoroutines))
					atomic.AddUint64(&c.activeGoroutines, ^uint64(0))
					c.mu.Lock()
					c.cond.Signal()
					c.mu.Unlock()
				}(cmd)
				//TODO:remove
			case <-done:
				return
			}
		}
	}()

	go func() {
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
	}()

	<-done    // Block until done is closed
	wg.Wait() // Wait for all goroutines to finish
	fmt.Println("All goroutines finished. Exiting...")
	return nil
}

func (c *CLI) setWorkers(args []string) error {
	var ns string
	fs := flag.NewFlagSet(setWorkers, flag.ContinueOnError)
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
	c.mu.Lock()
	c.cond.Broadcast()
	c.mu.Unlock()

	fmt.Printf("Number of goroutines set to %d\n", n)
	return nil
}

func (c *CLI) processCommand(input string, gid uint64) {
	args := strings.Split(input, " ")
	commandName := args[0]

	fmt.Printf("\nGoroutine %d started %s task...\n", gid, commandName)
	defer func(gid uint64) {
		fmt.Printf("Goroutine %d finished %s task!\n", gid, commandName)
	}(gid)

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
	case setWorkers:
		if err := c.setWorkers(args[1:]); err != nil {
			log.Println(err)
		}
	case exit:
		os.Exit(0)
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
	time.Sleep(500 * time.Millisecond)
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
