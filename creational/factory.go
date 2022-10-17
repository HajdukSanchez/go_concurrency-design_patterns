package main

import "fmt"

// Getters ans setters
type Product interface {
	getStock() int
	setStock(stock int)
	getName() string
	setName(name string)
}

/// New kind of product

type Computer struct {
	name  string
	stock int
}

// ToString function
func (c *Computer) String() string {
	s := fmt.Sprintf("Product: %s, with Stock: %d", c.name, c.stock)
	return s
}

func (c *Computer) getStock() int {
	return c.stock
}

func (c *Computer) setStock(stock int) {
	c.stock = stock
}

func (c *Computer) getName() string {
	return c.name
}

func (c *Computer) setName(name string) {
	c.name = name
}

/// New kind of computer

type Laptop struct {
	Computer
}

func NewLaptop() Product {
	return &Laptop{
		Computer: Computer{
			name:  "Laptop computer",
			stock: 25,
		},
	}
}

/// New kind of computer

type Desktop struct {
	Computer
}

func NewDesktop() Product {
	return &Desktop{
		Computer: Computer{
			name:  "Desktop computer",
			stock: 35,
		},
	}
}

// This is the factory
func ComputerFactory(computerType Product) (Product, error) {
	switch computerType.(type) {
	case *Desktop:
		return NewDesktop(), nil
	case *Laptop:
		return NewLaptop(), nil
	default:
		return nil, fmt.Errorf("Unknown computer type: %T", computerType)
	}
}

func main() {
	laptop, _ := ComputerFactory(&Desktop{})
	desktop, _ := ComputerFactory(&Laptop{})

	fmt.Println(laptop)
	fmt.Println(desktop)
}
