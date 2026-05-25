package bddtest

// bdd_test.go is the test runner that wires godog to Go's testing framework.
//
// Run S3 scenarios:   go test -v -run TestBDDS3 -timeout 30m
// Run VPC scenarios:  go test -v -run TestBDDVPC -timeout 30m
// Run all:            go test -v -timeout 30m

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/cucumber/godog"
)

func TestMain(m *testing.M) {
	profile := os.Getenv("AWS_PROFILE")
	if profile == "" {
		fmt.Fprintln(os.Stderr, "ERROR: AWS_PROFILE is not set.")
		os.Exit(1)
	}
	if err := os.Setenv("AWS_PROFILE", profile); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to set AWS_PROFILE: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

// TestBDDS3 runs the S3 bucket BDD feature specs.
func TestBDDS3(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeS3Scenario,
		Options: &godog.Options{
			Format:         "pretty",
			Paths:          []string{"features/s3_bucket.feature"},
			TestingT:       t,
			DefaultContext: WithTestingT(context.Background(), t),
		},
	}

	if suite.Run() != 0 {
		t.Fatal("BDD S3 tests failed")
	}
}

// TestBDDVPC runs the VPC BDD feature specs.
func TestBDDVPC(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeVPCScenario,
		Options: &godog.Options{
			Format:         "pretty",
			Paths:          []string{"features/vpc.feature"},
			TestingT:       t,
			DefaultContext: WithTestingT(context.Background(), t),
		},
	}

	if suite.Run() != 0 {
		t.Fatal("BDD VPC tests failed")
	}
}
