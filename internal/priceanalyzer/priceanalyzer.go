package priceanalyzer

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type PriceAnalyzer struct {
	tradesScanner  *bufio.Scanner
	uniquesScanner *bufio.Scanner
	setsScanner    *bufio.Scanner

	output io.StringWriter

	items map[string]*Item
}

func NewPriceAnalyzer(trades io.Reader, uniques io.Reader, sets io.Reader, output io.StringWriter) *PriceAnalyzer {
	return &PriceAnalyzer{
		tradesScanner:  bufio.NewScanner(trades),
		uniquesScanner: bufio.NewScanner(uniques),
		setsScanner:    bufio.NewScanner(sets),
		output:         output,
		items:          make(map[string]*Item),
	}
}

func (p *PriceAnalyzer) Analyze() error {
	// Load item lists
	if err := p.loadUniques(); err != nil {
		return errors.Wrap(err, "error loading uniques")
	}
	if err := p.loadSets(); err != nil {
		return errors.Wrap(err, "error loading sets")
	}

	// Parse trades text and calculate prices
	for p.tradesScanner.Scan() {
		line := p.tradesScanner.Text()
		fmt.Printf("Processing line: \"%s\"\n", line)

		item, price, valid := p.parseLine(line)
		if !valid {
			continue
		}

		if _, ok := item.PriceDistribution[runesByName[price]]; !ok {
			item.PriceDistribution[runesByName[price]] = 1
		} else {
			item.PriceDistribution[runesByName[price]] += 1
		}
	}

	if err := p.tradesScanner.Err(); err != nil {
		return errors.Wrap(err, "error during trades scanning")
	}

	// Write results
	p.writeOutput()

	return nil
}

func (p *PriceAnalyzer) loadUniques() error {
	for p.uniquesScanner.Scan() {
		itemName := p.uniquesScanner.Text()

		// create the searchable version of the item name to use as the map key
		itemNameSearchable := strings.ToLower(itemName)
		reg, err := regexp.Compile("[^a-z]+")
		if err != nil {
			log.Fatal(err)
		}
		itemNameSearchable = reg.ReplaceAllString(itemNameSearchable, "")

		p.items[itemNameSearchable] = NewItem(itemName)
	}

	return errors.Wrap(p.tradesScanner.Err(), "error scanning uniques")
}

func (p *PriceAnalyzer) loadSets() error {
	for p.setsScanner.Scan() {
		itemName := p.setsScanner.Text()

		// create the searchable version of the item name to use as the map key
		itemNameSearchable := strings.ToLower(itemName)
		reg, err := regexp.Compile("[^a-z]+")
		if err != nil {
			log.Fatal(err)
		}
		itemNameSearchable = reg.ReplaceAllString(itemNameSearchable, "")

		p.items[itemNameSearchable] = NewItem(itemName)
	}

	return errors.Wrap(p.tradesScanner.Err(), "error scanning sets")
}

// parseLine parses a single line and returns the item and price (the rune name as string).
// If the line doesn't have a valid item and price, it false for the valid flag.
func (p *PriceAnalyzer) parseLine(line string) (*Item, string, bool) {
	line = strings.ToLower(line)

	// convert all non-alpha to spaces
	reg, err := regexp.Compile("[^a-z]+")
	if err != nil {
		log.Fatal(err)
	}
	line = reg.ReplaceAllString(line, " ")

	// remove leading/trailing spaces
	line = strings.TrimLeft(line, " ")
	line = strings.TrimRight(line, " ")

	// split by spaces
	words := strings.Split(line, " ")

	// need at least an item and a price
	if len(words) < 2 {
		return nil, "", false
	}

	// remove any offer/need prefix
	if strings.EqualFold(words[0], "o") || strings.EqualFold(words[0], "offer") ||
		strings.EqualFold(words[0], "n") || strings.EqualFold(words[0], "need") {
		words = words[1:]
	}

	// need at least an item and a price
	if len(words) < 2 {
		return nil, "", false
	}

	// remove "obo" ending
	if strings.EqualFold(words[len(words)-1], "obo") {
		words = words[:len(words)-1]
	}

	// need at least an item and a price
	if len(words) < 2 {
		return nil, "", false
	}

	// extract price which must match a rune name
	price := words[len(words)-1]
	if _, ok := runesByName[price]; !ok {
		return nil, "", false
	}

	// search for an item name match
	// chop off parts from the end of the name until a match is found, or we run out of name parts
	itemNameFound := false
	var item *Item
	for numItemNameParts := len(words) - 1; numItemNameParts > 0; numItemNameParts-- {
		itemName := strings.Join(words[:numItemNameParts], "")
		var ok bool
		if item, ok = p.items[itemName]; ok {
			itemNameFound = true
			break
		}
	}

	if !itemNameFound {
		return nil, "", false
	}

	fmt.Printf("\tMatched item = \"%s\" with price = %s\n", item, price)

	return item, price, true
}

func (p *PriceAnalyzer) writeOutput() {
	// Print out all items and price distributions
	fmt.Printf("\n\nITEM PRICES\n===========\n\n")
	for _, item := range p.items {
		if len(item.PriceDistribution) < 1 {
			continue
		}

		fmt.Printf("%s", item)
		for runeNo, count := range item.PriceDistribution {
			fmt.Printf("\t%s:%d", runesByNumber[runeNo], count)
		}
		fmt.Println()
	}

	// Write all items and price distribution to output
	for _, item := range p.items {
		if len(item.PriceDistribution) < 1 {
			continue
		}

		p.output.WriteString(item.String())
		for runeNo, count := range item.PriceDistribution {
			p.output.WriteString(fmt.Sprintf("\t%s:%d", runesByNumber[runeNo], count))
		}
		p.output.WriteString("\n")
	}
}
