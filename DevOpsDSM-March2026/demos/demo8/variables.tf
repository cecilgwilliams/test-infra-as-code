variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "bucket_prefix" {
  description = "Prefix for the S3 bucket name"
  type        = string
  default     = "my-terraform-learning-bucket"
}

variable "environment" {
  description = "Environment tag"
  type        = string
  default     = "test"
}
