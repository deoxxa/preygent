package main

import (
	"math"
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
	a := &Agent{}

	a.Id = uuid.New()
	a.World = w
	a.Generation = g
	a.interp = gopush.NewInterpreter(gopush.DefaultOptions)
	a.interp.Options.EvalPushLimit = 100
	a.interp.Options.TopLevelPopCode = true

	stack := gopush.Stack{
		Functions: map[string]func(){
			"forward": func() {
				a.X, a.Y = a.FocusAt(1)

				a.World.active = true
			},
			"left": func() {
				a.Direction -= 1
				if a.Direction < North {
					a.Direction = West
				}
			},
			"right": func() {
				a.Direction += 1
				if a.Direction > West {
					a.Direction = North
				}
			},
			"move": func() {
				if !a.interp.StackOK("integer", 1) {
					return
				}

				distance := a.interp.Stacks["integer"].Pop().(int)

				a.World.Debugf("[%04d][%s] move %d (%d) %s", a.World.iterations, a.Id, distance, a.Points, DirectionNames[a.Direction])

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

				a.World.Debugf("[%04d][%s] turned from %s to %s", a.World.iterations, a.Id, DirectionNames[pd], DirectionNames[a.Direction])
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
						a.World.Debugf("[%04d][%s] (%d) consuming %s (%d)", a.World.iterations, a.Id, a.Points, o.Id, o.Points)

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

	a.Reset()

	return a
}

func (a *Agent) Load(program string) error {
	if code, err := gopush.ParseCode(program); err != nil {
		return err
	} else {
		a.code = code
	}

	return nil
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
	a.Points = 100

	for _, s := range a.interp.Stacks {
		s.Flush()
	}

	return a
}

func (a *Agent) Step() {
	a.interp.RunCode(a.code)
}

func (a *Agent) Mutate(generation int) *Agent {
	n := NewAgent(a.World, generation)

	n.code = compactCode(mutateCode(a.interp, copyCode(a.code)))

	return n
}

func mutateCode(interp *gopush.Interpreter, c gopush.Code) gopush.Code {
	for i := 0; i < len(c.List); i++ {
		c.List[i] = mutateCode(interp, c.List[i])
	}

	if rand.Intn(100) > 95 {
		return interp.RandomCode(5)
	}

	return c
}

func compactCode(c gopush.Code) gopush.Code {
	for len(c.List) == 1 {
		c = c.List[0]
	}

	for i := 0; i < len(c.List); i++ {
		c.List[i] = compactCode(c.List[i])
	}

	return c
}

func copyCode(c gopush.Code) gopush.Code {
	if len(c.List) == 0 {
		return c
	}

	list := make([]gopush.Code, len(c.List))

	for i, v := range c.List {
		list[i] = copyCode(v)
	}

	return c
}
