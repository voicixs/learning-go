package main

import (
	"fmt"
	"math"
)

// ============================================================================
// 1. BASIC INTERFACE & IMPLICIT IMPLEMENTATION
// ============================================================================
// In Go, an interface is a "contract" that defines method signatures without
// containing any implementation logic.
//
// Visibility rules apply here too:
// - Uppercase first letter (e.g., Shape) = Public / Exported
// - Lowercase first letter = Private / Package-internal

type Shape interface {
	Area() float64
}

// Struct 1: Rectangle
type Rectangle struct {
	Width  float64
	Height float64
}

// Struct 2: Circle
type Circle struct {
	Radius float64
}

// Implicit Implementation:
// There is NO "implements" keyword in Go. A struct implements an interface
// simply by implementing all the methods declared in the interface contract.

// Rectangle implements Shape by defining Area():
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Circle implements Shape by defining Area():
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// ============================================================================
// 2. POLYMORPHISM
// ============================================================================
// This generic function accepts ANY type that satisfies the Shape contract.
func CalculateArea(s Shape) float64 {
	return s.Area()
}

// ============================================================================
// 3. INTERFACE COMPOSITION / EMBEDDING
// ============================================================================
// Interfaces can embed other interfaces to combine multiple contracts.

type Measurable interface {
	Perimeter() float64
}

// Geometry embeds both Shape and Measurable.
// Any type satisfying Geometry must implement both Area() and Perimeter().
type Geometry interface {
	Shape
	Measurable
}

// Let's implement Perimeter() on Rectangle so it satisfies Measurable and Geometry:
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// A function requiring the combined Geometry interface:
func DescribeShape(g Geometry) {
	fmt.Printf("Area: %.2f\n", g.Area())
	fmt.Printf("Perimeter: %.2f\n", g.Perimeter())
}

// ============================================================================
// 4. THE EMPTY INTERFACE & TYPE ASSERTIONS
// ============================================================================
// An empty interface `interface{}` has zero methods, meaning EVERY type
// in Go satisfies it. It can hold values of any type.
// Note: In Go 1.18+, `any` is an alias for `interface{}`.

func DescribeValue(t interface{}) {
	// %T prints the underlying concrete type, %v prints the value
	fmt.Printf("Type: %T, Value: %v\n", t, t)
}

// ============================================================================
// 5. CUSTOM ERRORS VIA THE BUILT-IN ERROR INTERFACE
// ============================================================================
// In Go, `error` is simply a built-in interface defined as:
//
//	type error interface {
//	    Error() string
//	}
//
// By implementing an `Error() string` method on a custom struct, that struct
// automatically becomes a valid Go error!

type CalculationError struct {
	Message string
}

func (c CalculationError) Error() string {
	return c.Message
}

// Example function returning our custom error:
func PerformCalculation(val float64) (float64, error) {
	if val < 0 {
		return 0, CalculationError{Message: "invalid input: value must be non-negative"}
	}
	return math.Sqrt(val), nil
}

// ============================================================================
// 6. IDIOMATIC ADDITION: TYPE SWITCH
// ============================================================================
// While single type assertions (v, ok := val.(int)) work for one type,
// a "type switch" is the idiomatic Go way to handle multiple possible types.
func IdentifyType(val any) {
	switch v := val.(type) {
	case int:
		fmt.Printf("It's an integer: %d\n", v)
	case string:
		fmt.Printf("It's a string: %s\n", v)
	case bool:
		fmt.Printf("It's a boolean: %t\n", v)
	default:
		fmt.Printf("Unknown type: %T\n", v)
	}
}

// ============================================================================
// MAIN EXECUTION
// ============================================================================

func main() {
	fmt.Println("--- 1. Polymorphism with Interfaces ---")
	rect := Rectangle{Width: 5, Height: 4}
	circle := Circle{Radius: 2}

	// CalculateArea accepts both structs because both implement Shape:
	fmt.Printf("Rectangle Area: %.2f\n", CalculateArea(rect))
	fmt.Printf("Circle Area: %.2f\n", CalculateArea(circle))

	fmt.Println("\n--- 2. Interface Composition (Embedding) ---")
	// Rectangle implements both Area() and Perimeter(), so it satisfies Geometry:
	DescribeShape(rect)

	fmt.Println("\n--- 3. Empty Interface & Type Assertions ---")
	var mysteryBox interface{} = 10
	DescribeValue(mysteryBox)

	// Safe Type Assertion: value, ok := interfaceValue.(TargetType)
	retrievedInt, ok := mysteryBox.(int)
	if ok {
		fmt.Printf("Successfully retrieved int: %d\n", retrievedInt)
	} else {
		fmt.Println("Value is not an integer")
	}

	// Failed assertion check:
	mysteryBox = "hello"
	DescribeValue(mysteryBox)
	_, ok = mysteryBox.(int)
	if !ok {
		fmt.Println("Checked mysteryBox as int: Value is not an integer")
	}

	fmt.Println("\n--- 4. Type Switch (Idiomatic Addition) ---")
	IdentifyType(42)
	IdentifyType("Gopher")
	IdentifyType(true)

	fmt.Println("\n--- 5. Custom Error Implementation ---")
	result, err := PerformCalculation(16)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Calculation result: %.2f\n", result)
	}

	// Triggering the error condition:
	_, err = PerformCalculation(-5)
	if err != nil {
		fmt.Printf("Caught expected error: %v\n", err)
	}
}
