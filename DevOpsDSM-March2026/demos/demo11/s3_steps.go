package bddtest

// S3 BDD step definitions using godog + Terratest.
//
// These steps implement the Gherkin scenarios in features/s3_bucket.feature.
// They use Terratest to deploy real AWS resources and make assertions.

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cucumber/godog"
	terraaws "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

const awsRegion = "us-east-1"

// S3TestContext holds shared state between steps for S3 scenarios.
type S3TestContext struct {
	TerraformOptions *terraform.Options
	BucketName       string
}

// NewS3TestContext creates a new test context pointing at the S3 Terraform module.
func NewS3TestContext() *S3TestContext {
	return &S3TestContext{
		TerraformOptions: &terraform.Options{
			TerraformDir: "../demo8",
		},
	}
}

func (s *S3TestContext) theTerraformModuleHasBeenApplied(ctx context.Context) error {
	t := GetTestingT(ctx)
	terraform.InitAndApply(t, s.TerraformOptions)
	s.BucketName = terraform.Output(t, s.TerraformOptions, "bucket_name")
	return nil
}

func (s *S3TestContext) theS3BucketShouldExistIn(ctx context.Context, region string) error {
	t := GetTestingT(ctx)
	terraaws.AssertS3BucketExists(t, region, s.BucketName)
	return nil
}

func (s *S3TestContext) theS3BucketShouldHaveEncryption(ctx context.Context, algo string) error {
	t := GetTestingT(ctx)
	sess, err := terraaws.NewAuthenticatedSession(awsRegion)
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %s", err)
	}
	s3Client := s3.New(sess)
	out, err := s3Client.GetBucketEncryption(&s3.GetBucketEncryptionInput{
		Bucket: aws.String(s.BucketName),
	})
	if err != nil {
		return fmt.Errorf("failed to get bucket encryption: %s", err)
	}
	rules := out.ServerSideEncryptionConfiguration.Rules
	assert.Equal(t, 1, len(rules))
	actual := *rules[0].ApplyServerSideEncryptionByDefault.SSEAlgorithm
	assert.Equal(t, algo, actual)
	return nil
}

func (s *S3TestContext) theS3BucketShouldHaveVersioningEnabled(ctx context.Context) error {
	t := GetTestingT(ctx)
	versioning := terraaws.GetS3BucketVersioning(t, awsRegion, s.BucketName)
	if versioning != "Enabled" {
		return fmt.Errorf("expected versioning Enabled, got %s", versioning)
	}
	return nil
}

func (s *S3TestContext) theS3BucketShouldBlockPublicAccess(ctx context.Context) error {
	t := GetTestingT(ctx)
	sess, err := terraaws.NewAuthenticatedSession(awsRegion)
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %s", err)
	}
	s3Client := s3.New(sess)
	out, err := s3Client.GetPublicAccessBlock(&s3.GetPublicAccessBlockInput{
		Bucket: aws.String(s.BucketName),
	})
	if err != nil {
		return fmt.Errorf("failed to get public access block: %s", err)
	}
	cfg := out.PublicAccessBlockConfiguration
	assert.True(t, *cfg.BlockPublicAcls)
	assert.True(t, *cfg.IgnorePublicAcls)
	assert.True(t, *cfg.BlockPublicPolicy)
	assert.True(t, *cfg.RestrictPublicBuckets)
	return nil
}

// Cleanup destroys the Terraform resources.
func (s *S3TestContext) Cleanup(ctx context.Context) {
	t := GetTestingT(ctx)
	terraform.Destroy(t, s.TerraformOptions)
}

// InitializeS3Scenario registers all S3 step definitions with godog.
func InitializeS3Scenario(ctx *godog.ScenarioContext) {
	s := NewS3TestContext()

	ctx.Given(`^the Terraform module has been applied$`, s.theTerraformModuleHasBeenApplied)
	ctx.Then(`^the S3 bucket should exist in "([^"]*)"$`, s.theS3BucketShouldExistIn)
	ctx.Then(`^the S3 bucket should have "([^"]*)" encryption$`, s.theS3BucketShouldHaveEncryption)
	ctx.Then(`^the S3 bucket should have versioning enabled$`, s.theS3BucketShouldHaveVersioningEnabled)
	ctx.Then(`^the S3 bucket should block public access$`, s.theS3BucketShouldBlockPublicAccess)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		s.Cleanup(ctx)
		return ctx, nil
	})
}
