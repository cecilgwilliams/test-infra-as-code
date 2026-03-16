variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "subnet_count" {
  description = "Number of subnets to create"
  type        = number
  default     = 3
}

variable "name_prefix" {
  description = "Prefix for resource names"
  type        = string
  default     = "test"
}

variable "environment" {
  description = "Environment tag"
  type        = string
  default     = "test"
}
