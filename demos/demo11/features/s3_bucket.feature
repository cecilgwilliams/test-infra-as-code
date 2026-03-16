# features/s3_bucket.feature
# BDD feature spec for S3 bucket security requirements.
# Run: go test -v -run TestBDDS3

Feature: Secure S3 bucket
  As a platform engineer
  I need S3 buckets to follow security best practices
  So that our data is protected at rest

  Background:
    Given the Terraform module has been applied

  Scenario: Bucket exists
    Then the S3 bucket should exist in "us-east-1"

  Scenario: Bucket has encryption
    Then the S3 bucket should have "AES256" encryption

  Scenario: Bucket has versioning
    Then the S3 bucket should have versioning enabled

  Scenario: Bucket blocks public access
    Then the S3 bucket should block public access
