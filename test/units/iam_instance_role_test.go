package test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIAMInstanceRole(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir:    "../../examples/units/iam_instance_role",
		TerraformBinary: "terragrunt",
	}

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	profileName := terraform.Output(t, terraformOptions, "instance_profile_name")
	roleNameOut := terraform.Output(t, terraformOptions, "role_name")
	roleARN := terraform.Output(t, terraformOptions, "role_arn")

	assert.NotEmpty(t, profileName)
	assert.NotEmpty(t, roleNameOut)
	assert.NotEmpty(t, roleARN)

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	require.NoError(t, err)

	iamClient := iam.NewFromConfig(cfg)

	// Role exists with the correct name
	roleOut, err := iamClient.GetRole(ctx, &iam.GetRoleInput{
		RoleName: &roleNameOut,
	})
	require.NoError(t, err)
	assert.Equal(t, roleARN, *roleOut.Role.Arn)

	// Instance profile exists and is linked to the role
	profileOut, err := iamClient.GetInstanceProfile(ctx, &iam.GetInstanceProfileInput{
		InstanceProfileName: &profileName,
	})
	require.NoError(t, err)
	require.Len(t, profileOut.InstanceProfile.Roles, 1)
	assert.Equal(t, roleNameOut, *profileOut.InstanceProfile.Roles[0].RoleName)
}
