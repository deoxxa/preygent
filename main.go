package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	w := NewWorld(true)
	g := 0

	var topAgent *Agent

	for {
		w.Reset()

		g++

		for i := 0; i < 20; i++ {
			w.AddAgent(NewAgent(w, g).Reset())
		}
		if topAgent != nil {
			w.AddAgent(topAgent)
		}

		for w.Step() {
		}

		sort.Sort(AgentList(w.agents))
		if topAgent != w.agents[0] {
			topAgent = w.agents[0]

			fmt.Printf("the world ran for %d ticks in generation %d\n", w.ticks, g)
			fmt.Printf("the new top agent is %s (from generation %d) with %d points\n", topAgent.Id, topAgent.Generation, topAgent.Points)
			fmt.Printf("\n")

			fmt.Printf("%s\n", topAgent.code)
			fmt.Printf("\n")

			if len(w.lines) > 0 {
				for _, l := range w.lines {
					fmt.Println(l)
				}

				fmt.Printf("\n")
			}
		}
	}
}
