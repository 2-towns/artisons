package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"gifthub/conf"
	"gifthub/console/parser"
	"gifthub/console/populate"
	"gifthub/logs"
	"gifthub/notifications/mails"
	"gifthub/notifications/vapid"
	"gifthub/orders"
	"gifthub/products"
	"gifthub/users"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/text/message"
)

// var (
// 	printer = message.NewPrinter(locales.Console)
// )

func main() {
	logs.Init()

	start := time.Now()

	if len(os.Args) == 1 {
		log.Fatalln("The command is required, here are the possibilities: import")
	}

	command := os.Args[len(os.Args)-1]
	ctx := context.Background()

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
				slog.LogAttrs(ctx, slog.LevelError, "cannot parse the csv", slog.String("error", err.Error()))
				log.Fatal()
			}

			lines, err := parser.Import(data, conf.DefaultMID)
			if err != nil {
				slog.LogAttrs(ctx, slog.LevelError, "cannot import the csv", slog.String("error", err.Error()))
				log.Fatal()
			}

			slog.LogAttrs(ctx, slog.LevelInfo, "import successful", slog.Int("lines", lines))

		}

	case "populate":
		{
			err := populate.Run()

			if err != nil {
				slog.LogAttrs(ctx, slog.LevelError, "cannot populate", slog.String("error", err.Error()))
				log.Fatal()
			}
		}

	case "orderstatus":
		{
			id := flag.String("id", "", "The order id")
			status := flag.String("status", "", "The new order status")

			flag.Parse()

			err := orders.UpdateStatus(ctx, *id, *status)
			if err != nil {
				slog.LogAttrs(ctx, slog.LevelError, "cannot update the order", slog.String("oid", *id), slog.String("error", err.Error()))
				log.Fatal()
			}

			order, err := orders.Find(ctx, *id)
			if err != nil {
				slog.LogAttrs(ctx, slog.LevelError, "cannot find the order", slog.String("oid", *id), slog.String("error", err.Error()))
				log.Fatal()
			}

			user, err := users.Get(ctx, order.UID)
			if err != nil {
				slog.LogAttrs(ctx, slog.LevelError, "cannot get the user", slog.String("oid", *id), slog.String("error", err.Error()))
				log.Fatal()
			}

			p := message.NewPrinter(user.Lang)
			msg := p.Sprintf("mail_magic_link", id, status)
			mails.Send(ctx, user.Email, msg)

			for _, value := range user.Devices {
				vapid.Send(ctx, value, msg)
			}
		}

	case "orderdetail":
		{
			id := flag.String("id", "", "The order id")

			flag.Parse()

			o, err := orders.Find(ctx, *id)
			if err != nil {
				slog.LogAttrs(ctx, slog.LevelError, "cannot find the order", slog.String("id", *id), slog.String("error", err.Error()))
				log.Fatal()
			}

			empJSON, err := json.MarshalIndent(o, "", "  ")
			if err != nil {
				slog.LogAttrs(ctx, slog.LevelError, "cannot parse the object", slog.String("id", *id), slog.String("error", err.Error()))
				log.Fatal()
			}

			log.Printf("%s\n", string(empJSON))
		}

	case "ordernote":
		{
			id := flag.String("id", "", "The order id")
			note := flag.String("note", "", "The note to attach")

			flag.Parse()

			err := orders.AddNote(ctx, *id, *note)
			if err != nil {
				slog.LogAttrs(ctx, slog.LevelError, "cannot add a note", slog.String("id", *id), slog.String("error", err.Error()))
				log.Fatalln()
			}

			slog.LogAttrs(ctx, slog.LevelInfo, "note added to the order")
		}

	case "userlist":
		{
			page := flag.Int("page", 0, "The page used in pagination")

			flag.Parse()

			u, err := users.List(ctx, *page)
			if err != nil {
				slog.LogAttrs(ctx, slog.LevelError, "cannot list the users", slog.Int("page", *page), slog.String("error", err.Error()))
				log.Fatalln()
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"ID", "Email", "Updated at"})

			for _, user := range u {
				t.AppendRow([]interface{}{user.ID, user.Email, user.UpdatedAt})
			}

			t.Render()
		}

	case "productdetail":
		{
			pid := flag.String("pid", "", "The product id")

			flag.Parse()

			p, err := products.Find(ctx, *pid)
			if err != nil {
				slog.LogAttrs(ctx, slog.LevelError, "cannot get the product detail", slog.String("pid", *pid), slog.String("error", err.Error()))
				log.Fatalln()
			}

			pjson, err := json.MarshalIndent(p, "", "  ")
			if err != nil {
				slog.LogAttrs(ctx, slog.LevelError, "cannot parse the object", slog.String("pid", *pid), slog.String("error", err.Error()))
				log.Fatal()
			}

			log.Printf("%s\n", string(pjson))
		}
	default:
		{
			slog.LogAttrs(ctx, slog.LevelError, "the command is not supported", slog.String("command", command))
		}
	}

	// Code to measure
	duration := time.Since(start)

	slog.LogAttrs(ctx, slog.LevelInfo, "command done", slog.Duration("duration", duration))
}
