package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"gifthub/conf"
	"gifthub/console/parser"
	"log"
	"os"
	"time"
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
	default:
		{
			log.Panicf("The commands %s is not supported!\n", command)
		}
	}

	// Code to measure
	duration := time.Since(start)

	fmt.Printf("Command done in %s.\n", duration)
}
