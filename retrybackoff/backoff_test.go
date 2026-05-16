package retrybackoff_test

import (
	"context"
	"errors"
	"testing"

	"github.com/davidcediel12/go-engineering/retrybackoff"
	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	type testCase struct {
		backoff          retrybackoff.Backoff
		argCtx           context.Context
		argOperation     func() error
		expectedErrorIs  error
		attempts         func() int
		expectedAttempts int
	}

	setupTests := map[string]func(t *testing.T) *testCase{
		"should_not_retry_when_operation_succeed": func(t *testing.T) *testCase {
			t.Helper()
			attempts := 0
			operation := func() error {
				attempts += 1
				return nil
			}
			return &testCase{
				argOperation:     operation,
				argCtx:           context.Background(),
				backoff:          *retrybackoff.New(),
				expectedErrorIs:  nil,
				attempts:         func() int { return attempts },
				expectedAttempts: 1,
			}
		},
		"should_retry_once_and_then_succeed_if_operation_fails": func(t *testing.T) *testCase {
			t.Helper()
			attempts := 0
			operation := func() error {
				defer func() { attempts += 1 }()
				if attempts < 1 {
					return errors.New("exploded")
				}
				return nil
			}
			return &testCase{
				argOperation:     operation,
				argCtx:           context.Background(),
				backoff:          *retrybackoff.New(),
				expectedErrorIs:  nil,
				attempts:         func() int { return attempts },
				expectedAttempts: 2,
			}
		},
		"should_exhaust_default_retries_and_fail_when_operation_fails_indefenitely": func(t *testing.T) *testCase {
			t.Helper()
			expectedError := errors.New("exploded")

			attempts := 0
			operation := func() error {
				attempts += 1
				return expectedError

			}
			return &testCase{
				argOperation:     operation,
				argCtx:           context.Background(),
				backoff:          *retrybackoff.New(),
				expectedErrorIs:  expectedError,
				attempts:         func() int { return attempts },
				expectedAttempts: retrybackoff.DefaultRetries + 1, // initial attempt  + retries
			}
		},
		"should_exhaust_custom_retries_and_fail_when_operation_fails_indefenitely": func(t *testing.T) *testCase {
			t.Helper()
			expectedError := errors.New("exploded")

			attempts := 0
			operation := func() error {
				attempts += 1
				return expectedError

			}
			return &testCase{
				argOperation:     operation,
				argCtx:           context.Background(),
				backoff:          *retrybackoff.New(retrybackoff.WithRetries(6)),
				expectedErrorIs:  expectedError,
				attempts:         func() int { return attempts },
				expectedAttempts: 7, // initial attempt  + retries
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
		})
	}
}
