package yamlmerge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"

	"gopkg.in/yaml.v3"
)

// MergeWithComments merges a template YAML (with ${VAR} placeholders) with a vars YAML (with comments)
func MergeWithComments(templateYAML string, varsYAML []byte) (string, error) {
	var varsNode yaml.Node
	if err := yaml.Unmarshal(varsYAML, &varsNode); err != nil {
		return "", fmt.Errorf("parse vars yaml: %w", err)
	}

	flatVars := make(map[string]string)
	comments := make(map[string]string)
	extractVarsWithComments(&varsNode, "", flatVars, comments)

	// Replace variables in template
	var templateNode yaml.Node
	if err := yaml.Unmarshal([]byte(templateYAML), &templateNode); err != nil {
		return "", fmt.Errorf("parse template yaml: %w", err)
	}
	applyReplacements(&templateNode, flatVars, comments)

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(&templateNode); err != nil {
		return "", fmt.Errorf("encode merged yaml: %w", err)
	}
	return buf.String(), nil
}

// extractVarsWithComments flattens YAML and collects comments
func extractVarsWithComments(node *yaml.Node, prefix string, out map[string]string, comments map[string]string) {
	if node.Kind == yaml.MappingNode {
		for i := 0; i < len(node.Content); i += 2 {
			k := node.Content[i]
			v := node.Content[i+1]
			key := k.Value
			if prefix != "" {
				key = prefix + "." + key
			}
			if v.LineComment != "" {
				comments[key] = v.LineComment
			} else if k.HeadComment != "" {
				comments[key] = k.HeadComment
			}
			switch v.Kind {
			case yaml.MappingNode, yaml.SequenceNode:
				extractVarsWithComments(v, key, out, comments)
			case yaml.ScalarNode:
				out[key] = v.Value
			}
		}
	} else if node.Kind == yaml.SequenceNode {
		var list []string
		for i, item := range node.Content {
			key := fmt.Sprintf("%s.%d", prefix, i)
			out[key] = item.Value
			list = append(list, item.Value)
		}
		if prefix != "" {
			j, _ := json.Marshal(list)
			out[prefix] = string(j)
		}
	}
}

// applyReplacements walks the template tree and replaces ${KEY} with values and comments
func applyReplacements(node *yaml.Node, values map[string]string, comments map[string]string) {
	if node.Kind == yaml.ScalarNode {
		node.Value = replaceVars(node.Value, values)
		if cmt, ok := comments[extractKeyName(node.Value)]; ok {
			node.LineComment = cmt
		}
		return
	}
	for _, child := range node.Content {
		applyReplacements(child, values, comments)
	}
}

// replaceVars replaces ${KEY} with its value
func replaceVars(input string, vars map[string]string) string {
	re := regexp.MustCompile(`\$\{([A-Za-z0-9_.]+)\}`)
	return re.ReplaceAllStringFunc(input, func(m string) string {
		key := re.FindStringSubmatch(m)[1]
		if val, ok := vars[key]; ok {
			return val
		}
		return m
	})
}

func extractKeyName(s string) string {
	re := regexp.MustCompile(`^\$\{([A-Za-z0-9_.]+)\}$`)
	if m := re.FindStringSubmatch(s); len(m) == 2 {
		return m[1]
	}
	return ""
}
