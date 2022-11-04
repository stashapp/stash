package astits

type wrappingCounter struct {
	value  int
	wrapAt int
}

func newWrappingCounter(wrapAt int) wrappingCounter {
	return wrappingCounter{
		value:  wrapAt + 1,
		wrapAt: wrapAt,
	}
}

func (c *wrappingCounter) get() int {
	return c.value
}

func (c *wrappingCounter) inc() int {
	c.value++
	if c.value > c.wrapAt {
		c.value = 0
	}
	return c.value
}
