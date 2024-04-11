package server

import "fmt"

type Command struct {
	// add stuff
}

func parseCommand(msg string) (Command, error) {
	t := msg[0]
	fmt.Println(t)
	switch t {
	case '*':
		fmt.Println(msg[1:])
	}
	return Command{}, nil
}
