package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3039{})
}

// E3039 validates DynamoDB AttributeDefinitions and KeySchema match.
type E3039 struct{}

func (r *E3039) ID() string {
	return "E3039"
}

func (r *E3039) ShortDesc() string {
	return "AttributeDefinitions / KeySchemas mismatch"
}

func (r *E3039) Description() string {
	return "Verifies attribute sets align between definitions and schemas"
}

func (r *E3039) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-dynamodb-table.html"
}

func (r *E3039) Tags() []string {
	return []string{"resources", "dynamodb", "attributes"}
}

func (r *E3039) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::DynamoDB::Table" {
			continue
		}

		// Get AttributeDefinitions
		attrDefsRaw, hasAttrDefs := res.Properties["AttributeDefinitions"]
		if !hasAttrDefs {
			continue
		}

		// Skip intrinsic functions
		if isIntrinsicFunction(attrDefsRaw) {
			continue
		}

		attrDefs, ok := attrDefsRaw.([]any)
		if !ok {
			continue
		}

		// Build set of defined attributes
		definedAttrs := make(map[string]bool)
		for _, attrDefRaw := range attrDefs {
			if isIntrinsicFunction(attrDefRaw) {
				continue
			}
			attrDef, ok := attrDefRaw.(map[string]any)
			if !ok {
				continue
			}
			if attrName, ok := attrDef["AttributeName"].(string); ok {
				definedAttrs[attrName] = true
			}
		}

		// Collect attributes used in key schemas
		usedAttrs := make(map[string]bool)

		// Check KeySchema
		if keySchemaRaw, hasKeySchema := res.Properties["KeySchema"]; hasKeySchema {
			if !isIntrinsicFunction(keySchemaRaw) {
				if keySchema, ok := keySchemaRaw.([]any); ok {
					for _, keyRaw := range keySchema {
						if isIntrinsicFunction(keyRaw) {
							continue
						}
						if key, ok := keyRaw.(map[string]any); ok {
							if attrName, ok := key["AttributeName"].(string); ok {
								usedAttrs[attrName] = true
							}
						}
					}
				}
			}
		}

		// Check GlobalSecondaryIndexes
		if gsiRaw, hasGSI := res.Properties["GlobalSecondaryIndexes"]; hasGSI {
			if !isIntrinsicFunction(gsiRaw) {
				if gsis, ok := gsiRaw.([]any); ok {
					for _, gsiItemRaw := range gsis {
						if isIntrinsicFunction(gsiItemRaw) {
							continue
						}
						if gsi, ok := gsiItemRaw.(map[string]any); ok {
							if keySchemaRaw, hasKeySchema := gsi["KeySchema"]; hasKeySchema {
								if !isIntrinsicFunction(keySchemaRaw) {
									if keySchema, ok := keySchemaRaw.([]any); ok {
										for _, keyRaw := range keySchema {
											if isIntrinsicFunction(keyRaw) {
												continue
											}
											if key, ok := keyRaw.(map[string]any); ok {
												if attrName, ok := key["AttributeName"].(string); ok {
													usedAttrs[attrName] = true
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}

		// Check LocalSecondaryIndexes
		if lsiRaw, hasLSI := res.Properties["LocalSecondaryIndexes"]; hasLSI {
			if !isIntrinsicFunction(lsiRaw) {
				if lsis, ok := lsiRaw.([]any); ok {
					for _, lsiItemRaw := range lsis {
						if isIntrinsicFunction(lsiItemRaw) {
							continue
						}
						if lsi, ok := lsiItemRaw.(map[string]any); ok {
							if keySchemaRaw, hasKeySchema := lsi["KeySchema"]; hasKeySchema {
								if !isIntrinsicFunction(keySchemaRaw) {
									if keySchema, ok := keySchemaRaw.([]any); ok {
										for _, keyRaw := range keySchema {
											if isIntrinsicFunction(keyRaw) {
												continue
											}
											if key, ok := keyRaw.(map[string]any); ok {
												if attrName, ok := key["AttributeName"].(string); ok {
													usedAttrs[attrName] = true
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}

		// Check for mismatches
		// All used attributes must be defined
		for attr := range usedAttrs {
			if !definedAttrs[attr] {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("DynamoDB Table '%s' uses attribute '%s' in KeySchema but it's not defined in AttributeDefinitions", resName, attr),
					Line:    res.Node.Line,
					Column:  res.Node.Column,
					Path:    []string{"Resources", resName, "Properties", "AttributeDefinitions"},
				})
			}
		}

		// All defined attributes should be used
		for attr := range definedAttrs {
			if !usedAttrs[attr] {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("DynamoDB Table '%s' defines attribute '%s' in AttributeDefinitions but it's not used in any KeySchema", resName, attr),
					Line:    res.Node.Line,
					Column:  res.Node.Column,
					Path:    []string{"Resources", resName, "Properties", "AttributeDefinitions"},
				})
			}
		}
	}

	return matches
}
