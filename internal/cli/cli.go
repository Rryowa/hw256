package cli

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"homework-1/internal/entities"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Service interface {
	UpdateCache() error
	AcceptOrder(order entities.Order) error
	ReturnOrderToCourier(orderID string) error
	IssueOrders(orderIDs []string) error
	AcceptReturn(orderID, userID string) error
	ListReturns(page, pageSize int) error
	ListOrders(userID string, limit int) error
}

type CLI struct {
	Service
	commandList []command
}

func NewCLI(s Service) CLI {
	return CLI{
		Service: s,
		commandList: []command{
			{
				name:        help,
				description: "Cправка",
			},
			{
				name:        acceptOrder,
				description: "Принять заказ: accept --id=12345 --r_id=54321 --date=2077-06-06",
			},
			{
				name:        returnOrderToCourier,
				description: "Вернуть заказ курьеру: return_courier --id=12345",
			},
			{
				name:        issueOrders,
				description: "Выдать заказ клиенту: issue --ids=1,2,3",
			},
			{
				name:        acceptReturn,
				description: "Принять возврат: accept_return --id=1 --r_id=2",
			},
			{
				name:        listReturns,
				description: "Вывести список возвратов: list_returns --page=1 --size=10",
			},
			{
				name:        listOrders,
				description: "Вывести список заказов отсортированный\n	 	по Сроку хранения: list_orders --u_id=1 --limit=3",
			},
			{
				name:        exit,
				description: "Выход",
			},
		},
	}
}

func (c CLI) Run() error {
	scanner := bufio.NewScanner(os.Stdin)
	c.updateCache()
	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()
		if len(input) == 0 {
			fmt.Println("command isn't set")
			continue
		}

		args := strings.Split(input, " ")
		commandName := args[0]
		switch commandName {
		case help:
			c.help()
		case acceptOrder:
			if err := c.acceptOrder(args[1:]); err != nil {
				log.Println("Error executing command:", err)
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
		case "exit":
			return nil
		default:
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}
}

func (c CLI) updateCache() error {
	return c.Service.UpdateCache()
}

func (c CLI) acceptOrder(args []string) error {
	var id, recipientId, dateStr string
	fs := flag.NewFlagSet(acceptOrder, flag.ContinueOnError)
	fs.StringVar(&id, "id", "0", "use --id=12345")
	fs.StringVar(&recipientId, "r_id", "0", "use --r_id=54321")
	fs.StringVar(&dateStr, "date", "0", "use --date=2024-06-06")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(id) == 0 {
		return errors.New("id is empty")
	}
	if len(recipientId) == 0 {
		return errors.New("recipient id is empty")
	}
	storageUntil, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		log.Println("Error parsing date:", err)
		return err
	}

	return c.Service.AcceptOrder(entities.Order{
		ID:           id,
		RecipientID:  recipientId,
		Issued:       false,
		Returned:     false,
		StorageUntil: storageUntil,
	})
}

func (c CLI) returnOrderToCourier(args []string) error {
	var id string
	fs := flag.NewFlagSet(returnOrderToCourier, flag.ContinueOnError)
	fs.StringVar(&id, "id", "0", "use --id=12345")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(id) == 0 {
		return errors.New("id is empty")
	}

	return c.Service.ReturnOrderToCourier(id)
}

func (c CLI) issueOrders(args []string) error {
	var idString string
	fs := flag.NewFlagSet(issueOrders, flag.ContinueOnError)
	fs.StringVar(&idString, "ids", "", "use --ids=1,2,3")
	if err := fs.Parse(args); err != nil {
		return err
	}
	ids := strings.Split(idString, ",")
	return c.Service.IssueOrders(ids)
}

func (c CLI) acceptReturn(args []string) error {
	var id, recipientId string
	fs := flag.NewFlagSet(acceptReturn, flag.ContinueOnError)
	fs.StringVar(&id, "id", "0", "use --id=12345")
	fs.StringVar(&recipientId, "r_id", "0", "use --r_id=54321")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(id) == 0 {
		return errors.New("id is empty")
	}
	if len(recipientId) == 0 {
		return errors.New("recipient id is empty")
	}

	return c.Service.AcceptReturn(id, recipientId)
}

func (c CLI) listReturns(args []string) error {
	var page, size string
	fs := flag.NewFlagSet(listReturns, flag.ContinueOnError)
	fs.StringVar(&page, "page", "0", "use --page=1")
	fs.StringVar(&size, "size", "0", "use --size=10")

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

	return c.Service.ListReturns(p, ps)
}

func (c CLI) listOrders(args []string) error {
	var userId, limit string
	fs := flag.NewFlagSet(listOrders, flag.ContinueOnError)
	fs.StringVar(&userId, "u_id", "0", "use --u_id=1")
	fs.StringVar(&limit, "limit", "0", "use --limit=3")

	if err := fs.Parse(args); err != nil {
		return err
	}

	l, err := strconv.Atoi(limit)
	if err != nil {
		return err
	}

	return c.Service.ListOrders(userId, l)
}

func (c CLI) help() {
	fmt.Println("command list:")
	for _, cmd := range c.commandList {
		fmt.Println("  ", cmd.name, cmd.description)
	}
	return
}
