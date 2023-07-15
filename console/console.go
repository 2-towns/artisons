package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"gifthub/conf"
	"gifthub/console/parser"
	"gifthub/console/populate"
	"gifthub/users"
	"log"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

func main() {
	start := time.Now()

	if len(os.Args) == 1 {
		log.Panicln("The command is required, here are the possibilities: import")
	}

	command := os.Args[1]

	switch command {
	case "import":
		{
			file := flag.String("file", "./web/testdata/data.csv", "The path to the csv file")

			log.Printf("Will try to import %s csv, hop!", *file)

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

	case "userlist":
		{
			page := flag.Int64("page", 0, "The page used in pagination")

			u, err := users.List(*page)
			if err != nil {
				log.Panic(err)
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"ID", "First Name", "Last Name", "Email", "Phone", "City", "Updated at"})

			for _, user := range u {
				t.AppendRow([]interface{}{user.ID, user.Firstname, user.Lastname, user.Email, user.Phone, user.City, user.UpdatedAt})
			}

			t.Render()
		}
	default:
		{
			log.Panicf("The commands %s is not supported!\n", command)
		}
	}

	// Code to measure
	duration := time.Since(start)

	fmt.Printf("Command done in %s.\n", duration)
}
