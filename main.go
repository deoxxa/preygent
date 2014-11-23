package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	w := NewWorld(true)
	g := 0

	bestAgents := []*Agent{}

	if program, err := ioutil.ReadFile("best.push"); err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
	} else {
		bestAgent := NewAgent(w, 0)
		if err := bestAgent.Load(string(program)); err != nil {
			panic(err)
		}
		bestAgents = append(bestAgents, bestAgent)
	}

	for {
		g++

		// fmt.Printf("running generation %d\n", g)

		w.Reset()

		if len(bestAgents) > 0 {
			for _, a := range bestAgents {
				w.AddAgent(a)
				w.AddAgent(a.Mutate(g))
				w.AddAgent(a.Mutate(g))
			}
		}
		for i := 0; i < 250; i++ {
			w.AddAgent(NewAgent(w, g).Reset())
		}

		for w.Step() {
		}

		sort.Sort(AgentList(w.agents))

		found := false
		for _, a := range bestAgents {
			if a == w.agents[0] {
				found = true
			}
		}

		if !found {
			fmt.Printf("the world ran for %d iterations in generation %d\n", w.iterations, g)
			fmt.Printf("the current best is %s (from generation %d) with %d points\n", w.agents[0].Id, w.agents[0].Generation, w.agents[0].Points)
			fmt.Printf("\n")

			fmt.Printf("%s\n", w.agents[0].code)
			fmt.Printf("\n")

			// if len(w.lines) > 0 {
			// 	for _, l := range w.lines {
			// 		fmt.Println(l)
			// 	}

			// 	fmt.Printf("\n")
			// }

			if err := ioutil.WriteFile("best.push", []byte(w.agents[0].code.String()), 0644); err != nil {
				panic(err)
			}
		}

		bestAgents = w.agents[0:25]
	}
}
