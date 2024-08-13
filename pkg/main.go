package pkg

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/awsrds"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	Input  *awsrds.AwsRdsStackInput
	Labels map[string]string
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	//create aws provider using the credentials from the input
	_, err := pulumiawsprovider.GetNative(ctx, s.Input.AwsCredential)
	if err != nil {
		return errors.Wrap(err, "failed to create aws provider")
	}

	return nil
}
