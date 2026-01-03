// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3047{})
}

// E3047 validates ECS Fargate CPU and memory combinations.
type E3047 struct{}

func (r *E3047) ID() string { return "E3047" }

func (r *E3047) ShortDesc() string {
	return "ECS Fargate CPU/memory combinations"
}

func (r *E3047) Description() string {
	return "Validates that ECS Fargate tasks use valid combinations of CPU and memory values."
}

func (r *E3047) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3047"
}

func (r *E3047) Tags() []string {
	return []string{"resources", "properties", "ecs", "fargate", "task"}
}

// Valid Fargate CPU/Memory combinations
var validFargateCombinations = map[string][]string{
	"256":   {"512", "1024", "2048"},
	"512":   {"1024", "2048", "3072", "4096"},
	"1024":  {"2048", "3072", "4096", "5120", "6144", "7168", "8192"},
	"2048":  {"4096", "5120", "6144", "7168", "8192", "9216", "10240", "11264", "12288", "13312", "14336", "15360", "16384"},
	"4096":  {"8192", "9216", "10240", "11264", "12288", "13312", "14336", "15360", "16384", "17408", "18432", "19456", "20480", "21504", "22528", "23552", "24576", "25600", "26624", "27648", "28672", "29696", "30720"},
	"8192":  {"16384", "20480", "24576", "28672", "32768", "36864", "40960", "45056", "49152", "53248", "57344", "61440"},
	"16384": {"32768", "40960", "49152", "57344", "65536", "73728", "81920", "90112", "98304", "106496", "114688", "122880"},
}

func (r *E3047) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::ECS::TaskDefinition" {
			continue
		}

		// Check if it requires Fargate compatibility
		reqCompat, hasReqCompat := res.Properties["RequiresCompatibilities"]
		isFargate := false

		if hasReqCompat {
			if compatList, ok := reqCompat.([]interface{}); ok {
				for _, compat := range compatList {
					if compatStr, ok := compat.(string); ok && compatStr == "FARGATE" {
						isFargate = true
						break
					}
				}
			}
		}

		if !isFargate {
			continue
		}

		cpu, hasCPU := res.Properties["Cpu"]
		memory, hasMemory := res.Properties["Memory"]

		if !hasCPU || !hasMemory {
			continue
		}

		cpuStr, ok1 := cpu.(string)
		memoryStr, ok2 := memory.(string)

		if !ok1 || !ok2 {
			continue
		}

		// Validate combination
		validMemories, cpuValid := validFargateCombinations[cpuStr]
		if !cpuValid {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': Invalid Fargate CPU value '%s'. Valid values: 256, 512, 1024, 2048, 4096, 8192, 16384",
					resName, cpuStr,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties", "Cpu"},
			})
			continue
		}

		// Check if memory is valid for this CPU
		memoryValid := false
		for _, validMem := range validMemories {
			if validMem == memoryStr {
				memoryValid = true
				break
			}
		}

		if !memoryValid {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': Invalid Fargate memory value '%s' for CPU '%s'. Valid values: %v",
					resName, memoryStr, cpuStr, validMemories,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties", "Memory"},
			})
		}
	}

	return matches
}
