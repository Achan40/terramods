package test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestGithubOIDCProvider(t *testing.T) {
	// Test oidc provider is created
	// DO NOT run this test in CI/CD, there can be only one oidc_provider_arn per provider url.
	// run go with -short flag to skip

	if testing.Short() {
		t.Skip("Skipping OIDC test in short mode.")
	}

	// Running test in pipeline will delete the creditials that is needed to run the pipeline...
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir:    "../../examples/required/github_oidc",
		TerraformBinary: "terragrunt",
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	oidcArn := terraform.Output(t, terraformOptions, "oidc_provider_arn")

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	assert.NoError(t, err)

	client := iam.NewFromConfig(cfg)

	// Get OIDC provider details
	resp, err := client.GetOpenIDConnectProvider(context.TODO(), &iam.GetOpenIDConnectProviderInput{
		OpenIDConnectProviderArn: &oidcArn,
	})

	assert.NoError(t, err)
	assert.Equal(t, "token.actions.githubusercontent.com", *resp.Url)
}
