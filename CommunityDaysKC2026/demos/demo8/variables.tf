variable "bucket_name" {
  description = "Name of the S3 bucket — provided by the test helper module"
  type        = string
}

variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}
