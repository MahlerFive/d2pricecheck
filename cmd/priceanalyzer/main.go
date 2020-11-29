package main

import (
	"flag"
	"log"
	"os"

	"github.com/MahlerFive/d2pricecheck/internal/priceanalyzer"
)

func main() {
	// Read command-line args
	var tradesFilename string
	var outputFilename string
	var uniquesFilename string
	var setsFilename string
	flag.StringVar(&tradesFilename, "in", "input.txt", "Input file containing trade listings")
	flag.StringVar(&outputFilename, "out", "output.txt", "Output file which will contain all items and their price distributions")
	flag.StringVar(&uniquesFilename, "uniques", "data/uniques.txt", "File containing unique item names")
	flag.StringVar(&setsFilename, "sets", "data/sets.txt", "File containing set item names")
	flag.Parse()

	// Open input files
	tradeFile, err := os.Open(tradesFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer tradeFile.Close()

	uniquesFile, err := os.Open(uniquesFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer uniquesFile.Close()

	setsFile, err := os.Open(setsFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer setsFile.Close()

	// Open output file
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// Run the analyzer
	priceAnalyzer := priceanalyzer.NewPriceAnalyzer(tradeFile, uniquesFile, setsFile, outputFile)
	if err := priceAnalyzer.Analyze(); err != nil {
		log.Fatal(err)
	}
}
