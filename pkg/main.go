package pkg

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/aws-rds-pulumi-module/pkg/outputs"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/awsrds"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	StackInput *awsrds.AwsRdsStackInput
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	locals := initializeLocals(ctx, s.StackInput)

	awsCredential := s.StackInput.AwsCredential

	//create aws provider using the credentials from the input
	awsProvider, err := aws.NewProvider(ctx,
		"classic-provider",
		&aws.ProviderArgs{
			AccessKey: pulumi.String(awsCredential.Spec.AccessKeyId),
			SecretKey: pulumi.String(awsCredential.Spec.SecretAccessKey),
			Region:    pulumi.String(awsCredential.Spec.Region),
		})
	if err != nil {
		return errors.Wrap(err, "failed to create aws provider")
	}

	createdSecurityGroup, err := securityGroup(ctx, locals, awsProvider)
	if err != nil {
		return errors.Wrap(err, "failed to create default security group")
	}

	// Create RDS Instance
	createRdsInstance, err := rdsInstance(ctx, locals, awsProvider, createdSecurityGroup)
	if err != nil {
		return errors.Wrap(err, "failed to create rds instance")
	}

	// Export Outputs
	ctx.Export(outputs.RdsInstanceEndpoint, createRdsInstance.Endpoint)
	ctx.Export(outputs.RdsInstanceId, createRdsInstance.ResourceId)
	ctx.Export(outputs.RdsInstanceArn, createRdsInstance.Arn)
	ctx.Export(outputs.RdsInstanceAddress, createRdsInstance.Address)
	ctx.Export(outputs.RdsSecurityGroup, createdSecurityGroup.Name)
	ctx.Export(outputs.RdsParameterGroup, createRdsInstance.ParameterGroupName)
	ctx.Export(outputs.RdsSubnetGroup, createRdsInstance.DbSubnetGroupName)
	ctx.Export(outputs.RdsOptionsGroup, createRdsInstance.OptionGroupName)
	return nil
}
