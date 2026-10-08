package main

import (
	"fmt"
	"math/rand"
	"time"
)

// ============================================================================
// 1. BASIC ENUMS WITH IOTA
// ============================================================================
// Go does not have a built-in 'enum' keyword. Instead, we group constants
// using the 'const' block and the 'iota' keyword.
// 'iota' is an auto-incrementing integer index that resets to 0 in each const block.

const (
	Pending  = iota // 0
	Approved        // 1 (automatically increments)
	Refused         // 2
	Unknown         // 3
)

// You can also apply expressions to iota (e.g., bitshifts or multiplication):
const (
	StepZero = iota * 2 // 0 * 2 = 0
	StepOne             // 1 * 2 = 2
	StepTwo             // 2 * 2 = 4
)

// ============================================================================
// 2. STRING-BASED ENUMS VIA DEFINED TYPES
// ============================================================================
// You can define a custom type based on 'string' to restrict function inputs
// and avoid raw string typos across your application.

type Protocol string

const (
	TCP   Protocol = "TCP"
	UDP   Protocol = "UDP"
	HTTP  Protocol = "HTTP"
	HTTPS Protocol = "HTTPS"
)

// ============================================================================
// 3. INTEGER ENUMS WITH CUSTOM STRING CONVERSION (fmt.Stringer)
// ============================================================================
// Defining a custom type on top of 'int' provides compile-time type safety.
type PacketStatus int

const (
	Accepted PacketStatus = iota // 0
	Rejected                     // 1
	Suspicious                   // 2
)

// By implementing String() string, PacketStatus satisfies the fmt.Stringer
// interface. When printed with fmt.Println or %s, Go automatically calls
// this method instead of printing the underlying raw integer (0, 1, 2).
func (ps PacketStatus) String() string {
	names := [...]string{"Accepted", "Rejected", "Suspicious"}

	// Guard against unmapped integer values
	if int(ps) < 0 || int(ps) >= len(names) {
		return "UnknownStatus"
	}
	return names[ps]
}

// ============================================================================
// 4. REAL-WORLD TUTORIAL EXAMPLE: PACKET ANALYZER
// ============================================================================

// Packet represents a simulated network packet.
type Packet struct {
	ID            int
	SourceIP      string
	DestinationIP string
	Size          int
	Protocol      Protocol
}

// PacketAnalyzer aggregates network metrics using our enums.
type PacketAnalyzer struct {
	TotalPackets      int
	AcceptedPackets   int
	RejectedPackets   int
	SuspiciousPackets int
}

// Analyze inspects packet attributes and assigns a status enum.
func (pa *PacketAnalyzer) Analyze(p Packet) PacketStatus {
	pa.TotalPackets++

	// 1. Mark unencrypted HTTP packets as Suspicious
	if p.Protocol == HTTP {
		pa.SuspiciousPackets++
		return Suspicious
	}

	// 2. Reject packets larger than 1500 bytes (typical MTU limit)
	if p.Size > 1500 {
		pa.RejectedPackets++
		return Rejected
	}

	// 3. Accept all other packets
	pa.AcceptedPackets++
	return Accepted
}

// Helper function to generate simulated packets
func generateRandomPacket(id int, rng *rand.Rand) Packet {
	protocols := []Protocol{TCP, UDP, HTTP, HTTPS}

	return Packet{
		ID:            id,
		SourceIP:      fmt.Sprintf("192.168.1.%d", rng.Intn(254)+1),
		DestinationIP: fmt.Sprintf("10.0.0.%d", rng.Intn(254)+1),
		Size:          rng.Intn(1800) + 100, // Size between 100 and 1900 bytes
		Protocol:      protocols[rng.Intn(len(protocols))],
	}
}

// ============================================================================
// MAIN EXECUTION
// ============================================================================

func main() {
	fmt.Println("--- 1. Basic Iota Enums ---")
	fmt.Printf("Pending: %d, Approved: %d, Refused: %d\n", Pending, Approved, Refused)
	fmt.Printf("StepZero: %d, StepOne: %d, StepTwo: %d\n", StepZero, StepOne, StepTwo)

	fmt.Println("\n--- 2. Stringer Interface in Action ---")
	status := Suspicious
	// Notice that %d prints the raw int, while %s or Println uses the String() method:
	fmt.Printf("Underlying integer: %d\n", status)
	fmt.Printf("Stringer representation: %s\n", status)

	fmt.Println("\n--- 3. Running the Packet Analyzer Simulation ---")
	analyzer := &PacketAnalyzer{}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Analyze 10 sample packets
	for i := 1; i <= 10; i++ {
		packet := generateRandomPacket(i, rng)
		resultStatus := analyzer.Analyze(packet)

		fmt.Printf("Packet #%02d | Proto: %-5s | Size: %4d bytes | Status: %s\n",
			packet.ID, packet.Protocol, packet.Size, resultStatus)
	}

	fmt.Println("\n--- Analyzer Summary Metrics ---")
	fmt.Printf("Total Packets Analyzed : %d\n", analyzer.TotalPackets)
	fmt.Printf("Accepted Packets       : %d\n", analyzer.AcceptedPackets)
	fmt.Printf("Rejected Packets (>1.5KB): %d\n", analyzer.RejectedPackets)
	fmt.Printf("Suspicious Packets (HTTP): %d\n", analyzer.SuspiciousPackets)
}
// ```

// ---

// ### Core Concepts Explained

// 1. **Why Avoid Magic Numbers?**  
//    Raw integers such as `0`, `1`, or `2` provide no context to developers reading your code. Giving them named constants like `Accepted`, `Rejected`, or `Suspicious` communicates clear intent and prevents typos.
// 2. **The `iota` Counter**  
//    The `iota` identifier represents successive integer constants in a `const` group. It starts at index `0` and increments by 1 for each line, allowing you to generate sequential values without hardcoding integers.
// 3. **Defined Types for Type Safety**  
//    Declaring `type PacketStatus int` or `type Protocol string` prevents developers from accidentally passing an arbitrary integer or string to a function that requires a specific enum.
// 4. **Implementing `fmt.Stringer`**  
//    By adding a `String() string` method to the custom enum type, functions like `fmt.Println` and formatting verbs like `%s` output readable textual labels rather than raw numbers.
