package rulessection

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1702{})
}

// E1702 validates the configuration of RuleCondition.
type E1702 struct{}

func (r *E1702) ID() string { return "E1702" }

func (r *E1702) ShortDesc() string {
	return "Validate RuleCondition configuration"
}

func (r *E1702) Description() string {
	return "Validates that RuleCondition in Rules is properly configured using valid rule-specific intrinsic functions."
}

func (r *E1702) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E1702"
}

func (r *E1702) Tags() []string {
	return []string{"rules", "rulecondition"}
}

// Valid rule-specific intrinsic functions per CloudFormation spec
var validRuleFunctions = map[string]bool{
	"Fn::And":              true,
	"Fn::Contains":         true,
	"Fn::EachMemberEquals": true,
	"Fn::EachMemberIn":     true,
	"Fn::Equals":           true,
	"Fn::If":               true,
	"Fn::Not":              true,
	"Fn::Or":               true,
	"Fn::RefAll":           true,
	"Fn::ValueOf":          true,
	"Fn::ValueOfAll":       true,
}

func (r *E1702) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for ruleName, rule := range tmpl.Rules {
		// Only validate if RuleCondition exists
		if rule.RuleCondition != nil {
			// RuleCondition must be a map with exactly one intrinsic function
			condMap, ok := rule.RuleCondition.(map[string]any)
			if !ok {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Rule '%s' RuleCondition must be an intrinsic function", ruleName),
					Line:    rule.Node.Line,
					Column:  rule.Node.Column,
					Path:    []string{"Rules", ruleName, "RuleCondition"},
				})
				continue
			}

			// Validate that the function used is a valid rule function
			for fnName := range condMap {
				if !validRuleFunctions[fnName] {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Rule '%s' RuleCondition uses invalid function '%s'. Valid functions are: Fn::And, Fn::Contains, Fn::EachMemberEquals, Fn::EachMemberIn, Fn::Equals, Fn::If, Fn::Not, Fn::Or, Fn::RefAll, Fn::ValueOf, Fn::ValueOfAll", ruleName, fnName),
						Line:    rule.Node.Line,
						Column:  rule.Node.Column,
						Path:    []string{"Rules", ruleName, "RuleCondition"},
					})
				}
			}
		}
	}

	return matches
}
