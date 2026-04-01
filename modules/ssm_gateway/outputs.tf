output "instance_id" {
  description = "ID of the SSM gateway EC2 instance"
  value       = aws_instance.ssm_gateway.id
}

output "instance_private_ip" {
  description = "Private IP of the SSM gateway EC2 instance"
  value       = aws_instance.ssm_gateway.private_ip
}

output "security_group_id" {
  description = "ID of the SSM gateway security group"
  value       = aws_security_group.ssm_gateway.id
}
