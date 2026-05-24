package retrybackoff_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/davidcediel12/go-engineering/retrybackoff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDo(t *testing.T) {
	type testCase struct {
		backoff                     retrybackoff.Backoff
		argCtx                      context.Context
		argOperation                func() error
		expectedErrorIs             error
		attempts                    func() int
		expectedAttempts            int
		expectedTimesBetweenRetries []time.Duration
		timeBetweenAttempts         func() []time.Duration
	}

	setupTests := map[string]func(t *testing.T) *testCase{
		"should_not_retry_when_operation_succeed": func(t *testing.T) *testCase {
			t.Helper()
			attempts := 0
			operation := func() error {
				attempts += 1
				return nil
			}

			backoff := retrybackoff.New()
			return &testCase{
				argOperation:     operation,
				argCtx:           context.Background(),
				backoff:          *backoff,
				expectedErrorIs:  nil,
				attempts:         func() int { return attempts },
				expectedAttempts: 1,
			}
		},
		"should_retry_once_and_then_succeed_if_operation_fails": func(t *testing.T) *testCase {
			t.Helper()
			attempts := 0
			timeBetweenAttempts := make([]time.Duration, 1)
			var lastExecutionTime time.Time
			operation := func() error {
				defer func() {
					attempts += 1
					lastExecutionTime = time.Now()
				}()
				if attempts > 0 {
					timeBetweenAttempts[attempts-1] = time.Since(lastExecutionTime)
				}
				if attempts < 1 {
					return errors.New("exploded")
				}
				return nil
			}
			return &testCase{
				argOperation:                operation,
				argCtx:                      context.Background(),
				backoff:                     *retrybackoff.New(),
				expectedErrorIs:             nil,
				attempts:                    func() int { return attempts },
				expectedAttempts:            2,
				timeBetweenAttempts:         func() []time.Duration { return timeBetweenAttempts },
				expectedTimesBetweenRetries: []time.Duration{time.Duration(retrybackoff.DefaultBaseDelay)},
			}
		},
		"should_exhaust_default_retries_and_fail_when_operation_fails_indefenitely": func(t *testing.T) *testCase {
			t.Helper()
			expectedError := errors.New("exploded")

			backoff := retrybackoff.New()
			attempts := 0
			timeBetweenAttempts := make([]time.Duration, backoff.Retries())
			var lastExecutionTime time.Time
			operation := func() error {
				defer func() {
					lastExecutionTime = time.Now()
				}()

				if attempts > 0 {
					timeBetweenAttempts[attempts-1] = time.Since(lastExecutionTime)
				}
				attempts += 1
				return expectedError

			}
			return &testCase{
				argOperation:        operation,
				argCtx:              context.Background(),
				backoff:             *backoff,
				expectedErrorIs:     expectedError,
				attempts:            func() int { return attempts },
				expectedAttempts:    retrybackoff.DefaultRetries + 1, // initial attempt  + retries
				timeBetweenAttempts: func() []time.Duration { return timeBetweenAttempts },
				expectedTimesBetweenRetries: []time.Duration{
					backoff.BaseDelay(),
					minDuration(2*backoff.BaseDelay(), backoff.MaxDelay()),
					minDuration(4*backoff.BaseDelay(), backoff.MaxDelay()),
				},
			}
		},
		"should_exhaust_custom_retries_and_fail_when_operation_fails_indefenitely": func(t *testing.T) *testCase {
			t.Helper()
			expectedError := errors.New("exploded")

			backoff := retrybackoff.New(retrybackoff.WithRetries(6))
			timeBetweenAttempts := make([]time.Duration, backoff.Retries())
			var lastExecutionTime time.Time
			attempts := 0
			operation := func() error {
				defer func() {
					lastExecutionTime = time.Now()
				}()

				if attempts > 0 {
					timeBetweenAttempts[attempts-1] = time.Since(lastExecutionTime)
				}
				attempts += 1
				return expectedError

			}
			return &testCase{
				argOperation:        operation,
				argCtx:              context.Background(),
				backoff:             *backoff,
				expectedErrorIs:     expectedError,
				attempts:            func() int { return attempts },
				expectedAttempts:    7, // initial attempt  + retries
				timeBetweenAttempts: func() []time.Duration { return timeBetweenAttempts },
				expectedTimesBetweenRetries: []time.Duration{
					backoff.BaseDelay(),
					minDuration(2*backoff.BaseDelay(), backoff.MaxDelay()),
					minDuration(4*backoff.BaseDelay(), backoff.MaxDelay()),
					minDuration(8*backoff.BaseDelay(), backoff.MaxDelay()),
					minDuration(16*backoff.BaseDelay(), backoff.MaxDelay()),
					minDuration(32*backoff.BaseDelay(), backoff.MaxDelay()),
				},
			}
		},
	}

	for name, setupTest := range setupTests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			test := setupTest(t)

			err := test.backoff.Do(test.argCtx, test.argOperation)
			assert.ErrorIs(t, err, test.expectedErrorIs)
			assert.Equal(t, test.expectedAttempts, test.attempts())

			if test.timeBetweenAttempts == nil {
				return
			}
			timeBetweenAttempts := test.timeBetweenAttempts()
			for i := range test.expectedTimesBetweenRetries {
				fmt.Printf("expected time: %v, actual: %v, delta: %v\n",
					test.expectedTimesBetweenRetries[i], timeBetweenAttempts[i], test.backoff.BaseDelay())

				require.GreaterOrEqual(t,
					timeBetweenAttempts[i],
					test.expectedTimesBetweenRetries[i])

				require.LessOrEqual(t,
					timeBetweenAttempts[i],
					test.expectedTimesBetweenRetries[i]+test.backoff.BaseDelay()+5*time.Millisecond)
			}

		})
	}
}

func minDuration(d1, d2 time.Duration) time.Duration {
	if d1 < d2 {
		return d1
	}
	return d2
}
