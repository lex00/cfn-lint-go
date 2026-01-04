package errors

import (
	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E0010{})
}

// E0010 documents SAM transform failures.
// SAM transform errors are actually caught at transform time by the linter,
// but this rule exists for documentation and rule ID tracking.
type E0010 struct{}

func (r *E0010) ID() string { return "E0010" }

func (r *E0010) ShortDesc() string {
	return "SAM transform failed"
}

func (r *E0010) Description() string {
	return "Checks that SAM templates can be successfully transformed to CloudFormation. SAM transform errors include invalid SAM resource configurations, missing required properties, or malformed intrinsic functions."
}

func (r *E0010) Source() string {
	return "https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/what-is-sam.html"
}

func (r *E0010) Tags() []string {
	return []string{"base", "sam", "transform"}
}

func (r *E0010) Match(tmpl *template.Template) []rules.Match {
	// SAM transform errors are caught at transform time in pkg/lint/lint.go.
	// This rule exists for documentation purposes and to provide a rule ID
	// for SAM transform failures reported by the linter.
	return nil
}
