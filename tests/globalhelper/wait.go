package globalhelper

import (
	"context"
	"errors"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

// WaitForState polls fetch until it reports the awaited state, returns an error, or timeout elapses.
// fetch returns (false, nil) while the resource is still converging; returning a non-nil error stops
// polling immediately and that error is returned to the caller.
func WaitForState(ctx context.Context, what string, interval, timeout time.Duration,
	fetch wait.ConditionWithContextFunc,
) error {
	err := wait.PollUntilContextTimeout(ctx, interval, timeout, true, fetch)
	if err != nil && errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("timed out after %v waiting for %s", timeout, what)
	}

	return err
}
