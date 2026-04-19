# Go Concurrency Patterns: Worker Pool / Binary Tree

A minimal example of two concurrency patterns for tree traversal with **early cancellation** – comparing two binary trees without walking them fully when a mismatch is found.

## What you get

- **Walk** – sends tree values (in‑order) to a channel.
- **Same** – walks two trees concurrently, stops at the first difference.

## Two patterns, two packages

| Package | Cancellation mechanism | Entry point |
|---------|------------------------|--------------|
| `quitchannel` | Dedicated `quit chan int` (closed on mismatch) | `MainQuit()` |
| `contextcancellation` | `context.Context` (canceled on mismatch) | `MainCtx()` |

## Running the examples
The project provides a single entry point that selects the implementation via command line:

```bash
go run main.go quit   # uses quitchannel.MainQuit()
go run main.go ctx    # uses contextcancellation.MainCtx()
```

Each MainQuit / MainCtx creates two trees, compares them, and prints the result.

## How cancellation works

**Quit channel** – `Same` creates a quit channel. If a value mismatch occurs or one walk finishes early, it closes quit. Walkers check quit before sending and exit if closed.

**Context** – `SameCtx` creates a cancellable context. On mismatch, it calls the cancel function. Walkers listen to ctx.Done() and stop the recursion.

Both patterns close the output channel after the walk completes to avoid goroutine leaks.