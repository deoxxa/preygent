package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	w := NewWorld(false)

	var topAgent *Agent

	for {
		w.Reset()

		for i := 0; i < 20; i++ {
			w.AddAgent(NewAgent(w).Reset())
		}
		if topAgent != nil {
			w.AddAgent(topAgent)
		}

		for w.Step() {
		}

		sort.Sort(AgentList(w.agents))
		if topAgent != w.agents[0] {
			topAgent = w.agents[0]
			fmt.Printf("the world ran for %d ticks\n", w.ticks)
			fmt.Printf("the new top agent is %s with %d points\n", topAgent.Id, topAgent.Points)
			fmt.Printf("\n")
			fmt.Printf("%s\n", topAgent.code)
			fmt.Printf("\n")
		}
	}
}
