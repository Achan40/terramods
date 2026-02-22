# terramods
A space to store my terragrunt modules 

Usage examples can be found in `infra/test`. Note that if remote_state is included and referenced, the resulting s3 bucket and dynamodb table created will need to be manually torn down as it cannot be tracked by terraform.

# Quick Notes
Testing
* `cd test`
* `go test -v`

Terragrunt
* On first run: `terragrunt apply --all --backend-bootstrap` to create the s3 remote backend if remote_state is referenced.