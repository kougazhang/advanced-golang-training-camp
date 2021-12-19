package main

import "fmt"

func main() {
	ben := &Ben{name: "Ben"}
	jerry := &Jerry{name: "Jerry"}
	var maker IceCreamMaker = ben

	var loop0, loop1 func()

	loop0 = func() {
		maker = ben
		go loop1()
	}

	loop1 = func() {
		maker = jerry
		go loop0()
	}

	go loop0()

	for {
		maker.Hello()
	}
}

type IceCreamMaker interface {
	Hello()
}

type Ben struct {
	id   int
	name string
}

func (b *Ben) Hello() {
	fmt.Printf("Ben: Hello, My Name is %s\n", b.name)
}

type Jerry struct {
	name string
}

func (j *Jerry) Hello() {
	fmt.Printf("Jerry: Hello, My Name is %s\n", j.name)
}
