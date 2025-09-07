package yamlmerge

import (
	"bytes"
	"fmt"
	"regexp"

	"gopkg.in/yaml.v3"
)

// MergeWithRegex 替换模板中的 ${VAR}，保留注释和顺序
func MergeWithRegex(template string, varsYAML []byte) (string, error) {
	var templateNode yaml.Node
	if err := yaml.Unmarshal([]byte(template), &templateNode); err != nil {
		return "", fmt.Errorf("parse template yaml: %w", err)
	}

	var varsNode yaml.Node
	if err := yaml.Unmarshal(varsYAML, &varsNode); err != nil {
		return "", fmt.Errorf("parse vars yaml: %w", err)
	}

	flatVars := make(map[string]*yaml.Node)
	if len(varsNode.Content) > 0 {
		flattenVars(varsNode.Content[0], "", flatVars)
	}

	// 递归替换模板 Node
	if err := replaceNode(templateNode.Content[0], flatVars); err != nil {
		return "", err
	}

	// 序列化回 YAML，保留注释
	out := &bytes.Buffer{}
	enc := yaml.NewEncoder(out)
	enc.SetIndent(2)
	if err := enc.Encode(templateNode.Content[0]); err != nil {
		return "", fmt.Errorf("encode yaml: %w", err)
	}

	return out.String(), nil
}

// replaceNode 递归替换节点中 ${VAR} 占位符
func replaceNode(node *yaml.Node, vars map[string]*yaml.Node) error {
	re := regexp.MustCompile(`\$\{([A-Za-z0-9_.]+)\}`)

	switch node.Kind {
	case yaml.ScalarNode:
		matches := re.FindStringSubmatch(node.Value)
		if len(matches) == 2 {
			key := matches[1]
			valNode, ok := vars[key]
			if !ok {
				return nil
			}
			// 用 vars 节点替换当前节点
			node.Kind = valNode.Kind
			node.Tag = valNode.Tag
			node.Value = valNode.Value
			node.Content = valNode.Content
			node.HeadComment = valNode.HeadComment
			node.LineComment = valNode.LineComment
		}
	case yaml.MappingNode, yaml.SequenceNode:
		for _, child := range node.Content {
			if err := replaceNode(child, vars); err != nil {
				return err
			}
		}
	}
	return nil
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
		if prefix != "" {
			out[prefix] = node
		}
	}
}
