package health

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestMultiChecker_Health(t *testing.T) {
	type tcase struct {
		checkers []Checker
		want     error
	}
	tests := map[string]tcase{
		"Nil checkers": {
			checkers: nil,
			want:     nil,
		},
		"Nil error": {
			checkers: []Checker{
				&testChecker{HealthFunc: func(ctx context.Context) error { return nil }},
				&testChecker{HealthFunc: func(ctx context.Context) error { return nil }},
			},
			want: nil,
		},
		"Non nil error": {
			checkers: []Checker{
				&testChecker{HealthFunc: func(ctx context.Context) error { return errTest }},
				&testChecker{HealthFunc: func(ctx context.Context) error { return nil }},
			},
			want: errTest,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			c := NewMultiChecker(tc.checkers...)
			if err := c.Health(context.Background()); err != tc.want {
				t.Errorf("Health() error = %v, want %v", err, tc.want)
			}
		})
	}
}

func TestMultiChecker_HealthDone(t *testing.T) {
	c := NewMultiChecker(&testChecker{
		HealthFunc: func(ctx context.Context) error { return nil }})
	waitExit := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		_ = c.Health(ctx)
		waitExit <- struct{}{}
	}()
	cancel()
	select {
	case <-time.After(time.Second * 1):
		t.Error("timepout")
	case <-waitExit:
	}
}

var errTest = errors.New("test error")

type testChecker struct {
	HealthFunc func(ctx context.Context) error
}

func (c *testChecker) Health(ctx context.Context) error {
	if c.HealthFunc == nil {
		return nil
	}
	return c.HealthFunc(ctx)
}
