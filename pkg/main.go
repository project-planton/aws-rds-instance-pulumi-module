package pkg

import (
	awsrdsinstancev1 "buf.build/gen/go/project-planton/apis/protocolbuffers/go/project/planton/provider/aws/awsrdsinstance/v1"
	"github.com/pkg/errors"
	"github.com/project-planton/aws-rds-instance-pulumi-module/pkg/outputs"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *awsrdsinstancev1.AwsRdsInstanceStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	awsCredential := stackInput.AwsCredential

	//create aws provider using the credentials from the input
	awsProvider, err := aws.NewProvider(ctx,
		"classic-provider",
		&aws.ProviderArgs{
			AccessKey: pulumi.String(awsCredential.AccessKeyId),
			SecretKey: pulumi.String(awsCredential.SecretAccessKey),
			Region:    pulumi.String(awsCredential.Region),
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
