package main

import (
	"fmt"
	"math/rand"

	"code.google.com/p/go-uuid/uuid"
	"github.com/DataWraith/gopush"
)

type Direction int8

const (
	North Direction = iota
	South Direction = iota
	East  Direction = iota
	West  Direction = iota
)

var DirectionNames = map[Direction]string{
	North: "North",
	South: "South",
	East:  "East",
	West:  "West",
}

type AgentList []*Agent

func (l AgentList) Len() int      { return len(l) }
func (l AgentList) Swap(a, b int) { l[a], l[b] = l[b], l[a] }
func (l AgentList) Less(a, b int) bool {
	if l[a].Points > l[b].Points {
		return true
	} else if l[a].Points < l[b].Points {
		return false
	} else {
		return l[a].code.Length < l[b].code.Length
	}
}

type Agent struct {
	Id        string
	World     *World
	X, Y      int
	Direction Direction
	Points    int

	interp *gopush.Interpreter
	code   gopush.Code
}

func NewAgent(w *World) *Agent {
	a := &Agent{
		Points: 5,
	}

	a.Id = uuid.New()
	a.World = w
	a.interp = gopush.NewInterpreter(gopush.DefaultOptions)
	a.interp.Options.EvalPushLimit = 10000

	stack := gopush.Stack{
		Functions: map[string]func(){
			"move": func() {
				if !a.interp.StackOK("integer", 1) {
					return
				}

				distance := a.interp.Stacks["integer"].Pop().(int)

				if a.World.Debug {
					fmt.Printf("[%s] move %d (%d) %s\n", a.Id, distance, a.Points, DirectionNames[a.Direction])
				}

				if distance > a.Points {
					distance = a.Points
				}

				switch a.Direction {
				case North:
					a.Y += distance
				case South:
					a.Y -= distance
				case East:
					a.X += distance
				case West:
					a.X -= distance
				}

				a.World.active = true
			},
			"turn": func() {
				if !a.interp.StackOK("integer", 1) {
					return
				}

				pd := a.Direction

				a.Direction += a.interp.Stacks["integer"].Pop().(Direction)
				if a.Direction < North {
					a.Direction = West
				}
				if a.Direction > West {
					a.Direction = North
				}

				if a.World.Debug {
					fmt.Printf("[%s] turned from %s to %s\n", a.Id, DirectionNames[pd], DirectionNames[a.Direction])
				}

				a.World.active = true
			},
			"available": func() {
				others := w.AgentsAt(a.X, a.Y)

				points := 0
				for _, o := range others {
					if o == a {
						continue
					}

					points += o.Points
				}

				if points != 0 {
					if a.World.Debug {
						fmt.Printf("[%s] found %d available points\n", a.Id, points)
					}
				}

				a.interp.Stacks["integer"].Push(points)
			},
			"consume": func() {
				others := w.AgentsAt(a.X, a.Y)

				for _, o := range others {
					if o == a {
						continue
					}

					if o.Points > 0 {
						if a.World.Debug {
							fmt.Printf("[%s] (%d points) is consuming a point from %s (%d points)\n", a.Id, a.Points, o.Id, o.Points)
						}

						o.Points--
						a.Points++

						break
					}
				}

				a.World.active = true
			},
		},
	}

	a.interp.Options.RegisterStack("preygent", &stack)
	a.interp.RegisterStack("preygent", &stack)

	a.code = a.interp.RandomCode(20)

	return a
}

func (a *Agent) Reset() *Agent {
	a.X = rand.Intn(10)
	a.Y = rand.Intn(10)
	a.Direction = North
	a.Points = 5

	for _, s := range a.interp.Stacks {
		s.Flush()
	}

	return a
}

func (a *Agent) Step() {
	a.interp.RunCode(a.code)
}
