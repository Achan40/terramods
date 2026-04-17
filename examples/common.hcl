locals {
  region = "us-east-2"
}

remote_state {
  backend = "s3"
  config = {
    bucket                 = "tfstate-bucket-${get_aws_account_id()}"
    region                 = local.region
    encrypt                = true
    key                    = "${path_relative_to_include()}/terraform.tfstate"
    dynamodb_table         = "tflock-table"
    skip_bucket_versioning = false
  }
}