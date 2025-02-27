package main

import (
	"fmt"
	"strings"
)

func main() {
	uploadFile(strings.Split(SelectFile(), "file://")[1])
}

func uploadFile(file string) {
	fmt.Println(file)
}
