package test

import (
	"fmt"
	"os"
	"testing"

	sdkaws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
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

func terraformOptions() *terraform.Options {
	return &terraform.Options{
		TerraformDir: "../",
	}
}

// TestS3BucketCreation deploys the S3 bucket module, verifies the bucket
// exists in AWS, then tears everything down.
func TestS3BucketCreation(t *testing.T) {
	t.Parallel()

	opts := terraformOptions()
	defer terraform.Destroy(t, opts)
	terraform.InitAndApply(t, opts)

	bucketName := terraform.Output(t, opts, "bucket_name")
	aws.AssertS3BucketExists(t, awsRegion, bucketName)
}

// TestS3BucketHasEncryption verifies the S3 bucket has server-side
// encryption enabled (AES256).
func TestS3BucketHasEncryption(t *testing.T) {
	t.Parallel()

	opts := terraformOptions()
	defer terraform.Destroy(t, opts)
	terraform.InitAndApply(t, opts)

	bucketName := terraform.Output(t, opts, "bucket_name")

	// Verify encryption is enabled (GetS3BucketEncryption removed in terratest v0.47+)
	sess, err := aws.NewAuthenticatedSession(awsRegion)
	if err != nil {
		t.Fatalf("Failed to create AWS session: %s", err)
	}
	s3Client := s3.New(sess)
	result, err := s3Client.GetBucketEncryption(&s3.GetBucketEncryptionInput{
		Bucket: sdkaws.String(bucketName),
	})
	if err != nil {
		t.Fatalf("Failed to get bucket encryption: %s", err)
	}
	rules := result.ServerSideEncryptionConfiguration.Rules
	assert.Equal(t, 1, len(rules))
	assert.Equal(t, "AES256", *rules[0].ApplyServerSideEncryptionByDefault.SSEAlgorithm)
}

// TestS3BucketHasVersioning verifies the S3 bucket has versioning enabled.
func TestS3BucketHasVersioning(t *testing.T) {
	t.Parallel()

	opts := terraformOptions()
	defer terraform.Destroy(t, opts)
	terraform.InitAndApply(t, opts)

	bucketName := terraform.Output(t, opts, "bucket_name")

	versioning := aws.GetS3BucketVersioning(t, awsRegion, bucketName)
	assert.Equal(t, "Enabled", versioning)
}
