# terramods
A space to store my terragrunt modules

Usage examples can be found in `infra/test`. Note that if remote_state is included and referenced, the resulting s3 bucket and dynamodb table created will need to be manually torn down as it cannot be tracked by terraform.

# Quick Notes
Testing
* `cd test`
* `go test -v`

Terragrunt
* On first run: `terragrunt apply --all --backend-bootstrap` to create the s3 remote backend if remote_state is referenced.

## Adding a New Module

When adding a new module, make sure to create a `README.md` with the
following template before opening a PR:
```markdown
# Module Name

Brief description of what the module does.

<!-- BEGIN_TF_DOCS -->
<!-- END_TF_DOCS -->

```

Then run terraform-docs locally in the repo root to populate it:
```bash
make docs
```

# Updating Modules

After making any changes to a module, regenerate the module documentation from the repo root before opening a PR:
```bash
make docs
```

The CI pipeline will fail if the README.md is out of date.