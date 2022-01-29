package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gwentgo/gwent"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type MpEvent int

const (
	PlayerJoin MpEvent = iota
	DeckChosen
	GameReady
)

type GameEvent int

const (
	GameFinished GameEvent = iota
	GameUpdated
	TurnChanged
	EnemyPassed
)

type MultiplayerGame struct {
	Key          string
	LobbyChannel [2]chan MpEvent
	GameChannel  [2]chan GameEvent
	sync.Mutex
	Game *gwent.Game

	ChanA chan bool
	UserA *gwent.PlayerData
	DeckA *gwent.PlayerDeck

	ChanB chan bool
	UserB *gwent.PlayerData
	DeckB *gwent.PlayerDeck
}
type mpstorage struct {
	Key  string
	Game *gwent.Game

	UserA *gwent.PlayerData
	DeckA *gwent.PlayerDeck

	UserB *gwent.PlayerData
	DeckB *gwent.PlayerDeck
}

var multiplayer_games = map[string]*MultiplayerGame{}

// createMPgame url: /multi/create
func createMPgame(c *gin.Context) {
	key := RandomString(32)
	mpg := &MultiplayerGame{
		Key: key,
	}
	mpg.LobbyChannel[gwent.PlayerA] = make(chan MpEvent, 3)
	mpg.LobbyChannel[gwent.PlayerB] = make(chan MpEvent, 3)

	mpg.GameChannel[gwent.PlayerA] = make(chan GameEvent, 4)
	mpg.GameChannel[gwent.PlayerB] = make(chan GameEvent, 4)
	multiplayer_games[key] = mpg

	setcoookie(c, GKEY, key)
	delcookie(c, SIDE)
	c.Redirect(http.StatusFound, "/multi/join")
}

const MP = "multiplayer_game"
const GKEY = "gamekey"

func MultiGameMiddleware(c *gin.Context) {
	cookie, err := c.Cookie(GKEY)
	if err != nil || cookie == "" {
		c.String(400, "missing game cookie")
		c.Abort()
		return
	}
	lobby, ok := multiplayer_games[cookie]
	if !ok {
		c.String(400, "no game found")
		c.Abort()
		return
	}

	if _side, err := c.Cookie(SIDE); err == nil {
		side := gwent.UnstringTurn(_side)
		if side != gwent.Tie {
			c.Set(SIDE, side)
		}
	}

	c.Set(MP, lobby)
	c.Next()
}

const SIDE = "side"

// joinGameP1 for the first player, who initiated the lobby
func joinGameP1(c *gin.Context) { // gwent.PlayerA
	user := c.MustGet("user").(*gwent.PlayerData)
	key, _ := c.Cookie(GKEY)
	lobby, ok := multiplayer_games[key]
	if !ok {
		c.String(400, "no game found, bad key")
		return
	}
	if side, _ := c.Cookie(SIDE); side == "" {
		setcoookie(c, SIDE, lobby.AssignUser(user).String())
	}
	c.HTML(200, "waiting.html", gin.H{"Key": key, "User": user})
}

func joinGameP2(c *gin.Context) { // gwent.PlayerB
	user := c.MustGet("user").(*gwent.PlayerData)
	key := c.Param(GKEY)
	lobby, ok := multiplayer_games[key]
	if !ok {
		c.String(400, "no game found")
		return
	}

	if ckey, _ := c.Cookie(GKEY); ckey != key { // first time visiting this URL
		setcoookie(c, GKEY, key)
		lobby.SignalJoined()
		setcoookie(c, SIDE, lobby.AssignUser(user).String())
	}
	lobby.LobbyChannel[gwent.PlayerB] <- PlayerJoin
	c.HTML(200, "waiting.html", gin.H{"Key": key, "User": user})
}

// waitingRoom url: /multi/wait
func waitingRoom(c *gin.Context) {
	lobby := c.MustGet(MP).(*MultiplayerGame)
	side := c.MustGet(SIDE).(gwent.Turn)
	c.Stream(func(w io.Writer) bool {
		event := <-lobby.LobbyChannel[side]
		switch event {
		case PlayerJoin:
			c.SSEvent("joined", "opponent ("+side.Enemy().String()+") joined !")
			break
		case DeckChosen:
			c.SSEvent("deck", "opponent ("+side.Enemy().String()+") has chosen its deck")
			break
		case GameReady:
			c.SSEvent("ready", "game is ready")
			return false
		}
		return true
	})
}

func chooseDeck(c *gin.Context) {
	user := c.MustGet(USER).(*gwent.PlayerData)
	lobby := c.MustGet(MP).(*MultiplayerGame)
	side := c.MustGet(SIDE).(gwent.Turn)

	sindex := c.Param("index")
	index, err := strconv.ParseInt(sindex, 10, 64)
	if err != nil {
		panic(err)
		return
	}
	deck := user.Decks.GetByIndex(int(index))
	if deck == nil {
		panic("invalid index")
		return
	}
	if lobby.SetDeck(deck, side) {
		c.String(201, "you chose %s", deck.Name)
	} else {
		c.String(200, "Error: that deck is not eligible for gameplay, check requirements (>22 unit cards)")
	}
}

func (mg *MultiplayerGame) IsFull() bool {
	mg.Lock()
	defer mg.Unlock()
	return mg.UserA != nil && mg.UserB != nil
}
func (mg *MultiplayerGame) AssignUser(user *gwent.PlayerData) gwent.Turn {
	mg.Lock()
	defer mg.Unlock()
	if mg.UserA == nil {
		mg.UserA = user
		return gwent.PlayerA
	}
	if mg.UserB == nil {
		mg.UserB = user
		return gwent.PlayerB
	}
	return gwent.Tie
}
func (mg *MultiplayerGame) SetDeck(deck *gwent.PlayerDeck, side gwent.Turn) bool {
	if !deck.Eligible() {
		return false
	}
	mg.Lock()
	if side == gwent.PlayerA {
		mg.DeckA = deck
	}
	if side == gwent.PlayerB {
		mg.DeckB = deck
	}
	mg.LobbyChannel[side.Enemy()] <- DeckChosen
	mg.Unlock()
	mg.TryStartGame()
	return true
}
func (mg *MultiplayerGame) TryStartGame() {
	mg.Lock()
	defer mg.Unlock()
	if mg.DeckA != nil && mg.DeckB != nil && mg.Game == nil {
		game, err := gwent.GameFromDecks(mg.DeckA, mg.DeckB)
		if err != nil {
			log.Println(err.Error())
			return
		}
		mg.Game = game
		mg.SignalReady()
		log.Println("game starting")
	}
	log.Println("not ready,not starting game")
}
func (mg *MultiplayerGame) Started() bool {
	mg.Lock()
	defer mg.Unlock()
	return mg.Game != nil
}
func (mg *MultiplayerGame) SignalJoined() {
	mg.LobbyChannel[gwent.PlayerA] <- PlayerJoin
}
func (mg *MultiplayerGame) SignalReady() {
	mg.LobbyChannel[gwent.PlayerA] <- GameReady
	mg.LobbyChannel[gwent.PlayerB] <- GameReady
}

func savempg() {
	f, err := os.OpenFile("games_db.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(&multiplayer_games)
	if err != nil {
		panic(err)
	}

	if _, err := f.Write(data); err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
}
func loadmpg() {
	f, err := os.Open("games_db.json")
	if err != nil {
		log.Printf("open error: %s", err.Error())
		return
	}
	if err := json.NewDecoder(f).Decode(&multiplayer_games); err != nil {
		log.Println(err)
	}
}

/*
func (ge GameEvent) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", int(ge))), nil
}
func (ge *GameEvent) UnmarshalJSON(data []byte) error {
	_, err := fmt.Sscanf("%d", string(data), ge)
	return err
}
func (e MpEvent) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", int(e))), nil
}
func (e *MpEvent) UnmarshalJSON(data []byte) error {
	_, err := fmt.Sscanf("%d", string(data), e)
	return err
}*/

func (mpg *MultiplayerGame) MarshalJSON() ([]byte, error) {
	alt := mpstorage{
		Key:   mpg.Key,
		Game:  mpg.Game,
		UserA: mpg.UserA,
		DeckA: mpg.DeckA,
		UserB: mpg.UserB,
		DeckB: mpg.DeckB,
	}
	return json.Marshal(alt)
}
func (mpg *MultiplayerGame) UnmarshalJSON(data []byte) error {
	var alt mpstorage
	err := json.Unmarshal(data, &alt)
	mpg.Key = alt.Key
	mpg.Game = alt.Game
	mpg.UserA = alt.UserA
	mpg.DeckA = alt.DeckA
	mpg.UserB = alt.UserB
	mpg.DeckB = alt.DeckB

	mpg.LobbyChannel[gwent.PlayerA] = make(chan MpEvent, 3)
	mpg.LobbyChannel[gwent.PlayerB] = make(chan MpEvent, 3)

	mpg.GameChannel[gwent.PlayerA] = make(chan GameEvent, 4)
	mpg.GameChannel[gwent.PlayerB] = make(chan GameEvent, 4)

	return err
}

// ======================================================================================
// gamplay part below

func gameHandler(c *gin.Context) {
	lobby := c.MustGet(MP).(*MultiplayerGame)
	side := c.MustGet(SIDE).(gwent.Turn)
	game := lobby.Game
	player := game.Side(side)
	enemy := game.Side(side.Enemy())
	println("ajax ?: ", c.GetHeader("HX-Request"))
	if c.GetHeader("HX-Request") == "true" {
		displayMPboard(c, "/multi/game", "table.html", false, game, player, enemy)
	} else {
		displayMPboard(c, "/multi/game", "mpgame.html", false, game, player, enemy)
	}
}

func moveHandler(c *gin.Context) {
	lobby := c.MustGet(MP).(*MultiplayerGame)
	side := c.MustGet(SIDE).(gwent.Turn)
	row, cardid := getRowAndCard(c)
	game := lobby.Game
	player := game.Side(side)
	enemy := game.Side(side.Enemy())

	if game.Turn != side {
		c.String(400, "please wait your turn, i see you trying to cheat !")
		return
	}

	if cardid < 0 {
		game.Pass(player)
		lobby.GameChannel[side.Enemy()] <- EnemyPassed
		lobby.GameChannel[side.Enemy()] <- TurnChanged
		lobby.GameChannel[side] <- TurnChanged
		lobby.GameChannel[side.Enemy()] <- GameUpdated
		displayMPboard(c, "/multi/game",
			"table.html", false,
			game, player, enemy)
		if game.Finished() {
			lobby.GameChannel[side.Enemy()] <- GameFinished
			lobby.GameChannel[side] <- GameFinished
		}
		return
	}

	card := player.GetCardById(cardid)
	if cardid >= 0 && card == nil {
		return
	}
	medicEffect := game.PlayMove(card, row, player, enemy)
	if medicEffect {
	} else {
		game.Switch()
		lobby.GameChannel[side.Enemy()] <- TurnChanged
		lobby.GameChannel[side] <- TurnChanged
	}

	lobby.GameChannel[side.Enemy()] <- GameUpdated // signal other side
	displayMPboard(c, "/multi/game",
		"table.html", medicEffect,
		game, player, enemy)
}

func gameEventListener(c *gin.Context) {
	lobby := c.MustGet(MP).(*MultiplayerGame)
	side := c.MustGet(SIDE).(gwent.Turn)
	channel := lobby.GameChannel[side]

	c.Stream(func(w io.Writer) bool {
		event := <-channel
		switch event {
		case TurnChanged:
			c.SSEvent("TurnChanged", lobby.Game.Turn.String())
		default:
			c.SSEvent(event.String(), "event "+event.String())
		}
		return true
	})
}

func (ge GameEvent) String() string {
	if ge == GameFinished {
		return "GameFinished"
	}
	if ge == GameUpdated {
		return "GameUpdated"
	}
	if ge == EnemyPassed {
		return "EnemyPassed"
	}
	return "error"
}

func displayMPboard(
	c *gin.Context,
	url string,
	template string,
	choice bool,
	game *gwent.Game,
	player *gwent.GameSide,
	enemy *gwent.GameSide) {
	weather := game.WeatherCards.Effects()
	c.HTML(200, template, gin.H{
		"Weather": gin.H{
			"CardList":     game.WeatherCards,
			"CloseCombat":  weather.Has(gwent.BitingFrost),
			"RangedCombat": weather.Has(gwent.ImpenetrableFog),
			"Siege":        weather.Has(gwent.TorrentialRain),
		},
		"Lives": gin.H{
			"Player": game.LivesLeft(player.Side),
			"Enemy":  game.LivesLeft(enemy.Side),
		},
		"MySide": player,
		"Player": player,
		"Enemy":  enemy,
		"Url":    url,
		"Choice": choice,
		"Round":  game.Round(),
		"Turn":   game.Turn.String(),
		"Side":   player.Side.String(),
	})
}
