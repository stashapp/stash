package models

import "errors"

var ErrScraperSource = errors.New("invalid ScraperSource")

type ScrapedItemReader interface {
	All() ([]*ScrapedItem, error)
}

type ScrapedItemWriter interface {
	Create(newObject ScrapedItem) (*ScrapedItem, error)
}

type ScrapedItemReaderWriter interface {
	ScrapedItemReader
	ScrapedItemWriter
}
