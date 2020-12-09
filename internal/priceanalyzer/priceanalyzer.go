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
	// each line represents a single item with potentially many aliases, comma-separated

	for p.uniquesScanner.Scan() {
		line := p.uniquesScanner.Text()
		itemNames := strings.Split(line, ",")
		if len(itemNames) < 1 {
			return nil
		}

		// create an item corresponding to this line - all aliases will reference it
		item := NewItem(itemNames[0])

		for _, itemName := range itemNames {
			// create the searchable version of the item name to use as the map key
			itemNameSearchable := strings.ToLower(itemName)
			reg, err := regexp.Compile("[^a-z]+")
			if err != nil {
				log.Fatal(err)
			}
			itemNameSearchable = reg.ReplaceAllString(itemNameSearchable, "")

			// map the item name or alias to the item
			p.items[itemNameSearchable] = item
		}
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

	// remove common words/phrases that should be ignored
	i := 0
	for _, word := range words {
		// remove any offer/need words
		if strings.EqualFold(word, "o") || strings.EqualFold(word, "offer") ||
			strings.EqualFold(word, "n") || strings.EqualFold(word, "need") {
			continue
		}

		// remove "obo"
		if strings.EqualFold(word, "obo") {
			continue
		}

		// remove "rune" ending (eg. when price is "ist rune", ignore the "rune" part)
		// NOTE: removing all instances of "rune" should be safe, at least for now, since no unique/sets have "rune" in the name as a standalone word
		if strings.EqualFold(word, "rune") {
			continue
		}

		// if we reached here, we don't want to remove this word so put it back in the array
		words[i] = word
		i++
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
		if item.Output {
			continue
		}
		if len(item.PriceDistribution) < 1 {
			continue
		}

		p.output.WriteString(item.String())
		for runeNo, count := range item.PriceDistribution {
			p.output.WriteString(fmt.Sprintf("\t%s:%d", runesByNumber[runeNo], count))
		}
		p.output.WriteString("\n")

		item.Output = true
	}
}
