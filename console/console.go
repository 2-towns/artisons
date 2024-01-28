package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"artisons/conf"
	"artisons/console/parser"
	"artisons/db"
	"artisons/logs"
	"artisons/notifications/mails"
	"artisons/notifications/vapid"
	"artisons/orders"
	"artisons/products"
	"artisons/users"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"golang.org/x/text/message"
)

func parseRedisFile(ctx context.Context, file string) [][]interface{} {
	f, err := os.ReadFile(file)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot open the file", slog.String("error", err.Error()))
		log.Fatal(err)

	}

	cmds := strings.Split(string(f), "\n")
	lines := [][]interface{}{}

	for _, line := range cmds {
		if line == "" {
			continue
		}

		args := []interface{}{}

		r := csv.NewReader(strings.NewReader(line))
		r.Comma = ' '
		fields, err := r.Read()

		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "error when parsing line", slog.String("line", line), slog.String("error", err.Error()))
			log.Fatalln(err)
		}

		for _, val := range fields {
			if val != "" {
				args = append(args, val)
			}
		}

		if strings.Contains(line, "updated_at") && !strings.Contains(line, "FT.CREATE") {
			args = append(args, "updated_at", time.Now().Unix())
		}

		lines = append(lines, args)
	}

	return lines
}

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
			file := flag.String("file", "./web/data/data.csv", "The path to the csv file")
			flag.Parse()

			f, err := os.Open(*file)
			if err != nil {
				slog.LogAttrs(ctx, slog.LevelError, "cannot open the file", slog.String("error", err.Error()))
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

	case "redis":
		{
			file := flag.String("file", "populate.redis", "The path to the populate file")

			flag.Parse()

			lines := parseRedisFile(ctx, *file)
			pipe := db.Redis.Pipeline()

			for _, line := range lines {
				pipe.Do(ctx, line...).Result()
			}

			_, err := pipe.Exec(ctx)

			if err != nil {
				slog.LogAttrs(ctx, slog.LevelError, "got error when populatin", slog.String("error", err.Error()))
				if err.Error() != "Unknown Index name" {
					log.Fatal(err)
				}
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
			msg := p.Sprintf("login_mail_otp", id, status)
			subject := p.Sprintf("email_order_update", id)
			mails.Send(ctx, user.Email, subject, msg)

			sessions, err := user.Sessions(ctx)
			if err != nil {
				slog.LogAttrs(ctx, slog.LevelError, "cannot get the session", slog.Int("uid", user.ID), slog.String("error", err.Error()))
				log.Fatal()
			}

			for _, session := range sessions {
				if session.WPToken != "" {
					vapid.Send(ctx, session.WPToken, msg)
				}
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

	// case "userlist":
	// 	{
	// 		page := flag.Int("page", 0, "The page used in pagination")

	// 		flag.Parse()

	// 		u, err := users.List(ctx, *page)
	// 		if err != nil {
	// 			slog.LogAttrs(ctx, slog.LevelError, "cannot list the users", slog.Int("page", *page), slog.String("error", err.Error()))
	// 			log.Fatalln()
	// 		}

	// 		t := table.NewWriter()
	// 		t.SetOutputMirror(os.Stdout)
	// 		t.AppendHeader(table.Row{"ID", "Email", "Updated at"})

	// 		for _, user := range u {
	// 			t.AppendRow([]interface{}{user.ID, user.Email, user.UpdatedAt})
	// 		}

	// 		t.Render()
	// 	}

	case "productdetail":
		{
			pid := flag.String("id", "", "The product id")

			flag.Parse()

			p, err := products.Find(ctx, *pid)
			if err != nil {
				slog.LogAttrs(ctx, slog.LevelError, "cannot get the product detail", slog.String("id", *pid), slog.String("error", err.Error()))
				log.Fatalln()
			}

			pjson, err := json.MarshalIndent(p, "", "  ")
			if err != nil {
				slog.LogAttrs(ctx, slog.LevelError, "cannot parse the object", slog.String("id", *pid), slog.String("error", err.Error()))
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
