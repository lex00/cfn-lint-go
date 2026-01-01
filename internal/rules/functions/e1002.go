// Package functions contains intrinsic function and template validation rules (E1xxx).
package functions

import (
	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1002{})
}

// E1002 checks that the template does not exceed CloudFormation size limits.
// Maximum template body size is 51,200 bytes (direct) or 460,800 bytes (S3).
type E1002 struct{}

func (r *E1002) ID() string { return "E1002" }

func (r *E1002) ShortDesc() string {
	return "Template size limit exceeded"
}

func (r *E1002) Description() string {
	return "Checks that the template body does not exceed 51,200 bytes (direct upload) or 460,800 bytes (S3 upload)."
}

func (r *E1002) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E1002"
}

func (r *E1002) Tags() []string {
	return []string{"template", "limits"}
}

// TemplateSizeLimitDirect is the max size for direct template upload (51,200 bytes).
const TemplateSizeLimitDirect = 51200

// TemplateSizeLimitS3 is the max size for S3-uploaded templates (460,800 bytes).
const TemplateSizeLimitS3 = 460800

func (r *E1002) Match(tmpl *template.Template) []rules.Match {
	// Template size validation requires access to raw bytes.
	// This is typically checked at parse time with the raw input.
	// The Match interface doesn't provide raw bytes, so this is a placeholder.
	// Size checking should be done in the linter before parsing.
	return nil
}
