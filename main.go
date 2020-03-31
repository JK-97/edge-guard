package main

//go:generate go run version_generate.go

import (
	"github.com/JK-97/edge-guard/cmd"
)

func main() {
	cmd.Execute()

}
