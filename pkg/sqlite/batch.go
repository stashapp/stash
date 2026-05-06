package sqlite

const defaultBatchSize = 1000

// batchExec executes the provided function in batches of the provided size.
func batchExec[T any](ids []T, batchSize int, fn func(batch []T) error) error {
	for i := 0; i < len(ids); i += batchSize {
		end := i + batchSize
		if end > len(ids) {
			end = len(ids)
		}

		batch := ids[i:end]
		if err := fn(batch); err != nil {
			return err
		}
	}

	return nil
}
