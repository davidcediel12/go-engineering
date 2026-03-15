package binarytree

import (
	"context"
	"fmt"

	"golang.org/x/tour/tree"
)

// walkRecursiveCtx walks the tree t sending all values
// from the tree to the channel ch, stopping if the context was cancelled
func walkRecursiveCtx(ctx context.Context, t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}
	walkRecursiveCtx(ctx, t.Left, ch)
	select {
	case ch <- t.Value: // Sent value
	case <-ctx.Done(): // signal to stop (context was cancelled)
		return
	}
	walkRecursiveCtx(ctx, t.Right, ch)
}

// WalkCtx is a facade to start recursion to walks the tree,
// it makes sure that once the recursion is done, the channel is closed
func WalkCtx(ctx context.Context, t *tree.Tree, ch chan int) {
	defer close(ch)
	walkRecursiveCtx(ctx, t, ch)
}

// SameCtx determines whether the trees
// t1 and t2 contain the same values, returning early and closing the context if it
// finds a mismatch
func SameCtx(t1, t2 *tree.Tree) bool {
	t1Chan := make(chan int)
	t2Chan := make(chan int)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go WalkCtx(ctx, t1, t1Chan)
	go WalkCtx(ctx, t2, t2Chan)

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

func mainCtx() {
	ch := make(chan int)
	go WalkCtx(context.Background(), tree.New(1), ch) // In Walk(), we don't need to quit early, so we send a background context
	for v := range ch {
		fmt.Println(v)
	}

	fmt.Println(SameCtx(tree.New(1), tree.New(1)))
}
