package main

import (
	"fmt"
)

type World struct {
	Debug bool

	agents     []*Agent
	active     bool
	iterations int
	lines      []string
}

func NewWorld(debug bool) *World {
	return &World{
		Debug: debug,
	}
}

func (w *World) Debugf(format string, v ...interface{}) {
	w.lines = append(w.lines, fmt.Sprintf(format, v...))
}

func (w *World) AddAgent(a *Agent) {
	w.agents = append(w.agents, a)
}

func (w *World) Reset() {
	w.agents = []*Agent{}
	w.iterations = 0
	w.lines = []string{}
}

func (w *World) Step() bool {
	w.active = false

	for _, a := range w.agents {
		// this agent is dead
		if a.Points == 0 {
			continue
		}

		a.Step()
	}

	w.iterations++

	return w.active
}

func (w *World) AgentsAt(x, y int) []*Agent {
	found := []*Agent{}

	for _, a := range w.agents {
		if a.X == x && a.Y == y {
			found = append(found, a)
		}
	}

	return found
}
