.PHONY: docs
docs:
	find modules -mindepth 1 -maxdepth 1 -type d | xargs -I {} terraform-docs --config .terraform-docs.yml {}