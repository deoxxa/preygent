package main

import (
	"fmt"
	"time"
)

type World struct {
	Debug bool

	agents []*Agent
	active bool
	ticks  int
}

func NewWorld(debug bool) *World {
	return &World{
		Debug: debug,
	}
}

func (w *World) AddAgent(a *Agent) {
	w.agents = append(w.agents, a)
}

func (w *World) Reset() {
	w.agents = []*Agent{}
	w.ticks = 0
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

	w.ticks++

	if w.Debug {
		fmt.Printf("\n")
		time.Sleep(time.Second)
	}

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
