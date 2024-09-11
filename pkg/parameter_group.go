package pkg

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func parameterGroup(ctx *pulumi.Context, locals *Locals, awsProvider *aws.Provider) (*rds.ParameterGroup, error) {
	var parameterGroupParameterArray = rds.ParameterGroupParameterArray{}
	for _, parameter := range locals.AwsRds.Spec.RdsInstance.Parameters {
		parameterGroupParameterArray = append(parameterGroupParameterArray, &rds.ParameterGroupParameterArgs{
			ApplyMethod: pulumi.String(parameter.ApplyMethod),
			Name:        pulumi.String(parameter.Name),
			Value:       pulumi.String(parameter.Value),
		})

	}

	parameterGroupArgs := &rds.ParameterGroupArgs{
		NamePrefix: pulumi.Sprintf("%s-", locals.AwsRds.Metadata.Id),
		Family:     pulumi.String(locals.AwsRds.Spec.RdsInstance.DbParameterGroup),
		Tags:       pulumi.ToStringMap(locals.Labels),
		Parameters: rds.ParameterGroupParameterArray{
			&rds.ParameterGroupParameterArgs{
				ApplyMethod: pulumi.String("pending-reboot"),
				Name:        pulumi.String("max_connections"),
				Value:       pulumi.String("100"),
			},
		},
	}
	// Create RDS Parameter Group
	rdsParameterGroup, err := rds.NewParameterGroup(ctx, "rds-parameter-group", parameterGroupArgs, pulumi.Provider(awsProvider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create rds parameter group")
	}

	return rdsParameterGroup, nil
}
