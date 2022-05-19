package models

type Element struct {
	ID int
}

func (e *Element) Mismatched() []Element {
	return []Element{*e}
}
