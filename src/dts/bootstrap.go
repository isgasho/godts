package main

import (
	"dts/zk"
	"fmt"
)

func main() {
	str := "Hello, world"
	fmt.Println(str)
	zk.Connect()
}
