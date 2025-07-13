package converter

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
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
		return strconv.FormatFloat(float64(v), 'f', 6, 32) // 6 位小数，符合 float32 精度
	case float64:
		return strconv.FormatFloat(v, 'f', 6, 64) // 6 位小数，符合 float64 精度
	case complex64:
		return fmt.Sprintf("%v", v)
	case complex128:
		return fmt.Sprintf("%v", v)
	case []byte:
		return string(v)
	case fmt.Stringer:
		return v.String()

	// 使用反射处理其他类型
	default:
		val := reflect.ValueOf(value)
		switch val.Kind() {
		case reflect.Ptr:
			if val.IsNil() {
				return "nil"
			}
			return ToString(val.Elem().Interface()) // 解引用指针
		case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct:
			return ToString(value)
		case reflect.Interface:
			if val.IsNil() {
				return "nil"
			}
			return ToString(val.Elem().Interface())
		default:
			return fmt.Sprintf("%v", value)
		}
	}
}

// ToJSON 将任意类型转换为 JSON 字符串
func ToJSON(value interface{}) (string, error) {
	if value == nil {
		return "null", nil
	}

	// 使用反射检查类型
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Func, reflect.Chan, reflect.UnsafePointer:
		return "", fmt.Errorf("type %s is not JSON-serializable", val.Kind())
	case reflect.Ptr:
		if val.IsNil() {
			return "null", nil
		}
		return ToJSON(val.Elem().Interface())
	case reflect.Interface:
		if val.IsNil() {
			return "null", nil
		}
		return ToJSON(val.Elem().Interface())
	case reflect.Complex64, reflect.Complex128:
		return "", fmt.Errorf("complex types are not JSON-serializable")
	}

	// 尝试 JSON 序列化
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("json marshal failed: %w", err)
	}
	return string(jsonBytes), nil
}
