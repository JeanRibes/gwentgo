package gwent

type BondMap map[string]int

func (bond BondMap) Add(card *Card) {
	if _, ok := bond[card.Name]; ok {
		bond[card.Name] += 1
	} else {
		bond[card.Name] = 1
	}
}

func (bond BondMap) Has(card Card) bool {
	_, ok := bond[card.Name]
	return ok
}

func (bond BondMap) Get(card *Card) int {
	return bond[card.Name]
}

func ToCards(list CardList) *Cards {
	out := &Cards{}
	for i, card := range list {
		card.Id = i
		out.Add(card)
	}
	return out
}

func CardPointers(in []Card) (out []*Card) {
	out = make([]*Card, len(in))
	for i, card := range in {
		out[i] = &card
	}
	return out
}
