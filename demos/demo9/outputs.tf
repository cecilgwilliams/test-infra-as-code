output "vpc_id" {
  description = "ID of the created VPC"
  value       = aws_vpc.this.id
}

output "vpc_cidr" {
  description = "CIDR block of the VPC"
  value       = aws_vpc.this.cidr_block
}

output "subnet_ids" {
  description = "IDs of the created subnets"
  value       = aws_subnet.this[*].id
}

output "subnet_count" {
  description = "Number of subnets created"
  value       = length(aws_subnet.this)
}
