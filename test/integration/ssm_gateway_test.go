package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestSSMGateway(t *testing.T) {
	t.Parallel()

	region := "us-east-2"
	uniqueID := fmt.Sprintf("%d", time.Now().Unix())

	// deploy a vpc first so the gateway has a subnet to land in
	vpcOptions := &terraform.Options{
		TerraformDir:    "../../infra/test/private_vpc",
		TerraformBinary: "terragrunt",
		Vars: map[string]interface{}{
			"vpc_name": "test-ssm-vpc-" + uniqueID,
			"region":   region,
			"vpc_cidr": "10.0.0.0/16",
			"availability_zones": []string{
				"us-east-2a",
			},
		},
	}

	defer terraform.Destroy(t, vpcOptions)
	terraform.InitAndApply(t, vpcOptions)

	vpcID := terraform.Output(t, vpcOptions, "vpc_id")
	privateSubnetIDs := terraform.OutputList(t, vpcOptions, "private_subnet_ids")

	gatewayOptions := &terraform.Options{
		TerraformDir:    "../../infra/test/ssm_gateway",
		TerraformBinary: "terragrunt",
		Vars: map[string]interface{}{
			"name":      "test-ssm-gateway-" + uniqueID,
			"vpc_id":    vpcID,
			"subnet_id": privateSubnetIDs[0],
		},
	}

	defer terraform.Destroy(t, gatewayOptions)
	terraform.InitAndApply(t, gatewayOptions)

	instanceID := terraform.Output(t, gatewayOptions, "instance_id")
	securityGroupID := terraform.Output(t, gatewayOptions, "security_group_id")

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	assert.NoError(t, err)

	ec2Client := ec2.NewFromConfig(cfg)
	iamClient := iam.NewFromConfig(cfg)

	// validate instance exists in the correct vpc and has no public ip
	instanceOut, err := ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	})
	assert.NoError(t, err)
	assert.Len(t, instanceOut.Reservations, 1)

	instance := instanceOut.Reservations[0].Instances[0]
	assert.Equal(t, vpcID, *instance.VpcId)
	assert.Equal(t, privateSubnetIDs[0], *instance.SubnetId)
	assert.Nil(t, instance.PublicIpAddress)

	// validate the instance profile is attached
	assert.NotNil(t, instance.IamInstanceProfile)

	// validate the instance profile has the ssm managed policy attached
	profileName := fmt.Sprintf("test-ssm-gateway-%s-ssm-gateway-profile", uniqueID)
	profileOut, err := iamClient.GetInstanceProfile(ctx, &iam.GetInstanceProfileInput{
		InstanceProfileName: &profileName,
	})
	assert.NoError(t, err)
	assert.Len(t, profileOut.InstanceProfile.Roles, 1)

	roleName := *profileOut.InstanceProfile.Roles[0].RoleName
	attachedPolicies, err := iamClient.ListAttachedRolePolicies(ctx, &iam.ListAttachedRolePoliciesInput{
		RoleName: &roleName,
	})
	assert.NoError(t, err)

	hasSsmPolicy := false
	for _, policy := range attachedPolicies.AttachedPolicies {
		if *policy.PolicyName == "AmazonSSMManagedInstanceCore" {
			hasSsmPolicy = true
			break
		}
	}
	assert.True(t, hasSsmPolicy, "expected AmazonSSMManagedInstanceCore to be attached")

	// validate security group has no inbound rules
	sgOut, err := ec2Client.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{
		GroupIds: []string{securityGroupID},
	})
	assert.NoError(t, err)
	assert.Len(t, sgOut.SecurityGroups, 1)

	sg := sgOut.SecurityGroups[0]
	assert.Empty(t, sg.IpPermissions, "expected no inbound rules on SSM gateway security group")
}
