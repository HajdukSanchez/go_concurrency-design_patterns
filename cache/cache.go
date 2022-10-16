package main

import (
	"fmt"
	"sync"
	"time"
)

func ExpensiveFibonacci(n int) int {
	fmt.Printf("Calculate expensive Fibonacci for %d \n", n)
	time.Sleep(5 * time.Second)
	return n
}

type Service struct {
	InProgress map[int]bool       // Key: Number to calculate Fibonacci, Value: if there is in progress
	IsPending  map[int][]chan int // key: Number to calculate Fibonacci, Value: Workers waiting to current process calculate
	mtx        sync.RWMutex
}

func (s *Service) Work(job int) {
	s.mtx.RLock()
	exists := s.InProgress[job]
	if exists {
		// If exists, someone takes this job and it is processing
		s.mtx.RUnlock()
		response := make(chan int)
		defer close(response)

		s.mtx.Lock()
		s.IsPending[job] = append(s.IsPending[job], response) // Add the created channel
		s.mtx.Unlock()
		fmt.Printf("Waiting for response job: %d\n", job)
		resp := <-response
		fmt.Printf("Response done, received: %d\n", resp)
		return
	}
	s.mtx.RUnlock()

	s.mtx.RLock()
	s.InProgress[job] = true
	s.mtx.RUnlock()
	fmt.Printf("Calculate Fibonacci for: %d\n", job)
	result := ExpensiveFibonacci(job)

	s.mtx.RLock()
	pendingWorkers, exists := s.IsPending[job]
	s.mtx.RUnlock()

	if exists {
		for _, pendingWorker := range pendingWorkers {
			pendingWorker <- result
			fmt.Printf("Result sent - all pending workers ready job: %d\n", job)
		}
	}

	s.mtx.Lock()
	s.InProgress[job] = false
	s.IsPending[job] = make([]chan int, 0) // Slice with no values
	s.mtx.Unlock()
}

func NewService() *Service {
	return &Service{
		InProgress: make(map[int]bool),
		IsPending:  make(map[int][]chan int),
	}
}

func main() {
	service := NewService()
	jobs := []int{3, 4, 5, 5, 4, 8, 8, 8}
	var wg sync.WaitGroup
	wg.Add(len(jobs))
	for _, job := range jobs {
		go func(n int) {
			defer wg.Done()
			service.Work(n)
		}(job)
	}
	wg.Wait()
}
