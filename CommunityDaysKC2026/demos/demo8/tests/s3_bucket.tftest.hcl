# s3_bucket.tftest.hcl
# Integration test: deploys real AWS resources and asserts their configuration.
# Requires AWS credentials — run with: AWS_PROFILE=<profile> terraform test

# Step 1: helper module generates a unique bucket name
run "setup_tests" {
  module {
    source = "./tests/setup"
  }
}

# Step 2: apply the module and assert the resulting infrastructure
run "create_bucket" {
  command = apply

  variables {
    bucket_name = "${run.setup_tests.bucket_prefix}-tftest-demo"
  }

  assert {
    condition     = aws_s3_bucket.this.bucket == "${run.setup_tests.bucket_prefix}-tftest-demo"
    error_message = "Bucket name does not match expected value: ${aws_s3_bucket.this.bucket}"
  }

  assert {
    condition     = aws_s3_bucket_versioning.this.versioning_configuration[0].status == "Enabled"
    error_message = "Expected versioning to be Enabled, got: ${aws_s3_bucket_versioning.this.versioning_configuration[0].status}"
  }

  assert {
    condition     = one(one(aws_s3_bucket_server_side_encryption_configuration.this.rule).apply_server_side_encryption_by_default).sse_algorithm == "AES256"
    error_message = "Expected bucket to use AES256 server-side encryption"
  }

  assert {
    condition     = aws_s3_bucket_public_access_block.this.block_public_acls == true
    error_message = "Public ACLs should be blocked"
  }
}
