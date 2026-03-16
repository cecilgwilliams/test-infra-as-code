package test

// TDD Demo — Step 1: Write the test FIRST.
//
// AWS_PROFILE=your-profile go test -v -run TestS3BucketHasVersioning -timeout 30m   # FAIL

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

const awsRegion = "us-east-1"

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

func TestS3BucketHasVersioning(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../",
	}

	// Always clean up
	defer terraform.Destroy(t, terraformOptions)

	// Deploy
	terraform.InitAndApply(t, terraformOptions)

	bucketName := terraform.Output(t, terraformOptions, "bucket_name")

	versioning := aws.GetS3BucketVersioning(t, awsRegion, bucketName)
	assert.Equal(t, "Enabled", versioning)
}
