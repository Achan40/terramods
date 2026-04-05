output "oidc_role_arn" {
  description = "ARN of the CI/CD role for github"
  value       = aws_iam_role.ci_cd_role.arn
}

output "oidc_provider_arn" {
  description = "ARN of the oidc provider"
  value       = aws_iam_openid_connect_provider.github_oidc.arn
}