package astits

type wrappingCounter struct {
	wrapAt int
	value  int
}

func newWrappingCounter(wrapAt int) wrappingCounter {
	return wrappingCounter{
		wrapAt: wrapAt,
	}
}

// returns current counter state and increments internal value
func (c *wrappingCounter) get() int {
	ret := c.value
	c.value++
	if c.value > c.wrapAt {
		c.value = 0
	}
	return ret
}
