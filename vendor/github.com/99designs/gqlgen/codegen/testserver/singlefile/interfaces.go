package singlefile

import "math"

type Shape interface {
	Area() float64
	isShape()
}

type ShapeUnion interface {
	Area() float64
	isShapeUnion()
}

type Circle struct {
	Radius float64
	Coordinates
}

func (c *Circle) Area() float64 {
	return c.Radius * math.Pi * math.Pi
}

func (c *Circle) isShapeUnion() {}
func (c *Circle) isShape()      {}

type Rectangle struct {
	Length, Width float64
	Coordinates
}

func (r *Rectangle) Area() float64 {
	return r.Length * r.Width
}
func (r *Rectangle) isShapeUnion() {}
func (r *Rectangle) isShape()      {}

type Node interface {
	Child() (Node, error)
}

type ConcreteNodeA struct {
	ID    string
	Name  string
	child Node
}

func (n *ConcreteNodeA) Child() (Node, error) {
	return n.child, nil
}

//  Implements the Node interface with another interface
type ConcreteNodeInterface interface {
	Node
	ID() string
}

type ConcreteNodeInterfaceImplementor struct{}

func (c ConcreteNodeInterfaceImplementor) ID() string {
	return "CNII"
}

func (c ConcreteNodeInterfaceImplementor) Child() (Node, error) {
	return &ConcreteNodeA{
		ID:   "Child",
		Name: "child",
	}, nil
}

type BackedByInterface interface {
	ThisShouldBind() string
	ThisShouldBindWithError() (string, error)
}

type BackedByInterfaceImpl struct {
	Value string
	Error error
}

func (b *BackedByInterfaceImpl) ThisShouldBind() string {
	return b.Value
}

func (b *BackedByInterfaceImpl) ThisShouldBindWithError() (string, error) {
	return b.Value, b.Error
}
