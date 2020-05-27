package mafia

import (
	"errors"
	"math/rand"
	"strings"
	"time"
)

type Game struct {
	Day     int
	Players []*Player

	killQueue []*Player
}

type Player struct {
	ID       string
	Role     string
	Index    int
	immunity bool
}

func NewGame() *Game {
	return new(Game)
}

func (g *Game) Add(id string) {
	player := new(Player)
	player.ID = id

	g.Players = append(g.Players, player)
}

func (g *Game) Pop(id string) *Player {
	i, player := g.GetPlayer(id)
	if player != nil {
		g.Players = append(g.Players[:i], g.Players[i+1:]...)
		return player
	}
	return nil
}

func (g *Game) Random(roles []string) error {
	if len(roles) != len(g.Players) {
		return errors.New("Количество ролей не соответствует колическу игроков")
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(roles), func(i, j int) { roles[i], roles[j] = roles[j], roles[i] })

	for _, player := range g.Players {
		player.Role, roles = strings.ToLower(roles[0]), roles[1:]
	}

	return nil
}

func (g *Game) Next() []*Player {
	g.Day++
	for _, player := range g.Players {
		player.immunity = false
	}
	defer func() { g.killQueue = nil }()
	return g.killQueue
}

func (g *Game) Kill(id string) error {
	if player := g.Pop(id); player != nil {
		g.killQueue = append(g.killQueue, player)
		return nil
	}
	return errors.New("Нет такого игрока")
}

func (g *Game) Immunity(id string) error {
	if _, player := g.GetPlayer(id); player != nil {
		player.immunity = true
		return nil
	}
	return errors.New("Нет такого игрока")
}

func (g *Game) Jail(id string) error {
	if _, player := g.GetPlayer(id); player != nil {
		if player.immunity {
			return errors.New("У данного игрока иммунитет")
		}
		g.Pop(id)
		return nil
	}
	return errors.New("Нет такого игрока")
}

func (g *Game) GetPlayer(id string) (int, *Player) {
	for i, player := range g.Players {
		if player.ID == id {
			return i, player
		}
	}
	return 0, nil
}

func (g *Game) GetPlayerByIndex(i int) (int, *Player) {
	for _, player := range g.Players {
		if player.Index == i {
			return i, player
		}
	}
	return 0, nil
}
