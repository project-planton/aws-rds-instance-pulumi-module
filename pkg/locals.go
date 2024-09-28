package pkg

import (
	"github.com/plantoncloud/project-planton/apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/awsrdsinstance"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsRdsInstance *awsrdsinstance.AwsRdsInstance
	Labels         map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *awsrdsinstance.AwsRdsInstanceStackInput) *Locals {
	locals := &Locals{}

	//assign value for the locals variable to make it available across the project
	locals.AwsRdsInstance = stackInput.Target

	return locals
}
