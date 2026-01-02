// Package schema provides CloudFormation resource schema validation.
package schema

// PropertyConstraints defines validation constraints for a property.
type PropertyConstraints struct {
	// Pattern is a regex pattern the string value must match.
	Pattern string
	// MinLength is the minimum string length.
	MinLength *int
	// MaxLength is the maximum string length.
	MaxLength *int
	// MinValue is the minimum numeric value.
	MinValue *float64
	// MaxValue is the maximum numeric value.
	MaxValue *float64
	// MinItems is the minimum array length.
	MinItems *int
	// MaxItems is the maximum array length.
	MaxItems *int
	// UniqueItems indicates that list items must be unique.
	UniqueItems bool
}

// ResourceConstraints defines resource-level constraints.
type ResourceConstraints struct {
	// ReadOnlyProperties are properties that cannot be set by users (create-only or read-only).
	ReadOnlyProperties []string
	// MutuallyExclusive defines sets of properties where only one can be specified.
	MutuallyExclusive [][]string
	// DependentRequired maps a property to other properties that must be present if it is.
	DependentRequired map[string][]string
	// DependentExcluded maps a property to other properties that must NOT be present if it is.
	DependentExcluded map[string][]string
	// OneOf defines sets where exactly one property must be present.
	OneOf [][]string
	// AnyOf defines sets where at least one property must be present.
	AnyOf [][]string
}

// intPtr returns a pointer to an int.
func intPtr(v int) *int { return &v }

// floatPtr returns a pointer to a float64.
func floatPtr(v float64) *float64 { return &v }

// resourceConstraints maps resource type -> property name -> constraints.
// Based on CloudFormation Registry Schemas.
var resourceConstraints = map[string]map[string]PropertyConstraints{
	"AWS::Lambda::Function": {
		"FunctionName": {
			Pattern:   `^[a-zA-Z0-9-_]+$`,
			MinLength: intPtr(1),
			MaxLength: intPtr(64),
		},
		"Description": {
			MaxLength: intPtr(256),
		},
		"Handler": {
			MaxLength: intPtr(128),
			Pattern:   `^[^\s]+$`,
		},
		"MemorySize": {
			MinValue: floatPtr(128),
			MaxValue: floatPtr(10240),
		},
		"Timeout": {
			MinValue: floatPtr(1),
			MaxValue: floatPtr(900),
		},
		"ReservedConcurrentExecutions": {
			MinValue: floatPtr(0),
		},
	},
	"AWS::S3::Bucket": {
		"BucketName": {
			Pattern:   `^[a-z0-9][a-z0-9.-]*[a-z0-9]$`,
			MinLength: intPtr(3),
			MaxLength: intPtr(63),
		},
	},
	"AWS::EC2::SecurityGroup": {
		"GroupName": {
			MaxLength: intPtr(255),
		},
		"GroupDescription": {
			MaxLength: intPtr(255),
		},
	},
	"AWS::EC2::Instance": {
		"InstanceType": {
			Pattern: `^[a-z][a-z0-9-]+\.[a-z0-9]+$`,
		},
	},
	"AWS::IAM::Role": {
		"RoleName": {
			Pattern:   `^[\w+=,.@-]+$`,
			MinLength: intPtr(1),
			MaxLength: intPtr(64),
		},
		"Description": {
			MaxLength: intPtr(1000),
		},
		"MaxSessionDuration": {
			MinValue: floatPtr(3600),
			MaxValue: floatPtr(43200),
		},
		"Path": {
			Pattern:   `^\/.*\/$`,
			MinLength: intPtr(1),
			MaxLength: intPtr(512),
		},
	},
	"AWS::IAM::Policy": {
		"PolicyName": {
			Pattern:   `^[\w+=,.@-]+$`,
			MinLength: intPtr(1),
			MaxLength: intPtr(128),
		},
	},
	"AWS::IAM::User": {
		"UserName": {
			Pattern:   `^[\w+=,.@-]+$`,
			MinLength: intPtr(1),
			MaxLength: intPtr(64),
		},
		"Path": {
			Pattern:   `^\/.*\/$`,
			MinLength: intPtr(1),
			MaxLength: intPtr(512),
		},
	},
	"AWS::DynamoDB::Table": {
		"TableName": {
			Pattern:   `^[a-zA-Z0-9_.-]+$`,
			MinLength: intPtr(3),
			MaxLength: intPtr(255),
		},
	},
	"AWS::SQS::Queue": {
		"QueueName": {
			MaxLength: intPtr(80),
		},
		"DelaySeconds": {
			MinValue: floatPtr(0),
			MaxValue: floatPtr(900),
		},
		"MaximumMessageSize": {
			MinValue: floatPtr(1024),
			MaxValue: floatPtr(262144),
		},
		"MessageRetentionPeriod": {
			MinValue: floatPtr(60),
			MaxValue: floatPtr(1209600),
		},
		"VisibilityTimeout": {
			MinValue: floatPtr(0),
			MaxValue: floatPtr(43200),
		},
	},
	"AWS::SNS::Topic": {
		"TopicName": {
			MaxLength: intPtr(256),
		},
		"DisplayName": {
			MaxLength: intPtr(100),
		},
	},
	"AWS::ECS::Service": {
		"ServiceName": {
			MaxLength: intPtr(255),
		},
		"DesiredCount": {
			MinValue: floatPtr(0),
		},
		"HealthCheckGracePeriodSeconds": {
			MinValue: floatPtr(0),
			MaxValue: floatPtr(2147483647),
		},
	},
	"AWS::ECS::TaskDefinition": {
		"Family": {
			MaxLength: intPtr(255),
		},
		"Cpu": {
			Pattern: `^(256|512|1024|2048|4096|8192|16384)$`,
		},
		"Memory": {
			Pattern: `^[0-9]+$`,
		},
	},
	"AWS::Logs::LogGroup": {
		"LogGroupName": {
			Pattern:   `^[\.\-_/#A-Za-z0-9]+$`,
			MinLength: intPtr(1),
			MaxLength: intPtr(512),
		},
		"RetentionInDays": {
			// Valid values: 1, 3, 5, 7, 14, 30, 60, 90, 120, 150, 180, 365, 400, 545, 731, 1096, 1827, 2192, 2557, 2922, 3288, 3653
			MinValue: floatPtr(1),
		},
	},
	"AWS::CloudWatch::Alarm": {
		"AlarmName": {
			MinLength: intPtr(1),
			MaxLength: intPtr(255),
		},
		"AlarmDescription": {
			MaxLength: intPtr(1024),
		},
		"EvaluationPeriods": {
			MinValue: floatPtr(1),
		},
		"Period": {
			MinValue: floatPtr(1),
		},
	},
	"AWS::ApiGateway::RestApi": {
		"Name": {
			MinLength: intPtr(1),
		},
		"Description": {
			MaxLength: intPtr(1024),
		},
	},
	"AWS::EC2::VPC": {
		"CidrBlock": {
			Pattern: `^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(\/([0-9]|[1-2][0-9]|3[0-2]))$`,
		},
	},
	"AWS::EC2::Subnet": {
		"CidrBlock": {
			Pattern: `^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(\/([0-9]|[1-2][0-9]|3[0-2]))$`,
		},
	},
	"AWS::KMS::Key": {
		"Description": {
			MaxLength: intPtr(8192),
		},
	},
	"AWS::SecretsManager::Secret": {
		"Name": {
			MinLength: intPtr(1),
			MaxLength: intPtr(256),
		},
		"Description": {
			MaxLength: intPtr(2048),
		},
	},
	"AWS::SSM::Parameter": {
		"Name": {
			MinLength: intPtr(1),
			MaxLength: intPtr(2048),
		},
		"Description": {
			MaxLength: intPtr(1024),
		},
	},
	"AWS::StepFunctions::StateMachine": {
		"StateMachineName": {
			MinLength: intPtr(1),
			MaxLength: intPtr(80),
		},
	},
	"AWS::Events::Rule": {
		"Name": {
			MinLength: intPtr(1),
			MaxLength: intPtr(64),
			Pattern:   `^[\.\-_A-Za-z0-9]+$`,
		},
		"Description": {
			MaxLength: intPtr(512),
		},
	},
}

// GetPropertyConstraints returns constraints for a resource property.
// Returns nil if no constraints are defined.
func GetPropertyConstraints(resourceType, propertyName string) *PropertyConstraints {
	if props, ok := resourceConstraints[resourceType]; ok {
		if c, ok := props[propertyName]; ok {
			return &c
		}
	}
	return nil
}

// HasConstraints returns true if any constraints are defined for the resource type.
func HasConstraints(resourceType string) bool {
	_, ok := resourceConstraints[resourceType]
	return ok
}

// resourceLevelConstraints maps resource type -> resource-level constraints.
var resourceLevelConstraints = map[string]ResourceConstraints{
	"AWS::Lambda::Function": {
		ReadOnlyProperties: []string{"Arn", "SnapStartResponse"},
		MutuallyExclusive: [][]string{
			{"Code.S3Bucket", "Code.ImageUri"},
			{"Code.ZipFile", "Code.S3Bucket"},
		},
		DependentRequired: map[string][]string{
			"Code.S3Bucket": {"Code.S3Key"},
		},
	},
	"AWS::S3::Bucket": {
		ReadOnlyProperties: []string{"Arn", "DomainName", "DualStackDomainName", "RegionalDomainName", "WebsiteURL"},
	},
	"AWS::EC2::Instance": {
		ReadOnlyProperties: []string{"AvailabilityZone", "PrivateDnsName", "PrivateIp", "PublicDnsName", "PublicIp"},
		MutuallyExclusive: [][]string{
			{"SecurityGroups", "SecurityGroupIds"},
			{"SubnetId", "NetworkInterfaces"},
		},
	},
	"AWS::EC2::SecurityGroup": {
		ReadOnlyProperties: []string{"GroupId", "VpcId"},
		DependentRequired: map[string][]string{
			"SecurityGroupIngress": {"GroupDescription"},
			"SecurityGroupEgress":  {"GroupDescription"},
		},
	},
	"AWS::IAM::Role": {
		ReadOnlyProperties: []string{"Arn", "RoleId"},
	},
	"AWS::IAM::User": {
		ReadOnlyProperties: []string{"Arn"},
	},
	"AWS::DynamoDB::Table": {
		ReadOnlyProperties: []string{"Arn", "StreamArn"},
		DependentRequired: map[string][]string{
			"GlobalSecondaryIndexes": {"AttributeDefinitions"},
			"LocalSecondaryIndexes":  {"AttributeDefinitions"},
		},
	},
	"AWS::SQS::Queue": {
		ReadOnlyProperties: []string{"Arn", "QueueUrl"},
		MutuallyExclusive: [][]string{
			{"FifoQueue", "ContentBasedDeduplication"},
		},
		DependentRequired: map[string][]string{
			"ContentBasedDeduplication": {"FifoQueue"},
		},
	},
	"AWS::SNS::Topic": {
		ReadOnlyProperties: []string{"TopicArn"},
	},
	"AWS::RDS::DBInstance": {
		ReadOnlyProperties: []string{"Endpoint.Address", "Endpoint.Port", "Endpoint.HostedZoneId"},
		MutuallyExclusive: [][]string{
			{"DBSnapshotIdentifier", "SourceDBInstanceIdentifier"},
			{"DBSnapshotIdentifier", "MasterUsername"},
		},
		DependentRequired: map[string][]string{
			"MasterUsername": {"MasterUserPassword"},
		},
		DependentExcluded: map[string][]string{
			"DBSnapshotIdentifier": {"MasterUsername", "MasterUserPassword"},
		},
	},
	"AWS::ECS::Service": {
		ReadOnlyProperties: []string{"Name", "ServiceArn"},
		MutuallyExclusive: [][]string{
			{"LaunchType", "CapacityProviderStrategy"},
		},
		DependentRequired: map[string][]string{
			"LoadBalancers": {"Role"},
		},
	},
	"AWS::ECS::TaskDefinition": {
		ReadOnlyProperties: []string{"TaskDefinitionArn"},
		DependentRequired: map[string][]string{
			"Cpu":    {"NetworkMode"},
			"Memory": {"NetworkMode"},
		},
	},
	"AWS::Logs::LogGroup": {
		ReadOnlyProperties: []string{"Arn"},
	},
	"AWS::CloudWatch::Alarm": {
		MutuallyExclusive: [][]string{
			{"Metrics", "MetricName"},
			{"Metrics", "Statistic"},
		},
		OneOf: [][]string{
			{"MetricName", "Metrics"},
		},
	},
	"AWS::KMS::Key": {
		ReadOnlyProperties: []string{"Arn", "KeyId"},
	},
	"AWS::SecretsManager::Secret": {
		ReadOnlyProperties: []string{"Id"},
		MutuallyExclusive: [][]string{
			{"SecretString", "GenerateSecretString"},
		},
	},
	"AWS::SSM::Parameter": {
		MutuallyExclusive: [][]string{
			{"Value", "DataType"},
		},
	},
	"AWS::StepFunctions::StateMachine": {
		ReadOnlyProperties: []string{"Arn", "StateMachineRevisionId"},
		MutuallyExclusive: [][]string{
			{"Definition", "DefinitionS3Location"},
			{"Definition", "DefinitionString"},
		},
		OneOf: [][]string{
			{"Definition", "DefinitionS3Location", "DefinitionString"},
		},
	},
	"AWS::ApiGateway::RestApi": {
		ReadOnlyProperties: []string{"RootResourceId", "RestApiId"},
		MutuallyExclusive: [][]string{
			{"Body", "BodyS3Location"},
		},
	},
	"AWS::CloudFormation::Stack": {
		MutuallyExclusive: [][]string{
			{"TemplateBody", "TemplateURL"},
		},
		OneOf: [][]string{
			{"TemplateBody", "TemplateURL"},
		},
	},
	"AWS::ElasticLoadBalancingV2::TargetGroup": {
		DependentRequired: map[string][]string{
			"HealthCheckPort":     {"Protocol"},
			"HealthCheckProtocol": {"Protocol"},
		},
	},
	"AWS::Events::Rule": {
		MutuallyExclusive: [][]string{
			{"EventPattern", "ScheduleExpression"},
		},
		AnyOf: [][]string{
			{"EventPattern", "ScheduleExpression"},
		},
	},
}

// uniqueItemsProperties lists properties that require unique items.
var uniqueItemsProperties = map[string]map[string]bool{
	"AWS::EC2::SecurityGroup": {
		"SecurityGroupIngress": true,
		"SecurityGroupEgress":  true,
	},
	"AWS::IAM::Role": {
		"ManagedPolicyArns": true,
		"Policies":          true,
	},
	"AWS::IAM::User": {
		"Groups":            true,
		"ManagedPolicyArns": true,
		"Policies":          true,
	},
	"AWS::IAM::Group": {
		"ManagedPolicyArns": true,
		"Policies":          true,
	},
	"AWS::ECS::TaskDefinition": {
		"ContainerDefinitions": true,
		"Volumes":              true,
	},
	"AWS::ECS::Service": {
		"LoadBalancers":           true,
		"ServiceRegistries":       true,
		"PlacementConstraints":    true,
		"PlacementStrategies":     true,
		"CapacityProviderStrategy": true,
	},
	"AWS::Lambda::Function": {
		"VpcConfig.SecurityGroupIds": true,
		"VpcConfig.SubnetIds":        true,
		"Layers":                     true,
	},
	"AWS::DynamoDB::Table": {
		"AttributeDefinitions":    true,
		"KeySchema":               true,
		"GlobalSecondaryIndexes":  true,
		"LocalSecondaryIndexes":   true,
	},
	"AWS::SNS::Topic": {
		"Subscription": true,
	},
	"AWS::CloudWatch::Alarm": {
		"AlarmActions":            true,
		"OKActions":               true,
		"InsufficientDataActions": true,
		"Dimensions":              true,
	},
}

// GetResourceConstraints returns resource-level constraints.
// Returns nil if no constraints are defined.
func GetResourceConstraints(resourceType string) *ResourceConstraints {
	if c, ok := resourceLevelConstraints[resourceType]; ok {
		return &c
	}
	return nil
}

// RequiresUniqueItems returns true if the property requires unique items.
func RequiresUniqueItems(resourceType, propertyName string) bool {
	if props, ok := uniqueItemsProperties[resourceType]; ok {
		return props[propertyName]
	}
	return false
}
