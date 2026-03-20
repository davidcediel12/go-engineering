package main

import (
	"os"

	"github.com/davidcediel12/go-engineering/workerpool/binarytree/contextcancellation"
	"github.com/davidcediel12/go-engineering/workerpool/binarytree/quitchannel"
)

func main() {
	option := os.Args[1]

	if option == "quit" {
		quitchannel.MainQuit()
	}

	if option == "ctx" {
		contextcancellation.MainCtx()
	}
}
