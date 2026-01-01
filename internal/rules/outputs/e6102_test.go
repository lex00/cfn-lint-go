package outputs

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE6102_StringExportName(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
Outputs:
  BucketName:
    Value: !Ref MyBucket
    Export:
      Name: "my-export-name"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E6102{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for string export name, got %d: %v", len(matches), matches)
	}
}

func TestE6102_RefExportName(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  ExportName:
    Type: String
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
Outputs:
  BucketName:
    Value: !Ref MyBucket
    Export:
      Name: !Ref ExportName
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E6102{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for Ref export name, got %d: %v", len(matches), matches)
	}
}

func TestE6102_SubExportName(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
Outputs:
  BucketName:
    Value: !Ref MyBucket
    Export:
      Name: !Sub "${AWS::StackName}-bucket"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E6102{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for Sub export name, got %d: %v", len(matches), matches)
	}
}

func TestE6102_JoinExportName(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
Outputs:
  BucketName:
    Value: !Ref MyBucket
    Export:
      Name: !Join ["-", [!Ref "AWS::StackName", "bucket"]]
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E6102{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for Join export name, got %d: %v", len(matches), matches)
	}
}

func TestE6102_IfExportName(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Conditions:
  IsProd: !Equals [!Ref "AWS::StackName", "prod"]
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
Outputs:
  BucketName:
    Value: !Ref MyBucket
    Export:
      Name: !If [IsProd, "prod-export", "dev-export"]
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E6102{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for If export name, got %d: %v", len(matches), matches)
	}
}

func TestE6102_NumberExportName(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
Outputs:
  BucketName:
    Value: !Ref MyBucket
    Export:
      Name: 12345
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E6102{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for number export name, got %d", len(matches))
	}
}

func TestE6102_BooleanExportName(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
Outputs:
  BucketName:
    Value: !Ref MyBucket
    Export:
      Name: true
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E6102{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for boolean export name, got %d", len(matches))
	}
}

func TestE6102_NoExport(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
Outputs:
  BucketName:
    Value: !Ref MyBucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E6102{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for output without Export, got %d: %v", len(matches), matches)
	}
}

func TestE6102_MultipleOutputs(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
Outputs:
  ValidExport:
    Value: !Ref MyBucket
    Export:
      Name: "valid-name"
  InvalidExport1:
    Value: !Ref MyBucket
    Export:
      Name: 123
  InvalidExport2:
    Value: !Ref MyBucket
    Export:
      Name: false
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E6102{}
	matches := rule.Match(parsed)

	if len(matches) != 2 {
		t.Errorf("Expected 2 matches for invalid export names, got %d", len(matches))
	}
}
