package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

// Save a function reference and cache result for this function
type Memory struct {
	function Function
	cache    map[int]FunctionResult
	mtx      sync.Mutex // To avoid race condition
}

// Function to map a function with his int number
type Function func(key int) (interface{}, error)

// Needs to be the same type as Function return
type FunctionResult struct {
	value interface{} // Handle dynamic type of value (int, string, bool, etc)
	err   error
}

// Constructor for Memory struct
func NewCache(f Function) *Memory {
	return &Memory{
		function: f,
		cache:    make(map[int]FunctionResult),
	}
}

// Get value saved on cache
func (m *Memory) GetValue(key int) (interface{}, error) {
	m.mtx.Lock() // Block cache saved process for other subroutines
	result, exists := m.cache[key]
	m.mtx.Unlock()
	// If value is not found, calculate and saved it
	if !exists {
		m.mtx.Lock()
		result.value, result.err = m.function(key)
		m.cache[key] = result
		m.mtx.Unlock()
	}
	return result.value, result.err
}

// Same type as FUNCTION type defined above (to match)
func GetFibonacci(n int) (interface{}, error) {
	return Fibonacci(n), nil // Suppose there are no errors
}

func main() {
	cache := NewCache(GetFibonacci)                          // Create new cache fot Fibonacci function
	fib := []int{42, 40, 41, 42, 38, 42, 42, 38, 42, 42, 42} // Numbers to calculate
	var wg sync.WaitGroup
	channel := make(chan int, 2) // To limit the amount of subroutines calculating the same number
	for _, n := range fib {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			channel <- 1 // Start using channel
			start := time.Now()
			value, err := cache.GetValue(index) // Get value saved on cache based on Key
			if err != nil {
				log.Println(err)
			}
			fmt.Printf("Number: %d, Time taken: %s, Result: %d\n", index, time.Since(start), value)
			<-channel // Make the channel available
		}(n)
	}
	wg.Wait()
}
