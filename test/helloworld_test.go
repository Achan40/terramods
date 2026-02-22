package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestHelloWorld(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../modules/helloworld/",
	}

	// Clean up after test
	defer terraform.Destroy(t, terraformOptions)

	// Run terraform init and apply
	terraform.InitAndApply(t, terraformOptions)

	// Get output
	output := terraform.Output(t, terraformOptions, "hello_world")

	// Assert result
	assert.Equal(t, "Hello, World!", output)
}
