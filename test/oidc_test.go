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
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir:    "../infra/test/oidc",
		TerraformBinary: "terragrunt",
		Vars: map[string]interface{}{
			"github_repo": "Achan40/terramods",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	oidcArn := terraform.Output(t, terraformOptions, "oidc_provider_arn")

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	assert.NoError(t, err)

	client := iam.NewFromConfig(cfg)

	resp, err := client.GetOpenIDConnectProvider(context.TODO(), &iam.GetOpenIDConnectProviderInput{
		OpenIDConnectProviderArn: &oidcArn,
	})

	assert.NoError(t, err)
	assert.Equal(t, "token.actions.githubusercontent.com", *resp.Url)
}
