package main

import (
	"math"

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
	} else if l[a].Generation < l[b].Generation {
		return true
	} else if l[a].Generation > l[b].Generation {
		return false
	} else {
		return l[a].code.Length < l[b].code.Length
	}
}

type Agent struct {
	Id         string
	World      *World
	Generation int
	X, Y       int
	Direction  Direction
	Points     int

	interp *gopush.Interpreter
	code   gopush.Code
}

func NewAgent(w *World, g int) *Agent {
	a := &Agent{
		Points: 5,
	}

	a.Id = uuid.New()
	a.World = w
	a.Generation = g
	a.interp = gopush.NewInterpreter(gopush.DefaultOptions)
	a.interp.Options.EvalPushLimit = 100

	stack := gopush.Stack{
		Functions: map[string]func(){
			"move": func() {
				if !a.interp.StackOK("integer", 1) {
					return
				}

				distance := a.interp.Stacks["integer"].Pop().(int)

				a.World.Debugf("[%04d][%s] move %d (%d) %s", a.World.ticks, a.Id, distance, a.Points, DirectionNames[a.Direction])

				if distance > a.Points {
					distance = a.Points
				}

				a.X, a.Y = a.FocusAt(distance)

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

				a.World.Debugf("[%04d][%s] turned from %s to %s", a.World.ticks, a.Id, DirectionNames[pd], DirectionNames[a.Direction])

				a.World.active = true
			},
			"nearest": func() {
				nearest := int(100000000)

				for _, o := range w.agents {
					x, y := math.Abs(float64(a.X-o.X)), math.Abs(float64(a.Y-o.Y))

					distance := int(math.Sqrt(x*x + y*y))

					if distance < nearest {
						nearest = distance
					}
				}

				a.interp.Stacks["integer"].Push(nearest)
			},
			"available": func() {
				x, y := a.FocusAt(1)

				others := w.AgentsAt(x, y)

				points := 0
				for _, o := range others {
					if o == a {
						continue
					}

					points += o.Points
				}

				if points != 0 {
					a.World.Debugf("[%04d][%s] found %d points at %d,%d", a.World.ticks, a.Id, points, x, y)
				}

				a.interp.Stacks["integer"].Push(points)
			},
			"consume": func() {
				x, y := a.FocusAt(1)

				others := w.AgentsAt(x, y)

				for _, o := range others {
					if o == a {
						continue
					}

					if o.Points > 0 {
						a.World.Debugf("[%04d][%s] (%d) consuming %s (%d)", a.World.ticks, a.Id, a.Points, o.Id, o.Points)

						o.Points--
						a.Points++

						a.World.active = true

						break
					}
				}
			},
		},
	}

	a.interp.Options.RegisterStack("preygent", &stack)
	a.interp.RegisterStack("preygent", &stack)

	a.code = a.interp.RandomCode(20)

	return a
}

func (a *Agent) FocusAt(distance int) (int, int) {
	x, y := a.X, a.Y

	switch a.Direction {
	case East:
		x += distance
	case West:
		x -= distance
	case North:
		y += distance
	case South:
		y -= distance
	}

	return x, y
}

func (a *Agent) Reset() *Agent {
	a.X = 0
	a.Y = 0
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
