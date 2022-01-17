package gwent

var AllCardsMap map[string]*Card

func init() {
	AllCardsMap = map[string]*Card{}
	for key, card := range AllCards {
		c := AllCards[key]
		AllCardsMap[card.Name] = &c
	}
}
