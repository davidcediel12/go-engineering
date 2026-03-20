package quitchannel

import (
	"fmt"

	"golang.org/x/tour/tree"
)

// walkRecursive walks the tree t sending all values
// from the tree to the channel ch.
func walkRecursive(t *tree.Tree, ch, quit chan int) {
	if t == nil {
		return
	}
	walkRecursive(t.Left, ch, quit)
	select {
	case ch <- t.Value: // Sent value
	case _, ok := <-quit:
		if !ok { // signal to stop (quit was closed)
			return
		}
	}
	walkRecursive(t.Right, ch, quit)
}

// Walk is a wrapper that starts the recursion to walks the tree
// it closes the channel to avoid goroutine leaks when the recursion is done
func Walk(t *tree.Tree, ch, quit chan int) {
	defer close(ch)
	walkRecursive(t, ch, quit)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
// It uses a quit channel to notify the producers to stop when returning early
func Same(t1, t2 *tree.Tree) bool {
	t1Chan := make(chan int)
	t2Chan := make(chan int)
	quit := make(chan int)
	defer close(quit) // Close the quit channel to notify producers to stop

	go Walk(t1, t1Chan, quit)
	go Walk(t2, t2Chan, quit)

	for {
		t1Val, ok1 := <-t1Chan // Blocks until receiving value from t1
		t2Val, ok2 := <-t2Chan // Blocks until receiving value from t2

		if !ok1 && !ok2 {
			return true
		}
		if !ok1 || !ok2 {
			return false
		}
		if t1Val != t2Val {
			return false
		}
	}
}

func MainQuit() {
	ch := make(chan int)

	go Walk(tree.New(1), ch, make(chan int)) // In Walk(), we don't need to quit early, so we send a channel without any control
	for v := range ch {
		fmt.Println(v)
	}

	fmt.Println(Same(tree.New(1), tree.New(1)))
}
