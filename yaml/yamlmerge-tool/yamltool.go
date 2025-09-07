package yamlmerge

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// MergeWithRegex 替换模板中的 ${VAR}，只保留 vars 的行内注释
func MergeWithRegex(template string, varsYAML []byte) (string, error) {
	var varsNode yaml.Node
	if err := yaml.Unmarshal(varsYAML, &varsNode); err != nil {
		return "", fmt.Errorf("parse vars yaml: %w", err)
	}

	flatVars := make(map[string]*yaml.Node)
	if len(varsNode.Content) > 0 {
		flattenVars(varsNode.Content[0], "", flatVars)
	}

	re := regexp.MustCompile(`\$\{([A-Za-z0-9_.]+)\}`)
	lines := strings.Split(template, "\n")
	for i, line := range lines {
		lines[i] = re.ReplaceAllStringFunc(line, func(s string) string {
			key := re.FindStringSubmatch(s)[1]
			valNode, ok := flatVars[key]
			if !ok {
				return s
			}

			indent := getLineIndent(line)
			var valStr string
			if valNode.Kind == yaml.ScalarNode {
				valStr = nodeToYAMLStringWithLineComment(valNode, indent)
				return valStr
			} else {
				valStr = nodeToYAMLStringOnlyValue(valNode, indent+1)
				return "\n" + valStr
			}
		})
	}

	return strings.Join(lines, "\n"), nil
}

// flattenVars 扁平化 vars Node
func flattenVars(node *yaml.Node, prefix string, out map[string]*yaml.Node) {
	switch node.Kind {
	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			k := node.Content[i]
			v := node.Content[i+1]
			key := k.Value
			if prefix != "" {
				key = prefix + "." + key
			}
			out[key] = v
			flattenVars(v, key, out)
		}
	case yaml.SequenceNode:
		for i, item := range node.Content {
			key := fmt.Sprintf("%s.%d", prefix, i)
			out[key] = item
			flattenVars(item, key, out)
		}
	case yaml.ScalarNode:
		out[prefix] = node
	}
}

// 标量节点 + 行内注释
func nodeToYAMLStringWithLineComment(node *yaml.Node, indent int) string {
	indentStr := strings.Repeat("  ", indent)
	val := node.Value
	if node.LineComment != "" {
		val += " " + node.LineComment
	}
	return indentStr + val
}

// 序列或映射只序列化值部分，不包含模板 key
func nodeToYAMLStringOnlyValue(node *yaml.Node, indent int) string {
	indentStr := strings.Repeat("  ", indent)
	var buf bytes.Buffer

	switch node.Kind {
	case yaml.ScalarNode:
		val := node.Value
		if node.LineComment != "" {
			val += " " + node.LineComment
		}
		buf.WriteString(indentStr + val)
	case yaml.SequenceNode:
		for _, item := range node.Content {
			if item.Kind == yaml.ScalarNode {
				buf.WriteString(indentStr + "- " + item.Value)
				if item.LineComment != "" {
					buf.WriteString(" " + item.LineComment)
				}
				buf.WriteString("\n")
			} else {
				// 嵌套 Sequence/Mapping
				buf.WriteString(indentStr + "- " + "\n" + nodeToYAMLStringOnlyValue(item, indent+1) + "\n")
			}
		}
	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			k := node.Content[i]
			v := node.Content[i+1]
			buf.WriteString(fmt.Sprintf("%s%s: %s\n", indentStr, k.Value, nodeToYAMLStringOnlyValue(v, indent+1)))
		}
	}

	return strings.TrimRight(buf.String(), "\r\n")
}

// 获取行缩进空格数（每2个空格为一级）
func getLineIndent(line string) int {
	count := 0
	for _, ch := range line {
		if ch == ' ' {
			count++
		} else {
			break
		}
	}
	return count / 2
}
