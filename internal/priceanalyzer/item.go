package priceanalyzer

type Item struct {
	FriendlyName string
	// PriceDistribution is a map from the rune # to amount
	PriceDistribution map[int]int
	// Output marks whether this Item was already output or not
	Output bool
}

func NewItem(itemFriendlyName string) *Item {
	return &Item{
		FriendlyName:      itemFriendlyName,
		PriceDistribution: make(map[int]int),
	}
}

func (i Item) String() string {
	return i.FriendlyName
}
