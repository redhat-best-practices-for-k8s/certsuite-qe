package globalhelper

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWaitForState(t *testing.T) {
	t.Run("returns as soon as fetch reports the awaited state", func(t *testing.T) {
		calls := 0

		err := WaitForState(context.TODO(), "a resource", time.Millisecond, time.Second,
			func(context.Context) (bool, error) {
				calls++

				return true, nil
			})
		require.NoError(t, err)
		assert.Equal(t, 1, calls)
	})

	t.Run("retries while fetch reports the resource is still converging", func(t *testing.T) {
		calls := 0

		err := WaitForState(context.TODO(), "a resource", time.Millisecond, time.Second,
			func(context.Context) (bool, error) {
				calls++

				return calls == 3, nil
			})
		require.NoError(t, err)
		assert.Equal(t, 3, calls)
	})

	t.Run("returns fetch's error when polling cannot make progress", func(t *testing.T) {
		injected := errors.New("not schedulable")

		err := WaitForState(context.TODO(), "a resource", time.Millisecond, time.Second,
			func(context.Context) (bool, error) {
				return false, injected
			})
		assert.ErrorIs(t, err, injected)
	})

	t.Run("times out when the state never arrives", func(t *testing.T) {
		err := WaitForState(context.TODO(), "a resource", time.Millisecond, 5*time.Millisecond,
			func(context.Context) (bool, error) {
				return false, nil
			})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "timed out after")
		assert.Contains(t, err.Error(), "a resource")
	})

	t.Run("returns the context error when the caller cancels", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := WaitForState(ctx, "a resource", time.Millisecond, time.Second,
			func(context.Context) (bool, error) {
				return false, nil
			})
		assert.ErrorIs(t, err, context.Canceled)
	})
}
