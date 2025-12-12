package main

import "fmt"

func main() {
	for k, v := range all {
		fmt.Println(k, v)
		if k == "foo" {
			break
		}
	}
}

func all(field func(string, any) bool) {
	_ = field("hello", "world") &&
		field("foo", "bar") &&
		field("baz", "done")
}
