package main

import (
	"encoding/json"
	"fmt"
)

// ============================================================================
// 1. DEFINING STRUCTS & FIELD VISIBILITY
// ============================================================================
// A struct groups related data together into a single custom data type.
//
// In Go, casing controls visibility:
// - Uppercase first letter = Exported / Public (accessible outside the package)
// - Lowercase first letter = Unexported / Private (only accessible within the package)

type Address struct {
	// Struct tags (using backticks) provide metadata, commonly used for
	// JSON serialization, XML, or database ORMs.
	Street string `json:"my_street"`
	City   string `json:"my_city"`
}

type Employee struct {
	// Fields with uppercase letters are visible to external packages like `encoding/json`.
	Name     string `json:"name"`
	Age      int    `json:"age"`
	IsRemote bool   `json:"remote"` // Custom key name in JSON output

	// 5. EMBEDDED STRUCTS / COMPOSITION
	// Go uses composition instead of inheritance.
	// By embedding Address without giving it a field name, its fields and
	// methods are "promoted" directly to Employee.
	Address
}

// ============================================================================
// 2. METHODS: VALUE RECEIVER VS. POINTER RECEIVER
// ============================================================================

// Method on Address: Because Address is embedded into Employee, Employee
// can call this method directly!
func (a Address) PrintAddress() {
	fmt.Printf("Address: %s, %s\n", a.Street, a.City)
}

// VALUE RECEIVER: (e Employee)
// Operates on a COPY of the struct.
// Modifying fields here only affects the local copy inside this function,
// NOT the original struct instance outside.
func (e Employee) UpdateNameWrong(newName string) {
	e.Name = newName
	fmt.Printf("Inside UpdateNameWrong: name is %s (local copy only)\n", e.Name)
}

// POINTER RECEIVER: (e *Employee)
// Receives a pointer (memory address) to the original struct.
// Any modification directly mutates the original struct.
// Note: Use pointer receivers when mutating data or when structs are large
// to avoid expensive memory copying.
func (e *Employee) UpdateName(newName string) {
	e.Name = newName
}

// ============================================================================
// 3. IDIOMATIC ADDITION: CONSTRUCTOR FUNCTIONS (FACTORY PATTERN)
// ============================================================================
// Go does not have built-in class constructors like Java or C++.
// A common Go idiom is writing a `New...` function to initialize structs
// with default values or validations.
func NewEmployee(name string, age int, isRemote bool, street, city string) *Employee {
	return &Employee{
		Name:     name,
		Age:      age,
		IsRemote: isRemote,
		Address: Address{
			Street: street,
			City:   city,
		},
	}
}

// ============================================================================
// MAIN EXECUTION
// ============================================================================

func main() {
	fmt.Println("--- 1. Struct Initialization & Access ---")

	// Initializing a struct using named field syntax
	emp1 := Employee{
		Name:     "Alice",
		Age:      30,
		IsRemote: true,
		Address: Address{
			Street: "123 Main St",
			City:   "New York",
		},
	}

	// Accessing struct fields using dot notation
	fmt.Printf("Employee Name: %s, Age: %d, Remote: %t\n", emp1.Name, emp1.Age, emp1.IsRemote)

	// Due to struct embedding (composition), we can access Address fields directly:
	fmt.Printf("Street: %s, City: %s\n", emp1.Street, emp1.City)

	// We can also call the embedded method directly on emp1:
	emp1.PrintAddress()

	fmt.Println("\n--- 2. Anonymous Structs ---")
	// Used for one-time data structures without creating a global named type.
	// Great for quick encapsulation inside a function or test payloads.
	job := struct {
		Title  string
		Salary int
	}{
		Title:  "Software Engineer",
		Salary: 100000,
	}

	fmt.Printf("Job: %s, Salary: $%d\n", job.Title, job.Salary)

	fmt.Println("\n--- 3. Value Receiver vs Pointer Receiver ---")

	// Calling the value receiver method does NOT change emp1.Name
	emp1.UpdateNameWrong("Bob")
	fmt.Printf("After UpdateNameWrong: Name is still %s\n", emp1.Name) // Still "Alice"

	// Calling the pointer receiver method updates the original struct
	emp1.UpdateName("Bob")
	fmt.Printf("After UpdateName: Name is now %s\n", emp1.Name) // Updated to "Bob"

	// Direct pointer manipulation:
	empPtr := &emp1
	empPtr.Age = 31 // Go automatically dereferences the pointer to access the field
	fmt.Printf("Updated Age via pointer: %d\n", emp1.Age)

	fmt.Println("\n--- 4. Struct Tags & JSON Serialization ---")
	// Serialization: Convert struct to JSON using json.Marshal
	// Notice that the exported field names map to the lowercase/custom tags in the output.
	jsonData, err := json.Marshal(emp1)
	if err != nil {
		fmt.Printf("Error marshalling JSON: %v\n", err)
		return
	}

	// Convert byte slice to readable string
	fmt.Printf("Serialized JSON: %s\n", string(jsonData))

	fmt.Println("\n--- 5. Idiomatic Additions: Zero Values & Constructors ---")
	// Zero values: In Go, uninitialized struct fields default to their zero values
	// (strings = "", numbers = 0, bools = false, pointers/slices/maps = nil).
	var emptyEmp Employee
	fmt.Printf("Zero Value Employee: Name='%s', Age=%d, IsRemote=%t\n", emptyEmp.Name, emptyEmp.Age, emptyEmp.IsRemote)

	// Using the constructor function:
	emp2 := NewEmployee("Charlie", 28, false, "456 Elm St", "San Francisco")
	fmt.Printf("Created via Constructor: %+v\n", emp2)
}
