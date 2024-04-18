package gateway

import (
	"context"
	"errors"
	"testing"
	"time"
)

var errSomethingWrong = errors.New("something went wrong")

func Test_wrapper(t *testing.T) {
	// Test case when function f returns an error
	t.Run("function error", func(t *testing.T) {
		err := wrapper(context.Background(), func() error {
			return errSomethingWrong
		})

		if !errors.Is(err, errSomethingWrong) {
			t.Errorf("expected %v, got %v", errSomethingWrong, err)
		}
	})

	// Test case when context is done
	t.Run("context done", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := wrapper(ctx, func() error {
			time.Sleep(100 * time.Millisecond)
			return nil
		})

		if !errors.Is(err, context.Canceled) {
			t.Errorf("expected %v, got %v", context.Canceled, err)
		}
	})

	// Test case when function f completes before context is done
	t.Run("function completes", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := wrapper(ctx, func() error {
			time.Sleep(50 * time.Millisecond)
			return nil
		})

		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})
}
