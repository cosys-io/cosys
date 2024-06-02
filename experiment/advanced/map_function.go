package main

import "fmt"

func f(p string) {
	fmt.Println("function f parameter:", p)
}

func g(p string, q int) {
	fmt.Println("function g parameters:", p, q)
}

func main() {
	m := map[string]interface{}{
		"f": f,
		"g": g,
	}
	for k, v := range m {
		switch k {
		case "f":
			v.(func(string))("astring")
		case "g":
			v.(func(string, int))("astring", 42)
		}
	}
}
