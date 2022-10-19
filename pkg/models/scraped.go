package models

import (
	"context"
	"errors"
)

var ErrScraperSource = errors.New("invalid ScraperSource")

type ScrapedItemReader interface {
	All(ctx context.Context) ([]*ScrapedItem, error)
}

type ScrapedItemWriter interface {
	Create(ctx context.Context, newObject ScrapedItem) (*ScrapedItem, error)
}

type ScrapedItemReaderWriter interface {
	ScrapedItemReader
	ScrapedItemWriter
}
