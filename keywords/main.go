package main

import (
	"fmt"
	// 5c. Blank import: Registers driver/initialization side-effects without direct symbol usage
	_ "database/sql"
)

// ============================================================================
// CONCEPT 1: THE IOTA KEYWORD & BITWISE FLAGS
// ============================================================================

// Standard automatic numbering starting at 0:
const (
	Monday = iota // 0
	Tuesday       // 1
	Wednesday     // 2
	Thursday      // 3
	Friday        // 4
	Saturday      // 5
	Sunday        // 6
)

// Starting at 1 using arithmetic expressions:
const (
	Jan = iota + 1 // 1
	Feb            // 2
	Mar            // 3
)

// Bit flags using bitwise shifts (1 << iota):
// Commonly used for UNIX file permissions or boolean setting masks.
const (
	Readable   = 1 << iota // 1 << 0 = 1 (binary 001)
	Writable               // 1 << 1 = 2 (binary 010)
	Executable             // 1 << 2 = 4 (binary 100)
)

// ============================================================================
// CONCEPT 2: THE `new` KEYWORD
// ============================================================================
// new(T) allocates zero-initialized memory and returns a pointer (*T).

type Counter struct {
	Count int
	Name  string
}

func (c *Counter) Increment() {
	c.Count++
}

// Constructor returning a pointer using new():
func NewCounter() *Counter {
	// Allocates memory where Count is initialized to 0 and Name to ""
	return new(Counter)
}

// ============================================================================
// CONCEPT 3: LABELED BREAK IN NESTED LOOPS
// ============================================================================
// A normal 'break' only exits the innermost loop.
// A labeled break exits the specified outer loop directly.

func demonstrateLabeledBreak() {
	fmt.Println("Starting nested loop search:")

OuterLoop:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 1 && j == 1 {
				fmt.Printf("Breaking directly out of OuterLoop at i=%d, j=%d\n", i, j)
				break OuterLoop
			}
			fmt.Printf("Processing i=%d, j=%d\n", i, j)
		}
	}
	fmt.Println("Exited outer loop successfully.")
}

// ============================================================================
// CONCEPT 4: THE `goto` KEYWORD & LABELS
// ============================================================================
// goto jumps execution directly to a label.
// Useful for state machines or low-level performance-critical loops.

func demonstrateGoto() {
	i := 0

StartLoop:
	if i >= 3 {
		goto Finished
	}

	fmt.Printf("goto loop iteration: %d\n", i)
	i++
	goto StartLoop

Finished:
	fmt.Println("Reached Finished label via goto.")
}

// ============================================================================
// CONCEPT 5: THE BLANK IDENTIFIER (`_`)
// ============================================================================

func sampleFunction() ([]int, error) {
	return []int{10, 20, 30}, nil
}

func demonstrateBlankIdentifier() {
	// 5a. Discarding unwanted return values:
	numbers, _ := sampleFunction()

	// 5b. Discarding slice indices in 'for range' loops:
	fmt.Print("Iterating slice values only: ")
	for _, val := range numbers {
		fmt.Printf("%d ", val)
	}
	fmt.Println()
}

// ============================================================================
// MAIN EXECUTION
// ============================================================================

func main() {
	fmt.Println("--- 1. iota & Bit Flags ---")
	fmt.Printf("Monday: %d, Sunday: %d\n", Monday, Sunday)
	fmt.Printf("Jan: %d, Feb: %d\n", Jan, Feb)
	fmt.Printf("Readable: %b (int %d)\n", Readable, Readable)
	fmt.Printf("Writable: %b (int %d)\n", Writable, Writable)
	fmt.Printf("Executable: %b (int %d)\n", Executable, Executable)

	fmt.Println("\n--- 2. new() Keyword ---")
	c := NewCounter()
	c.Increment()
	fmt.Printf("Counter pointer: %p, Count: %d, Name: %q\n", c, c.Count, c.Name)

	fmt.Println("\n--- 3. Labeled Break ---")
	demonstrateLabeledBreak()

	fmt.Println("\n--- 4. goto Keyword ---")
	demonstrateGoto()

	fmt.Println("\n--- 5. Blank Identifier ---")
	demonstrateBlankIdentifier()
}
// ```

// ---

// ### Key Concepts Breakdown

// 1. **The `iota` Keyword & Bitwise Flags**:
//    * Eliminates manual constant sequencing by automatically generating sequential integers starting at 0.
//    * Can be combined with arithmetic expressions (`iota + 1`) or bit-shifting (`1 << iota`) to generate bitmask flags like `Readable`, `Writable`, and `Executable`.

// 2. **The `new` Keyword**:
//    * Allocates zero-initialized memory for a given type and returns a pointer (`*T`) to it.
//    * Numeric fields default to `0`, strings default to `""`, and references default to `nil`.

// 3. **Labeled `break`**:
//    * A standard `break` statement only terminates the innermost loop.
//    * Attaching a label (e.g., `OuterLoop:`) to the outer loop allows you to break out of all nested loops simultaneously without needing extra boolean flags.

// 4. **The `goto` Keyword**:
//    * Jumps program execution directly to a defined label.
//    * While rarely needed in high-level application code, it is useful for low-level state machines, cleanup handlers, or performance-critical loops.

// 5. **The Blank Identifier (`_`)**:
//    * Represents an explicitly ignored value in Go.
//    * Can be used to discard unused return values or ignore index positions in `for _, val := range slice`.
//    * When placed before an import (e.g., `_ "database/sql"`), it executes the package's `init()` function for side effects without importing any exported symbols directly.
