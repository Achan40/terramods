output "security_group_id" {
  description = "ID of the shared instance security group"
  value       = aws_security_group.instance.id
}
