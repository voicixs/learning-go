package main

import (
	"fmt"
	"sync"
	"time"
)

// ============================================================================
// 1. BASIC GOROUTINES & SYNC.WAITGROUP
// ============================================================================
// A goroutine is a lightweight thread managed directly by the Go runtime.
// You start one simply by placing the 'go' keyword before any function call.
//
// Because the 'main' function is also a goroutine, if 'main' exits, all child
// goroutines are instantly terminated. We use sync.WaitGroup to wait for them.

func printWorker(id int, wg *sync.WaitGroup) {
	// defer wg.Done() guarantees that when this function exits, it decrements
	// the WaitGroup counter, signaling that this worker is finished.
	defer wg.Done()

	fmt.Printf("[Worker %d] Starting work...\n", id)
	time.Sleep(100 * time.Millisecond) // Simulating work
	fmt.Printf("[Worker %d] Finished work!\n", id)
}

// ============================================================================
// 2. UNBUFFERED CHANNELS (SYNCHRONOUS COMMUNICATION)
// ============================================================================
// Channels are typed conduits through which goroutines send and receive values.
// Rule: By default, sending blocks until a receiver is ready, and receiving
// blocks until a sender is ready. This synchronizes goroutines without locks.

func calculateSquare(num int, out chan int) {
	// Send result into channel using the arrow operator (ch <- value)
	out <- num * num
}

// ============================================================================
// 3. BUFFERED CHANNELS, RANGE, AND CLOSE
// ============================================================================
// A buffered channel has capacity: make(chan T, capacity).
// Sending only blocks when the buffer is full; receiving only blocks when empty.
//
// When a producer finishes sending, it calls close(ch).
// Receivers can loop over incoming values using 'for value := range ch'.

func produceNumbers(count int, ch chan int) {
	for i := 1; i <= count; i++ {
		ch <- i
	}
	// Only the sender should close a channel. Closing signals no more items.
	close(ch)
}

// ============================================================================
// 4. THE SELECT STATEMENT (MULTIPLEXING CHANNELS)
// ============================================================================
// 'select' lets a goroutine wait on multiple channel operations simultaneously.
// It executes whichever case is ready first.

func fastWorker(ch chan string) {
	time.Sleep(50 * time.Millisecond)
	ch <- "fast response"
}

func slowWorker(ch chan string) {
	time.Sleep(200 * time.Millisecond)
	ch <- "slow response"
}

// ============================================================================
// MAIN EXECUTION
// ============================================================================

func main() {
	fmt.Println("--- 1. Goroutines and WaitGroups ---")

	var wg sync.WaitGroup

	// Launch 3 concurrent workers
	for i := 1; i <= 3; i++ {
		wg.Add(1) // Tell WaitGroup to expect 1 more worker
		go printWorker(i, &wg)
	}

	// wg.Wait() blocks until the WaitGroup counter reaches zero
	wg.Wait()
	fmt.Println("All workers completed!")

	fmt.Println("\n--- 2. Unbuffered Channels ---")

	// Create an unbuffered channel of integers
	results := make(chan int)

	// Launch computation concurrently
	go calculateSquare(9, results)

	// Receive from channel (blocks until calculateSquare sends a value)
	squaredValue := <-results
	fmt.Printf("Received from channel: 9 squared is %d\n", squaredValue)

	fmt.Println("\n--- 3. Buffered Channels & Range/Close ---")

	// Create a channel that can buffer up to 3 values without blocking
	bufferedCh := make(chan int, 3)

	go produceNumbers(3, bufferedCh)

	// 'range' continuously pulls values until the channel is closed
	for val := range bufferedCh {
		fmt.Printf("Received from buffered channel: %d\n", val)
	}

	fmt.Println("\n--- 4. Select Statement (Channel Multiplexing) ---")

	ch1 := make(chan string)
	ch2 := make(chan string)

	go fastWorker(ch1)
	go slowWorker(ch2)

	// Wait for whichever channel delivers data first
	select {
	case msg1 := <-ch1:
		fmt.Printf("select received: %s\n", msg1)
	case msg2 := <-ch2:
		fmt.Printf("select received: %s\n", msg2)
	case <-time.After(300 * time.Millisecond):
		fmt.Println("select timed out")
	}
}
// ```

// ---

// Here's what I found on **Go Concurrency (Goroutines and Channels)**:

// I found that Go approaches concurrency using lightweight goroutines managed directly by the runtime, relying on channels for message passing rather than shared-memory locking ("share memory by communicating").

// **Key themes I noticed:**
// 1. **Goroutines vs. OS Threads**: Goroutines start with very small stacks (a few kilobytes) and scale dynamically, allowing programs to run thousands concurrently.
// 2. **Coordination with `sync.WaitGroup`**: Because the `main` routine does not wait for spawned goroutines by default, `sync.WaitGroup` is standard for coordinating batch tasks.
// 3. **Synchronization with Channels**: Unbuffered channels act as rendezvous points where senders and receivers block until both are ready.
// 4. **Multiplexing with `select`**: The `select` statement enables listening to multiple channels at once, providing timeouts and non-blocking reads.

// 📥 The full research report is included as the first source in the import card below — import it to chat with the findings directly.