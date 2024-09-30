package pkg

import (
	awsrdsinstancev1 "buf.build/gen/go/plantoncloud/project-planton/protocolbuffers/go/project/planton/apis/provider/aws/awsrdsinstance/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsRdsInstance *awsrdsinstancev1.AwsRdsInstance
	Labels         map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *awsrdsinstancev1.AwsRdsInstanceStackInput) *Locals {
	locals := &Locals{}

	//assign value for the locals variable to make it available across the project
	locals.AwsRdsInstance = stackInput.Target

	return locals
}
