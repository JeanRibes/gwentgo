package gwent

type Ability enum

const (
	ScorchSiege Ability = iota
	HornSiege
	ClearWeatherEffect
	PickFogCard
	DrawEnemyDiscard
	CancelEnemyAbility
	LookEnemyHand
	PickRainCard
	ScorchCloseCombat
	HornRangedCombat
	DrawExtraBeginning
	PickFrostCard
	MedicEffect
	Discard2HandPick1Deck
	HornCloseCombat
	PickWeatherAny
)

type Leader struct {
	Name               string
	DisplayName        string
	Title              string
	Faction            Faction
	Image              string
	AbilityDescription string
	Ability            Ability
}

func (ability Ability) String() string {
	return strAbilityMap[ability]
}
func AbilityFromString(str string) Ability {
	return abilityStrMap[str]
}

func (ability Ability) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ability.String() + `""`), nil
}
func (ability *Ability) UnmarshalJSON(data []byte) (err error) {
	*ability = AbilityFromString(string(data))
	return nil
}

func init() {
	for str, ability := range abilityStrMap {
		strAbilityMap[ability] = str
	}
}

var abilityStrMap = map[string]Ability{
	"ScorchSiege":           ScorchSiege,
	"HornSiege":             HornSiege,
	"ClearWeatherEffect":    ClearWeatherEffect,
	"PickFogCard":           PickFogCard,
	"DrawEnemyDiscard":      DrawEnemyDiscard,
	"CancelEnemyAbility":    CancelEnemyAbility,
	"LookEnemyHand":         LookEnemyHand,
	"PickRainCard":          PickRainCard,
	"ScorchCloseCombat":     ScorchCloseCombat,
	"HornRangedCombat":      HornRangedCombat,
	"DrawExtraBeginning":    DrawExtraBeginning,
	"PickFrostCard":         PickFrostCard,
	"MedicEffect":           MedicEffect,
	"Discard2HandPick1Deck": Discard2HandPick1Deck,
	"HornCloseCombat":       HornCloseCombat,
	"PickWeatherAny":        PickWeatherAny,
}

var strAbilityMap = map[Ability]string{}

var Leaders = map[string]*Leader{
	"foltest-the-steel-forged": {
		Faction:            NorthernRealms,
		Name:               "foltest-the-steel-forged",
		DisplayName:        "Foltest",
		Title:              "The Steel-Forged",
		Image:              "Foltest-the-Steel-Forged-gwent-card.jpg",
		AbilityDescription: "Scorch Siege if enemies Siege strengh is 10 or higher",
		Ability:            ScorchSiege,
	},
	"foltest-the-siegemaster": {
		Faction:            NorthernRealms,
		Name:               "foltest-the-siegemaster",
		DisplayName:        "Foltest",
		Title:              "The Siegemaster",
		Image:              "Foltest-the-Siegemaster-gwent-card.jpg",
		AbilityDescription: "Horn on Siege row",
		Ability:            HornSiege,
	},
	"foltest-lord-commander-of-the-north": {
		Faction:            NorthernRealms,
		Name:               "foltest-lord-commander-of-the-north",
		DisplayName:        "Foltest",
		Title:              "Lord Commander of the North",
		Image:              "Foltest-Lord-Commander-of-the-North-gwent-card.jpg",
		AbilityDescription: "Clear any Weather effects in game",
		Ability:            ClearWeatherEffect,
	},
	"foltest-king-of-temeria": {
		Faction:            NorthernRealms,
		Name:               "foltest-king-of-temeria",
		DisplayName:        "Foltest",
		Title:              "King of Temeria",
		Image:              "Foltest-King-of-Temeria-gwent-card.jpg",
		AbilityDescription: "Pick a Fog card from your deck and play it immediately",
		Ability:            PickFogCard,
	},
	"emhyr-var-emreis-the-relentless": {
		Faction:            Nilfgaard,
		Name:               "emhyr-var-emreis-the-relentless",
		DisplayName:        "Emhyr var Emreis",
		Title:              "The Relentless",
		Image:              "Emhyr-var-Emreis-the-Relentless-gwen-card.jpg",
		AbilityDescription: "Draw card from opponent’s discard pile",
		Ability:            DrawEnemyDiscard,
	},
	"emhyr-var-emreis-the-white-flame": {
		Faction:            Nilfgaard,
		Name:               "emhyr-var-emreis-the-white-flame",
		DisplayName:        "Emhyr var Emreis",
		Title:              "The White Flame",
		Image:              "Emhyr-var-Emreis-the-White-Flame-gwent-card.jpg",
		AbilityDescription: "Cancel your opponents Leader ability",
		Ability:            CancelEnemyAbility,
	},
	"emhyr-var-emreis-the-emperor-of-nilfgaard": {
		Faction:            Nilfgaard,
		Name:               "emhyr-var-emreis-the-emperor-of-nilfgaard",
		DisplayName:        "Emhyr var Emreis",
		Title:              "The Emperor of Nilfgaard",
		Image:              "Emhyr-var-Emreis-Emperor-of-Nilfgaard-gwent-card.jpg",
		AbilityDescription: "Look at 3 random cards of your opponents hand",
		Ability:            LookEnemyHand,
	},
	"emhyr-var-emreis-his-imperial-majesty": {
		Faction:            Nilfgaard,
		Name:               "emhyr-var-emreis-his-imperial-majesty",
		DisplayName:        "Emhyr var Emreis",
		Title:              "His Imperial Majesty",
		Image:              "Emhyr-var-Emreis-His-Imperial-Majesty-gwent-card.jpg",
		AbilityDescription: "Pick a Rain card from your deck and play it immediately",
		Ability:            PickRainCard,
	},
	"francesca-findabair-queen-of-dol-blathanna": {
		Faction:            ScoiaTael,
		Name:               "francesca-findabair-queen-of-dol-blathanna",
		DisplayName:        "Francesca Findabair",
		Title:              "Queen of Dol Blathanna",
		Image:              "Francesca-Findabair-Queen-of-Dol-Blathanna-gwent-card.jpg",
		AbilityDescription: "Scorch Close Combat if enemies Close Combat strengh is 10 or higher",
		Ability:            ScorchCloseCombat,
	},
	"francesca-findabair-the-beautiful": {
		Faction:            ScoiaTael,
		Name:               "francesca-findabair-the-beautiful",
		DisplayName:        "Francesca Findabair",
		Title:              "The Beautiful",
		Image:              "Francesca-Findabair-the-Beautiful-gwent-card.jpg",
		AbilityDescription: "Horn effect on your Ranged Combat row",
		Ability:            HornRangedCombat,
	},
	"francesca-findabair-daisy-of-the-valley": {
		Faction:            ScoiaTael,
		Name:               "francesca-findabair-daisy-of-the-valley",
		DisplayName:        "Francesca Findabair",
		Title:              "Daisy of the Valley",
		Image:              "Francesca-Findabair-Daisy-of-the-Valley-gwent-card.jpg",
		AbilityDescription: "Draw extra card at the beginning of the game",
		Ability:            DrawExtraBeginning,
	},
	"francesca-findabair-pureblood-elf": {
		Faction:            ScoiaTael,
		Name:               "francesca-findabair-pureblood-elf",
		DisplayName:        "Francesca Findabair",
		Title:              "Pureblood Elf",
		Image:              "Francesca-Findabair-Pureblood-Elf-gwent-card.jpg",
		AbilityDescription: "Pick a Frost card from your deck and play it immediately",
		Ability:            PickFrostCard,
	},
	"eredin-destroyer-of-worlds": {
		Faction:            Monsters,
		Name:               "eredin-destroyer-of-worlds",
		DisplayName:        "Eredin",
		Title:              "Destroyer of Worlds",
		Image:              "Eredin-Destroyer-of-Worlds-gwent-card.jpg",
		AbilityDescription: "Pick a card from your discard pile and put it back into your hand",
		Ability:            MedicEffect,
	},
	"eredin-bringer-of-death": {
		Faction:            Monsters,
		Name:               "eredin-bringer-of-death",
		DisplayName:        "Eredin",
		Title:              "Bringer of Death",
		Image:              "Eredin-Bringer-of-Death-gwent-card.jpg",
		AbilityDescription: "Discard 2 cards from your hand, Draw 1 card of your choice from you deck",
		Ability:            Discard2HandPick1Deck,
	},
	"eredin-king-of-the-wild-hunt": {
		Faction:            Monsters,
		Name:               "eredin-king-of-the-wild-hunt",
		DisplayName:        "Eredin",
		Title:              "King of the Wild Hunt",
		Image:              "Eredin-King-of-the-Wild-Hunt-gwent-card.jpg",
		AbilityDescription: "Horn effect on your Close Combat row",
		Ability:            HornCloseCombat,
	},
	"eredin-commander-of-the-red-riders": {
		Faction:            Monsters,
		Name:               "eredin-commander-of-the-red-riders",
		DisplayName:        "Eredin",
		Title:              "Commander of the Red Riders",
		Image:              "Eredin-Commander-of-the-Red-Riders-gwent-card.jpg",
		AbilityDescription: "Pick any weather card from your deck and play it immediately",
		Ability:            PickWeatherAny,
	},
}
