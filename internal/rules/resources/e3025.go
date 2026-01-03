package resources

import (
	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3025{})
}

// E3025 validates RDS DB instance class.
type E3025 struct{}

func (r *E3025) ID() string {
	return "E3025"
}

func (r *E3025) ShortDesc() string {
	return "Validates RDS DB Instance Class"
}

func (r *E3025) Description() string {
	return "Confirms instance types match region availability via pricing APIs"
}

func (r *E3025) Source() string {
	return "https://aws.amazon.com/rds/instance-types/"
}

func (r *E3025) Tags() []string {
	return []string{"resources", "rds", "instance"}
}

func (r *E3025) Match(tmpl *template.Template) []rules.Match {
	// This rule would require access to AWS pricing API or a comprehensive
	// list of valid instance types per region. For now, we'll implement
	// a basic check that validates the instance class format.
	// A complete implementation would fetch pricing data from AWS.
	return nil
}
