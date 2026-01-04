package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/sam"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W3100{})
}

// W3100 warns about SAM Function resources missing MemorySize property.
type W3100 struct{}

func (r *W3100) ID() string { return "W3100" }

func (r *W3100) ShortDesc() string {
	return "SAM Function missing MemorySize"
}

func (r *W3100) Description() string {
	return "Warns when an AWS::Serverless::Function resource does not specify MemorySize. The default is 128 MB which may not be optimal for your function. Consider explicitly setting MemorySize based on your function's requirements."
}

func (r *W3100) Source() string {
	return "https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-function.html"
}

func (r *W3100) Tags() []string {
	return []string{"resources", "sam", "serverless", "lambda"}
}

func (r *W3100) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check if Globals has MemorySize set (via Metadata parsing)
	hasGlobalMemorySize := checkGlobalsFunctionProperty(tmpl, "MemorySize")

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::Serverless::Function" {
			continue
		}

		// Skip if function or Globals has MemorySize
		if _, hasLocal := res.Properties["MemorySize"]; hasLocal || hasGlobalMemorySize {
			continue
		}

		line, column := 0, 0
		if res.Node != nil {
			line = res.Node.Line
			column = res.Node.Column
		}

		matches = append(matches, rules.Match{
			Message: fmt.Sprintf("SAM Function '%s' does not specify MemorySize. Default is 128 MB which may not be optimal. Consider setting an explicit value.", resName),
			Line:    line,
			Column:  column,
			Path:    []string{"Resources", resName, "Properties"},
		})
	}

	return matches
}

// checkGlobalsFunctionProperty checks if a property is set in the Globals.Function section.
// Globals is stored in Metadata as "Globals" in SAM templates.
func checkGlobalsFunctionProperty(tmpl *template.Template, property string) bool {
	if tmpl == nil || !sam.IsSAMTemplate(tmpl) {
		return false
	}

	// Check the root YAML for Globals section
	if tmpl.Root == nil || len(tmpl.Root.Content) == 0 {
		return false
	}

	doc := tmpl.Root.Content[0]
	for i := 0; i < len(doc.Content); i += 2 {
		if doc.Content[i].Value == "Globals" {
			globalsNode := doc.Content[i+1]
			// Look for Function in Globals
			for j := 0; j < len(globalsNode.Content); j += 2 {
				if globalsNode.Content[j].Value == "Function" {
					funcNode := globalsNode.Content[j+1]
					// Look for the property
					for k := 0; k < len(funcNode.Content); k += 2 {
						if funcNode.Content[k].Value == property {
							return true
						}
					}
				}
			}
		}
	}

	return false
}
