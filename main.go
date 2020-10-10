package main

import (
	"autobuy/pkg/autobuy"
	"fmt"
)

func main() {
	result := autobuy.Check()
	fmt.Println(result)
}
