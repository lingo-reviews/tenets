package main

import "fmt"

func main() {
	saySomething("hi")
}

func saySomething(something string) {
	fmt.Println("hi")
}

func saySomethingAgain(something string) {
	fmt.Println(something)
}

func saySomethingOther(something, otherthing string) {
	fmt.Println("hi")
}

func saySomethingElse(something, otherthing string) {
	fmt.Println(something)
}
