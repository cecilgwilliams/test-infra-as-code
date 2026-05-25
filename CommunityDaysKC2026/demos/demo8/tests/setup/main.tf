# Helper module: generates a unique bucket name prefix for each test run
terraform {
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "~> 3.5"
    }
  }
}

resource "random_pet" "bucket_prefix" {
  length = 4
}

output "bucket_prefix" {
  value = random_pet.bucket_prefix.id
}
