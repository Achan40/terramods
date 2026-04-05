output "vpc_id" {
  description = "Id of the VPC"
  value = aws_vpc.main.id
}

output "private_subnet_ids" {
  description = "Ids of the private subnets within the VPC"
  value = values(aws_subnet.private)[*].id
}

output "public_subnet_ids" {
  description = "Ids of the public subnets within the VPC"
  value = values(aws_subnet.public)[*].id
}

output "eice_security_group_id" {
  description = "ID of the EC2 Instance Connect Endpoint security group"
  value       = aws_security_group.eice.id
}

output "region" {
  description = "AWS region the VPC was deployed into"
  value       = var.region
}

output "vpc_cidr" {
  description = "CIDR block of the VPC"
  value       = aws_vpc.main.cidr_block
}