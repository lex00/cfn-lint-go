package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3023{})
}

// E3023 validates Route53 RecordSet configurations.
type E3023 struct{}

func (r *E3023) ID() string {
	return "E3023"
}

func (r *E3023) ShortDesc() string {
	return "Validate Route53 RecordSets"
}

func (r *E3023) Description() string {
	return "Checks RecordSet configurations for proper setup according to DNS standards"
}

func (r *E3023) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53-recordset.html"
}

func (r *E3023) Tags() []string {
	return []string{"resources", "route53", "recordset"}
}

func (r *E3023) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::Route53::RecordSet" {
			continue
		}

		// Must have either ResourceRecords or AliasTarget, but not both
		hasResourceRecords := false
		hasAliasTarget := false

		if _, ok := res.Properties["ResourceRecords"]; ok {
			hasResourceRecords = true
		}
		if _, ok := res.Properties["AliasTarget"]; ok {
			hasAliasTarget = true
		}

		if hasResourceRecords && hasAliasTarget {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("RecordSet '%s' cannot have both ResourceRecords and AliasTarget", resName),
				Line:    res.Node.Line,
				Column:  res.Node.Column,
				Path:    []string{"Resources", resName, "Properties"},
			})
		}

		if !hasResourceRecords && !hasAliasTarget {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("RecordSet '%s' must have either ResourceRecords or AliasTarget", resName),
				Line:    res.Node.Line,
				Column:  res.Node.Column,
				Path:    []string{"Resources", resName, "Properties"},
			})
		}

		// TTL is not allowed with AliasTarget
		if hasAliasTarget {
			if _, hasTTL := res.Properties["TTL"]; hasTTL {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("RecordSet '%s' with AliasTarget cannot have TTL", resName),
					Line:    res.Node.Line,
					Column:  res.Node.Column,
					Path:    []string{"Resources", resName, "Properties", "TTL"},
				})
			}
		}
	}

	return matches
}
