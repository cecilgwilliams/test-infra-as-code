# mock.tftest.hcl
# Unit test using mock providers — no AWS credentials or real resources needed.
# Run with: terraform test -filter=tests/mock.tftest.hcl

mock_provider "aws" {}

run "validates_bucket_name_length" {
  command = plan

  variables {
    bucket_name = "my-demo-bucket"
  }

  assert {
    condition     = length(var.bucket_name) >= 3
    error_message = "Bucket name must be at least 3 characters"
  }

  assert {
    condition     = length(var.bucket_name) <= 63
    error_message = "Bucket name must be 63 characters or fewer"
  }
}

run "validates_bucket_name_lowercase" {
  command = plan

  variables {
    bucket_name = "valid-lowercase-bucket-name"
  }

  assert {
    condition     = var.bucket_name == lower(var.bucket_name)
    error_message = "Bucket name must be lowercase, got: ${var.bucket_name}"
  }
}
