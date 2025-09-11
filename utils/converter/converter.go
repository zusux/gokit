package converter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/shopspring/decimal"
)

// ToString 将任意类型转换为字符串表示，优化 float32 处理
func ToString(value interface{}) string {
	if value == nil {
		return "nil"
	}
	// 类型断言处理常见类型
	switch v := value.(type) {
	// 基本类型
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		// 直接处理 float32，避免转换为 float64
		return decimal.NewFromFloat32(v).StringFixed(6) // 6 位小数，符合 float32 精度
	case float64:
		return decimal.NewFromFloat(v).StringFixed(6) // 6 位小数，符合 float64 精度
	case complex64:
		return fmt.Sprintf("%v", v)
	case complex128:
		return fmt.Sprintf("%v", v)
	case []byte:
		return string(v)
	case fmt.Stringer:
		return v.String()
	default:
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false) // ✅ 禁用默认转义
		if err := enc.Encode(value); err != nil {
			log.Println("json marshal failed: %w", err)
			return ""
		}
		return buf.String()
	}
}

// ToJSON 将任意类型转换为 JSON 字符串
func ToJSON(value interface{}) (string, error) {
	if value == nil {
		return "null", nil
	}
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false) // ✅ 禁用默认转义
	if err := enc.Encode(value); err != nil {
		log.Println("json marshal failed: %w", err)
		return "", fmt.Errorf("json marshal failed: %w", err)
	}
	return buf.String(), nil
}
