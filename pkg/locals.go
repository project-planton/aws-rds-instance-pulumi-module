package pkg

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/awsrds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsRds *awsrds.AwsRds
	Labels map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *awsrds.AwsRdsStackInput) *Locals {
	locals := &Locals{}

	//assign value for the locals variable to make it available across the project
	locals.AwsRds = stackInput.Target

	return locals
}
