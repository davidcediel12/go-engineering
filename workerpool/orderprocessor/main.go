package main

import (
	"os"
	"strconv"

	orderprocessor "github.com/davidcediel12/go-engineering/workerpool/orderprocessor/queue"
)

func main() {
	workers, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic("invalid number of workers")
	}
	queue := orderprocessor.New(int64(workers))
	queue.Start()
}
