# main.tf
# Intentionally bad Terraform config for demo purposes.
#
# Run these against it:
#   terraform validate       # ✓ Syntax check
#   terraform fmt -check     # ✓ Format check
#   tflint                   # ⚠ Lint warnings (needs .tflint.hcl + tflint --init)
#   tfsec .                  # ❌ Security issues found!
#
# Expected tflint output:
#   Warning: `bucket_name` variable has no type (terraform_typed_variables)
#
# Expected tfsec output:
#   CRITICAL AWS017: S3 Bucket does not have encryption enabled
#   HIGH     AWS002: S3 Bucket does not have logging enabled
#   MEDIUM   AWS077: S3 Data should be versioned

terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

variable "bucket_name" {
  # no type defined - tflint will warn: terraform_typed_variables
  default = "my-bucket"
}

resource "aws_s3_bucket" "bad" {
  bucket = var.bucket_name
  # Missing encryption, versioning, etc.
}
