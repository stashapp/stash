//+build !debug

package plugins

func log(_ string, fn func() error) error {
	return fn()
}
