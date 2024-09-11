package pkg

import (
	"crypto/rand"
	"github.com/oklog/ulid/v2"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"time"
)

func rdsInstance(ctx *pulumi.Context, locals *Locals, awsProvider *aws.Provider, createdSecurityGroup *ec2.SecurityGroup) (*rds.Instance, error) {
	rdsInstanceArgs := &rds.InstanceArgs{}

	if len(locals.AwsRds.Spec.RdsInstance.SubnetIds) > 0 && locals.AwsRds.Spec.RdsInstance.DbSubnetGroupName == "" {
		createdSubnetGroup, err := subnetGroup(ctx, locals, awsProvider)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create subnet group")
		}
		rdsInstanceArgs.DbSubnetGroupName = createdSubnetGroup.Name
	}
	if locals.AwsRds.Spec.RdsInstance.DbSubnetGroupName != "" {
		rdsInstanceArgs.DbSubnetGroupName = pulumi.String(locals.AwsRds.Spec.RdsInstance.DbSubnetGroupName)
	}

	if locals.AwsRds.Spec.RdsInstance.ParameterGroupName == "" {
		createdParameterGroup, err := parameterGroup(ctx, locals, awsProvider)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create parameter group")
		}
		rdsInstanceArgs.ParameterGroupName = createdParameterGroup.Name
	} else {
		rdsInstanceArgs.ParameterGroupName = pulumi.String(locals.AwsRds.Spec.RdsInstance.ParameterGroupName)
	}

	if locals.AwsRds.Spec.RdsInstance.OptionGroupName == "" {
		createdOptionGroup, err := optionGroup(ctx, locals, awsProvider)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create option group")
		}
		rdsInstanceArgs.OptionGroupName = createdOptionGroup.Name
	} else {
		rdsInstanceArgs.OptionGroupName = pulumi.String(locals.AwsRds.Spec.RdsInstance.OptionGroupName)
	}

	manageMasterUserPassword := locals.AwsRds.Spec.RdsInstance.ManageMasterUserPassword
	if locals.AwsRds.Spec.RdsInstance.ReplicateSourceDb == "" {
		rdsInstanceArgs.Engine = pulumi.String(locals.AwsRds.Spec.RdsInstance.Engine)
		rdsInstanceArgs.EngineVersion = pulumi.String(locals.AwsRds.Spec.RdsInstance.EngineVersion)
		rdsInstanceArgs.AllocatedStorage = pulumi.Int(locals.AwsRds.Spec.RdsInstance.AllocatedStorage)
		if manageMasterUserPassword {
			rdsInstanceArgs.ManageMasterUserPassword = pulumi.Bool(manageMasterUserPassword)
			rdsInstanceArgs.MasterUserSecretKmsKeyId = pulumi.String(locals.AwsRds.Spec.RdsInstance.MasterUserSecretKmsKeyId)
		} else {
			rdsInstanceArgs.Username = pulumi.String(locals.AwsRds.Spec.RdsInstance.Username)
			rdsInstanceArgs.Password = pulumi.String(locals.AwsRds.Spec.RdsInstance.Password)
		}
	} else {
		rdsInstanceArgs.ReplicateSourceDb = pulumi.String(locals.AwsRds.Spec.RdsInstance.ReplicateSourceDb)
	}

	if !locals.AwsRds.Spec.RdsInstance.IsMultiAz {
		rdsInstanceArgs.AvailabilityZone = pulumi.String(locals.AwsRds.Spec.RdsInstance.AvailabilityZone)
	}

	if locals.AwsRds.Spec.RdsInstance.StorageType == "gp3" {
		rdsInstanceArgs.StorageThroughput = pulumi.Int(locals.AwsRds.Spec.RdsInstance.StorageThroughput)
	}

	if !locals.AwsRds.Spec.RdsInstance.SkipFinalSnapshot {
		entropy := ulid.Monotonic(rand.Reader, 0)
		ulidValue := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
		rdsInstanceArgs.FinalSnapshotIdentifier = pulumi.Sprintf("%s-%s", locals.AwsRds.Metadata.Id, ulidValue)
	}

	performanceInsightsEnabled := false
	performanceInsightsKmsKeyId := ""
	performanceInsightsRetentionPeriod := 7
	if locals.AwsRds.Spec.RdsInstance.PerformanceInsights != nil {
		performanceInsightsEnabled = locals.AwsRds.Spec.RdsInstance.PerformanceInsights.IsEnabled
		performanceInsightsKmsKeyId = locals.AwsRds.Spec.RdsInstance.PerformanceInsights.KmsKeyId
		performanceInsightsRetentionPeriod = int(locals.AwsRds.Spec.RdsInstance.PerformanceInsights.RetentionPeriod)
	}

	if performanceInsightsEnabled {
		rdsInstanceArgs.PerformanceInsightsEnabled = pulumi.Bool(performanceInsightsEnabled)
		rdsInstanceArgs.PerformanceInsightsKmsKeyId = pulumi.String(performanceInsightsKmsKeyId)
		rdsInstanceArgs.PerformanceInsightsRetentionPeriod = pulumi.Int(performanceInsightsRetentionPeriod)
	}

	if locals.AwsRds.Spec.RdsInstance.Monitoring != nil {
		rdsInstanceArgs.MonitoringInterval = pulumi.Int(locals.AwsRds.Spec.RdsInstance.Monitoring.MonitoringInterval)
		rdsInstanceArgs.MonitoringRoleArn = pulumi.String(locals.AwsRds.Spec.RdsInstance.Monitoring.MonitoringRoleArn)
	}

	if locals.AwsRds.Spec.RdsInstance.SnapshotIdentifier == "" {
		restoreInTime := &rds.InstanceRestoreToPointInTimeArgs{}
		if locals.AwsRds.Spec.RdsInstance.RestoreToPointInTime != nil {
			restoreInTime = &rds.InstanceRestoreToPointInTimeArgs{
				RestoreTime:                         pulumi.String(locals.AwsRds.Spec.RdsInstance.RestoreToPointInTime.RestoreTime),
				SourceDbInstanceAutomatedBackupsArn: pulumi.String(locals.AwsRds.Spec.RdsInstance.RestoreToPointInTime.SourceDbInstanceAutomatedBackupsArn),
				SourceDbInstanceIdentifier:          pulumi.String(locals.AwsRds.Spec.RdsInstance.RestoreToPointInTime.SourceDbInstanceIdentifier),
				SourceDbiResourceId:                 pulumi.String(locals.AwsRds.Spec.RdsInstance.RestoreToPointInTime.SourceDbiResourceId),
				UseLatestRestorableTime:             pulumi.Bool(locals.AwsRds.Spec.RdsInstance.RestoreToPointInTime.UseLatestRestorableTime),
			}
			rdsInstanceArgs.RestoreToPointInTime = restoreInTime
		}
	}

	vpcSecurityGroupIds := pulumi.ToStringArray(locals.AwsRds.Spec.RdsInstance.AssociateSecurityGroupIds)
	vpcSecurityGroupIds = append(vpcSecurityGroupIds, createdSecurityGroup.ID())

	rdsInstanceArgs.Identifier = pulumi.String(locals.AwsRds.Metadata.Id)
	rdsInstanceArgs.DbName = pulumi.String(locals.AwsRds.Spec.RdsInstance.DbName)
	rdsInstanceArgs.Port = pulumi.Int(locals.AwsRds.Spec.RdsInstance.Port)
	rdsInstanceArgs.CharacterSetName = pulumi.String(locals.AwsRds.Spec.RdsInstance.CharacterSetName)
	rdsInstanceArgs.InstanceClass = pulumi.String(locals.AwsRds.Spec.RdsInstance.InstanceClass)
	rdsInstanceArgs.MaxAllocatedStorage = pulumi.Int(locals.AwsRds.Spec.RdsInstance.MaxAllocatedStorage)
	rdsInstanceArgs.StorageEncrypted = pulumi.Bool(locals.AwsRds.Spec.RdsInstance.StorageEncrypted)
	rdsInstanceArgs.KmsKeyId = pulumi.String(locals.AwsRds.Spec.RdsInstance.KmsKeyId)
	rdsInstanceArgs.VpcSecurityGroupIds = vpcSecurityGroupIds
	rdsInstanceArgs.MultiAz = pulumi.Bool(locals.AwsRds.Spec.RdsInstance.IsMultiAz)
	rdsInstanceArgs.CaCertIdentifier = pulumi.String(locals.AwsRds.Spec.RdsInstance.CaCertIdentifier)
	rdsInstanceArgs.LicenseModel = pulumi.String(locals.AwsRds.Spec.RdsInstance.LicenseModel)
	rdsInstanceArgs.StorageType = pulumi.String(locals.AwsRds.Spec.RdsInstance.StorageType)
	rdsInstanceArgs.Iops = pulumi.Int(locals.AwsRds.Spec.RdsInstance.Iops)
	rdsInstanceArgs.PubliclyAccessible = pulumi.Bool(locals.AwsRds.Spec.RdsInstance.IsPubliclyAccessible)
	rdsInstanceArgs.SnapshotIdentifier = pulumi.String(locals.AwsRds.Spec.RdsInstance.SnapshotIdentifier)
	rdsInstanceArgs.AllowMajorVersionUpgrade = pulumi.Bool(locals.AwsRds.Spec.RdsInstance.AllowMajorVersionUpgrade)
	rdsInstanceArgs.AutoMinorVersionUpgrade = pulumi.Bool(locals.AwsRds.Spec.RdsInstance.AutoMinorVersionUpgrade)
	rdsInstanceArgs.ApplyImmediately = pulumi.Bool(locals.AwsRds.Spec.RdsInstance.ApplyImmediately)
	rdsInstanceArgs.MaintenanceWindow = pulumi.String(locals.AwsRds.Spec.RdsInstance.MaintenanceWindow)
	rdsInstanceArgs.CopyTagsToSnapshot = pulumi.Bool(locals.AwsRds.Spec.RdsInstance.CopyTagsToSnapshot)
	rdsInstanceArgs.BackupRetentionPeriod = pulumi.Int(locals.AwsRds.Spec.RdsInstance.BackupRetentionPeriod)
	rdsInstanceArgs.BackupWindow = pulumi.String(locals.AwsRds.Spec.RdsInstance.BackupWindow)
	rdsInstanceArgs.DeletionProtection = pulumi.Bool(locals.AwsRds.Spec.RdsInstance.DeletionProtection)
	rdsInstanceArgs.SkipFinalSnapshot = pulumi.Bool(locals.AwsRds.Spec.RdsInstance.SkipFinalSnapshot)
	rdsInstanceArgs.Timezone = pulumi.String(locals.AwsRds.Spec.RdsInstance.Timezone)
	rdsInstanceArgs.IamDatabaseAuthenticationEnabled = pulumi.Bool(locals.AwsRds.Spec.RdsInstance.IamDatabaseAuthenticationEnabled)
	rdsInstanceArgs.EnabledCloudwatchLogsExports = pulumi.ToStringArray(locals.AwsRds.Spec.RdsInstance.EnabledCloudwatchLogsExports)
	rdsInstanceArgs.Tags = pulumi.ToStringMap(locals.Labels)

	// Create RDS Instance
	rdsInstance, err := rds.NewInstance(ctx, "rdsInstance", rdsInstanceArgs,
		pulumi.Provider(awsProvider), pulumi.DependsOn([]pulumi.Resource{createdSecurityGroup}))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create rds instance")
	}

	return rdsInstance, nil
}
