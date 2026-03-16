package bddtest

// helpers.go provides shared utilities for BDD step definitions.

import (
	"context"
	"testing"
)

type testingTKey struct{}

// WithTestingT stores *testing.T in a context (used in DefaultContext).
func WithTestingT(ctx context.Context, t *testing.T) context.Context {
	return context.WithValue(ctx, testingTKey{}, t)
}

// GetTestingT retrieves *testing.T from the context.
// Requires DefaultContext to be set via WithTestingT in the test suite (see bdd_test.go).
func GetTestingT(ctx context.Context) *testing.T {
	t, ok := ctx.Value(testingTKey{}).(*testing.T)
	if !ok {
		panic("*testing.T not found in context — set DefaultContext via WithTestingT in bdd_test.go")
	}
	return t
}
