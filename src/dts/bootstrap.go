package main

import (
	"fmt"
	"github.com/DaigangLi/godts/src/dts/zk"
)

func main() {
	str := "Hello, world"
	fmt.Println(str)
	zk.Connect()
}
