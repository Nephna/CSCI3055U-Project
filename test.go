package main

import (
	"fmt"
)

type thing struct {
	aNum int

	func something () (int) {
		return aNum
	}
}

func main () {
	aThing := thing{10}
	fmt.Println(aThing.something())
}