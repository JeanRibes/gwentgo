// Package gwent, API for a Gwynt game logic and data structures
// It aims to reproduce the mini-game found in The Witcher 3 Wild Hunt, property of CDProjectRed
package gwent

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"
)

var AllCardsMap map[string]*Card
var AllCardsList = &CardList{}

func init() {
	*AllCardsList = loadcards()
	AllCardsMap = map[string]*Card{}
	for _, card := range *AllCardsList {
		AllCardsMap[card.Name] = card
	}
}

func loadcards() CardList {
	file, err := os.Open("cards.csv")
	if err != nil {
		panic(err)
	}
	cards, err := LoadCardsCSV(file)
	if err != nil {
		panic(err)
	}
	return cards
}

func LoadCardsCSV(r io.Reader) (cards CardList, err error) {
	reader := csv.NewReader(r)
	reader.Read() //skip headers
	id := 0
	for {
		row, err := reader.Read()
		if err == io.EOF || row == nil {
			break
		}
		if err != nil {
			return nil, err
		}

		faction, err := FactionFromString(row[4])
		strength, err := strconv.ParseInt(row[2], 10, 64)
		if err != nil {
			return nil, err
		}

		effects := Effects{}
		for _, ability := range strings.Split(row[3], ",") {
			if ability == "" {
				continue
			}
			if effect, err := EffectFromString(ability); err != nil {
				return nil, err
			} else {
				effects = append(effects, effect)
			}
		}
		card := Card{
			Id:          id,
			Faction:     faction,
			Row:         RowFromString(row[5]),
			Name:        row[0],
			DisplayName: row[1],
			Image:       row[6],
			Strength:    int(strength),
			Effects:     effects,
			score:       0,
		}
		id += 1
		cards = append(cards, &card)
	}
	return cards, nil
}
