package health

import (
	"context"
	"sync"
)

// Checker represents health check logic.
type Checker interface {
	// Health performs health check and return error if it fails.
	Health(ctx context.Context) error
}

// NewMultiChecker returns a new instance of MultiChecker.
// Takes variadic length of Checker interfaces.
func NewMultiChecker(checks ...Checker) *MultiChecker {
	hc := &MultiChecker{
		checks: make([]Checker, 0, len(checks)),
	}
	hc.checks = append(hc.checks, checks...)
	return hc
}

// MultiChecker holds an array of checkers end error.
// Struct implements Checker interface.
type MultiChecker struct {
	checks []Checker
	wg     sync.WaitGroup
	so     sync.Once
	err    error
}

func (c *MultiChecker) Health(ctx context.Context) error {
	for _, check := range c.checks {
		c.wg.Add(1)

		go func(ctx context.Context, f func(ctx context.Context) error) {
			defer c.wg.Done()
			select {
			case <-ctx.Done():
				c.so.Do(func() {
					c.err = ctx.Err()
				})
				return
			default:
				if err := f(ctx); err != nil {
					c.so.Do(func() {
						c.err = err
					})
				}
			}
		}(ctx, check.Health)
	}
	c.wg.Wait()
	return c.err
}
