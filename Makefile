all: preygent

preygent: main.go agent.go world.go
	go build
