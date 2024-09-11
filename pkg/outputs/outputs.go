package outputs

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/awsrds"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/automationapi/autoapistackoutput"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

const (
	RdsInstanceEndpoint = "rds-instance-endpoint"
	RdsInstanceId       = "rds-instance-id"
	RdsInstanceArn      = "rds-instance-arn"
	RdsInstanceAddress  = "rds-instance-address"
	RdsSubnetGroup      = "rds-subnet-group"
	RdsSecurityGroup    = "rds-security-group"
	RdsParameterGroup   = "rds-parameter-group"
	RdsOptionsGroup     = "rds-options-group"
)

func PulumiOutputsToStackOutputsConverter(pulumiOutputs auto.OutputMap,
	input *awsrds.AwsRdsStackInput) *awsrds.AwsRdsStackOutputs {
	return &awsrds.AwsRdsStackOutputs{
		RdsInstanceEndpoint: autoapistackoutput.GetVal(pulumiOutputs, RdsInstanceEndpoint),
		RdsInstanceId:       autoapistackoutput.GetVal(pulumiOutputs, RdsInstanceId),
		RdsInstanceArn:      autoapistackoutput.GetVal(pulumiOutputs, RdsInstanceArn),
		RdsInstanceAddress:  autoapistackoutput.GetVal(pulumiOutputs, RdsInstanceAddress),
		RdsSubnetGroup:      autoapistackoutput.GetVal(pulumiOutputs, RdsSubnetGroup),
		RdsSecurityGroup:    autoapistackoutput.GetVal(pulumiOutputs, RdsSecurityGroup),
		RdsParameterGroup:   autoapistackoutput.GetVal(pulumiOutputs, RdsParameterGroup),
		RdsOptionsGroup:     autoapistackoutput.GetVal(pulumiOutputs, RdsOptionsGroup),
	}
}
