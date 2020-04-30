package tests

import (
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAWSModules(t *testing.T) {
	t.Parallel()

	name := random.UniqueId()
	environment := "test"
	used_by := []interface{}{"random_resource"}
	awsRegion := aws.GetRandomStableRegion(t, nil, nil)
	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/security_group",
		Vars: map[string]interface{}{
			"name":    name,
			"used_by": used_by,
			"tags":    map[string]interface{}{"env": environment},
			"ingress": []map[string]interface{}{
				map[string]interface{}{
					"to_port":   "22",
					"from_port": "22",
				},
			},
		},
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	output := terraform.OutputAll(t, terraformOptions)

	assert.Equal(t, name, output["sg_name"])
	assert.Equal(t, used_by, output["sg_used_by"])
	assert.Equal(t, environment, output["sg_tags"].(map[string]interface{})["env"].(string))
	assert.NotEqual(t, []interface{}{}, output["sg_ingress"])
	assert.Equal(t, []interface{}{}, output["sg_egress"])
	assert.True(t, strings.Contains(output["sg_arn"].(string), "arn"))
}
