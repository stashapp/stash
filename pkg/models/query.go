package models

type QueryOptions struct {
	FindFilter *FindFilterType
	Count      bool
}

type QueryResult[T comparable] struct {
	IDs   []T
	Count int
}
