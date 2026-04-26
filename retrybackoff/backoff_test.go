package retrybackoff_test

import (
	"context"
	"testing"

	"github.com/davidcediel12/go-engineering/retrybackoff"
	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	type testCase struct {
		backoff       retrybackoff.Backoff
		argCtx        context.Context
		argOperation  func() error
		expectedError error
	}

	setupTests := map[string]func(t *testing.T) *testCase{
		"should_not_retry_when_operation_succeed": func(t *testing.T) *testCase {

			retries := 0
			operation := func() error {
				if retries > 0 {
					t.Errorf("expected no retries")
				}
				retries++
				return nil
			}
			return &testCase{
				argOperation:  operation,
				argCtx:        context.Background(),
				backoff:       *retrybackoff.New(),
				expectedError: nil,
			}
		},
	}

	for name, setupTest := range setupTests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			test := setupTest(t)

			err := test.backoff.Do(test.argCtx, test.argOperation)
			assert.Equal(t, test.expectedError, err)
		})
	}
}
