# s3_encrypted.rego
# Open Policy Agent (OPA) policy that enforces S3 encryption.
#
# Usage:
#   terraform plan -out=tfplan
#   terraform show -json tfplan > tfplan.json
#   opa eval --data s3_encrypted.rego --input tfplan.json "data.terraform.deny"

package terraform

deny contains msg if {
  resource := input.resource_changes[_]
  resource.type == "aws_s3_bucket"
  not resource.change.after.server_side_encryption_configuration
  msg := sprintf("S3 bucket '%s' must have encryption enabled", [resource.name])
}

deny contains msg if {
  resource := input.resource_changes[_]
  resource.type == "aws_s3_bucket"
  not resource.change.after.versioning
  msg := sprintf("S3 bucket '%s' must have versioning enabled", [resource.name])
}
