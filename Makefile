.PHONY: docs
docs:
	@if [ ! -f modules/README.md ]; then \
		echo "Creating missing README.md in modules/"; \
		printf '# Modules\n\n<!-- BEGIN_TF_DOCS -->\n<!-- END_TF_DOCS -->\n' > modules/README.md; \
	fi
	find modules -mindepth 1 -maxdepth 1 -type d | xargs -I {} terraform-docs --config .terraform-docs.yml {}
	terraform-docs --config .terraform-docs.yml modules