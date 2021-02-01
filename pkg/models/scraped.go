package models

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
