package main

import "math/rand"

type fish int

const (
	onefish fish = iota
	twofish
	redfish
	bluefish
)

func drianPond() fish {
	return fish(rand.Intn(3))
}

func main() {

	p := drianPond()

	// lingo:exhaustive_switch
	switch p {
	case onefish:
	default:
	}

	// lingo:exhaustive_switch
	switch p {
	case onefish, twofish:
	default:
	}

	// lingo:exhaustive_switch
	switch p {
	case redfish:
	case bluefish:
	default:
	}

	// lingo:exhaustive_switch
	switch p {
	case onefish, twofish:
	case redfish:
	default:
	}

	// lingo:exhaustive_switch
	switch p {
	case onefish:
	default:
	}

	// should be ignored
	switch p {
	case onefish:
	default:
	}

}
