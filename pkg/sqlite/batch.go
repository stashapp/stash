package sqlite

const defaultBatchSize = 1000

// batchExec executes the provided function in batches of the provided size.
func batchExec(ids []int, batchSize int, fn func(batch []int) error) error {
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
