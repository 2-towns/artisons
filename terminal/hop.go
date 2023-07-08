package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"gifthub/conf"
	"gifthub/util"
	"log"
	"os"
	"time"
)

func main() {
	file := flag.String("file", "./static/fake/data.csv", "The path to the csv file")

	log.Printf("Will try to import %s csv, hop!", *file)

	f, err := os.Open(*file)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	start := time.Now()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	lines, err := util.CsvImport(data, conf.DefaultMID)
	if err != nil {
		log.Panicln(err)
	}

	// Code to measure
	duration := time.Since(start)

	fmt.Printf("Import successful, %d line(s) imported in %s.\n", lines, duration)
}
