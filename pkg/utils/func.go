package utils

// Do executes each function in the slice in order. If any function returns an error, it is returned immediately.
func Do(fn []func() error) error {
	for _, f := range fn {
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}
