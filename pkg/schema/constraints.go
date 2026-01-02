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
