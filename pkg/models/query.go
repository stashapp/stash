package models

type QueryOptions struct {
	FindFilter *FindFilterType
	Count      bool
}

type QueryResult struct {
	IDs   []int
	Count int
}
