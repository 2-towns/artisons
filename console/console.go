package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"gifthub/conf"
	"gifthub/console/parser"
	"gifthub/console/populate"
	"gifthub/locales"
	"gifthub/notifications/mails"
	"gifthub/notifications/vapid"
	"gifthub/orders"
	"gifthub/users"
	"log"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/text/message"
)

var (
	printer = message.NewPrinter(locales.Console)
)

func main() {
	start := time.Now()

	if len(os.Args) == 1 {
		log.Fatalln("The command is required, here are the possibilities: import")
	}

	command := os.Args[len(os.Args)-1]

	switch command {
	case "import":
		{
			file := flag.String("file", "./web/testdata/data.csv", "The path to the csv file")

			f, err := os.Open(*file)
			if err != nil {
				log.Fatal(err)
			}

			defer f.Close()

			reader := csv.NewReader(f)
			data, err := reader.ReadAll()
			if err != nil {
				log.Fatal(err)
			}

			lines, err := parser.Import(data, conf.DefaultMID)
			if err != nil {
				log.Panicln(err)
			}

			fmt.Printf("Import successful, %d line(s) imported.\n", lines)

		}

	case "populate":
		{
			err := populate.Run()

			if err != nil {
				log.Panic(err)
			}
		}

	case "orderstatus":
		{
			id := flag.String("id", "", "The order id")
			status := flag.String("status", "", "The new order status")

			flag.Parse()

			err := orders.UpdateStatus(*id, *status)
			if err != nil {
				log.Fatalln(err)
			}

			order, err := orders.Find(*id)
			if err != nil {
				log.Fatalln(err)
			}

			user, err := users.Get(order.UID)
			if err != nil {
				log.Fatalln(err)
			}

			p := message.NewPrinter(user.Lang)
			msg := p.Sprintf("mail_magic_link", id, status)
			mails.Send(user.Email, msg)

			for _, value := range user.Devices {
				vapid.Send(value, msg)
			}
		}

	case "orderdetail":
		{
			id := flag.String("id", "", "The order id")

			flag.Parse()

			o, err := orders.Find(*id)
			if err != nil {
				log.Fatalln(err)
			}

			empJSON, err := json.MarshalIndent(o, "", "  ")
			if err != nil {
				log.Fatalf(err.Error())
			}

			log.Printf("%s\n", string(empJSON))
		}

	case "ordernote":
		{
			id := flag.String("id", "", "The order id")
			note := flag.String("note", "", "The note to attach")

			flag.Parse()

			err := orders.AddNote(*id, *note)
			if err != nil {
				log.Fatalln(err)
			}

			log.Println("Note added to the order.")
		}

	case "userlist":
		{
			page := flag.Int64("page", 0, "The page used in pagination")

			u, err := users.List(*page)
			if err != nil {
				log.Fatalln(err)
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"ID", "Email", "Updated at"})

			for _, user := range u {
				t.AppendRow([]interface{}{user.ID, user.Email, user.UpdatedAt})
			}

			t.Render()
		}
	default:
		{
			log.Fatalf("The commands %s is not supported!\n", command)
		}
	}

	// Code to measure
	duration := time.Since(start)

	fmt.Printf("Command done in %s.\n", duration)
}
