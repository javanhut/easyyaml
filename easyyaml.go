package easyyaml

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/javanhut/easyjson"
	"gopkg.in/yaml.v3"
)

// YAMLValue represents a flexible YAML value that can be any type
type YAMLValue struct {
	data interface{}
}

// Q provides a fluent query interface for chaining access
// Usage: data.Q("name", 0, "hair_color").String()
func (yv *YAMLValue) Q(keys ...interface{}) *YAMLValue {
	current := yv
	for _, key := range keys {
		current = current.Get(key)
		if current.IsNull() {
			break
		}
	}
	return current
}

// New creates a new YAMLValue from any Go value
func New(data interface{}) *YAMLValue {
	return &YAMLValue{data: data}
}

// Loads parses a YAML string and returns a YAMLValue
func Loads(yamlStr string) (*YAMLValue, error) {
	var data interface{}
	err := yaml.Unmarshal([]byte(yamlStr), &data)
	if err != nil {
		return nil, err
	}
	return &YAMLValue{data: data}, nil
}

// Load parses YAML from a byte slice and returns a YAMLValue
func Load(yamlBytes []byte) (*YAMLValue, error) {
	var data interface{}
	err := yaml.Unmarshal(yamlBytes, &data)
	if err != nil {
		return nil, err
	}
	return &YAMLValue{data: data}, nil
}

// LoadFile parses YAML from a file and returns a YAMLValue
func LoadFile(filename string) (*YAMLValue, error) {
	yamlBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return Load(yamlBytes)
}

// Dumps converts the YAMLValue to a YAML string
func (yv *YAMLValue) Dumps() (string, error) {
	bytes, err := yaml.Marshal(yv.data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Dump converts the YAMLValue to YAML bytes
func (yv *YAMLValue) Dump() ([]byte, error) {
	return yaml.Marshal(yv.data)
}

// DumpFile writes the YAMLValue to a file
func (yv *YAMLValue) DumpFile(filename string) error {
	yamlBytes, err := yv.Dump()
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}
	
	err = os.WriteFile(filename, yamlBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	return nil
}

// Get retrieves a value by key (for objects) or index (for arrays)
func (yv *YAMLValue) Get(key interface{}) *YAMLValue {
	switch v := yv.data.(type) {
	case map[string]interface{}:
		if keyStr, ok := key.(string); ok {
			if val, exists := v[keyStr]; exists {
				return &YAMLValue{data: val}
			}
		}
	case map[interface{}]interface{}:
		if val, exists := v[key]; exists {
			return &YAMLValue{data: val}
		}
	case []interface{}:
		if keyInt, ok := key.(int); ok {
			if keyInt >= 0 && keyInt < len(v) {
				return &YAMLValue{data: v[keyInt]}
			}
		}
	}
	return &YAMLValue{data: nil}
}

// Set sets a value by key (for objects) or index (for arrays)
func (yv *YAMLValue) Set(key interface{}, value interface{}) error {
	switch v := yv.data.(type) {
	case map[string]interface{}:
		if keyStr, ok := key.(string); ok {
			v[keyStr] = value
			return nil
		}
		return fmt.Errorf("key must be string for string-keyed map")
	case map[interface{}]interface{}:
		v[key] = value
		return nil
	case []interface{}:
		if keyInt, ok := key.(int); ok {
			if keyInt >= 0 && keyInt < len(v) {
				v[keyInt] = value
				return nil
			}
			return fmt.Errorf("index out of range")
		}
		return fmt.Errorf("key must be int for array")
	default:
		return fmt.Errorf("cannot set on non-object/array type")
	}
}

// Has checks if a key exists (for objects) or index is valid (for arrays)
func (yv *YAMLValue) Has(key interface{}) bool {
	switch v := yv.data.(type) {
	case map[string]interface{}:
		if keyStr, ok := key.(string); ok {
			_, exists := v[keyStr]
			return exists
		}
	case map[interface{}]interface{}:
		_, exists := v[key]
		return exists
	case []interface{}:
		if keyInt, ok := key.(int); ok {
			return keyInt >= 0 && keyInt < len(v)
		}
	}
	return false
}

// Delete removes a key from an object or index from array
func (yv *YAMLValue) Delete(key interface{}) error {
	switch v := yv.data.(type) {
	case map[string]interface{}:
		if keyStr, ok := key.(string); ok {
			delete(v, keyStr)
			return nil
		}
		return fmt.Errorf("key must be string for string-keyed map")
	case map[interface{}]interface{}:
		delete(v, key)
		return nil
	case []interface{}:
		if keyInt, ok := key.(int); ok {
			if keyInt >= 0 && keyInt < len(v) {
				copy(v[keyInt:], v[keyInt+1:])
				v = v[:len(v)-1]
				yv.data = v
				return nil
			}
			return fmt.Errorf("index out of range")
		}
		return fmt.Errorf("key must be int for array")
	default:
		return fmt.Errorf("cannot delete from non-object/array type")
	}
}

// Keys returns all keys for an object
func (yv *YAMLValue) Keys() []interface{} {
	switch v := yv.data.(type) {
	case map[string]interface{}:
		keys := make([]interface{}, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		return keys
	case map[interface{}]interface{}:
		keys := make([]interface{}, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		return keys
	}
	return []interface{}{}
}

// Values returns all values for an object or array
func (yv *YAMLValue) Values() []*YAMLValue {
	switch v := yv.data.(type) {
	case map[string]interface{}:
		values := make([]*YAMLValue, 0, len(v))
		for _, val := range v {
			values = append(values, &YAMLValue{data: val})
		}
		return values
	case map[interface{}]interface{}:
		values := make([]*YAMLValue, 0, len(v))
		for _, val := range v {
			values = append(values, &YAMLValue{data: val})
		}
		return values
	case []interface{}:
		values := make([]*YAMLValue, len(v))
		for i, val := range v {
			values[i] = &YAMLValue{data: val}
		}
		return values
	}
	return []*YAMLValue{}
}

// Items returns key-value pairs for an object
func (yv *YAMLValue) Items() map[interface{}]*YAMLValue {
	items := make(map[interface{}]*YAMLValue)
	switch v := yv.data.(type) {
	case map[string]interface{}:
		for k, val := range v {
			items[k] = &YAMLValue{data: val}
		}
	case map[interface{}]interface{}:
		for k, val := range v {
			items[k] = &YAMLValue{data: val}
		}
	}
	return items
}

// Len returns the length of an array or object
func (yv *YAMLValue) Len() int {
	switch v := yv.data.(type) {
	case map[string]interface{}:
		return len(v)
	case map[interface{}]interface{}:
		return len(v)
	case []interface{}:
		return len(v)
	case string:
		return len(v)
	}
	return 0
}

// IsNull checks if the value is null
func (yv *YAMLValue) IsNull() bool {
	return yv.data == nil
}

// IsObject checks if the value is an object
func (yv *YAMLValue) IsObject() bool {
	switch yv.data.(type) {
	case map[string]interface{}, map[interface{}]interface{}:
		return true
	}
	return false
}

// IsArray checks if the value is an array
func (yv *YAMLValue) IsArray() bool {
	_, ok := yv.data.([]interface{})
	return ok
}

// IsString checks if the value is a string
func (yv *YAMLValue) IsString() bool {
	_, ok := yv.data.(string)
	return ok
}

// IsNumber checks if the value is a number
func (yv *YAMLValue) IsNumber() bool {
	switch yv.data.(type) {
	case float64, int, int64, float32:
		return true
	}
	return false
}

// IsBool checks if the value is a boolean
func (yv *YAMLValue) IsBool() bool {
	_, ok := yv.data.(bool)
	return ok
}

// AsString returns the value as a string
func (yv *YAMLValue) AsString() string {
	if str, ok := yv.data.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", yv.data)
}

// AsInt returns the value as an integer
func (yv *YAMLValue) AsInt() int {
	switch v := yv.data.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return 0
}

// AsFloat returns the value as a float64
func (yv *YAMLValue) AsFloat() float64 {
	switch v := yv.data.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return 0.0
}

// AsBool returns the value as a boolean
func (yv *YAMLValue) AsBool() bool {
	switch v := yv.data.(type) {
	case bool:
		return v
	case string:
		return strings.ToLower(v) == "true"
	case float64:
		return v != 0
	case int:
		return v != 0
	}
	return false
}

// AsArray returns the value as a slice of YAMLValues
func (yv *YAMLValue) AsArray() []*YAMLValue {
	if arr, ok := yv.data.([]interface{}); ok {
		result := make([]*YAMLValue, len(arr))
		for i, v := range arr {
			result[i] = &YAMLValue{data: v}
		}
		return result
	}
	return []*YAMLValue{}
}

// AsObject returns the value as a map of YAMLValues
func (yv *YAMLValue) AsObject() map[interface{}]*YAMLValue {
	result := make(map[interface{}]*YAMLValue)
	switch v := yv.data.(type) {
	case map[string]interface{}:
		for k, val := range v {
			result[k] = &YAMLValue{data: val}
		}
	case map[interface{}]interface{}:
		for k, val := range v {
			result[k] = &YAMLValue{data: val}
		}
	}
	return result
}

// Raw returns the underlying Go value
func (yv *YAMLValue) Raw() interface{} {
	return yv.data
}

// String returns a string representation of the YAMLValue
func (yv *YAMLValue) String() string {
	if str, err := yv.Dumps(); err == nil {
		return str
	}
	return fmt.Sprintf("%v", yv.data)
}

// Append adds a value to an array
func (yv *YAMLValue) Append(value interface{}) error {
	if arr, ok := yv.data.([]interface{}); ok {
		yv.data = append(arr, value)
		return nil
	}
	return fmt.Errorf("cannot append to non-array type")
}

// Extend adds multiple values to an array
func (yv *YAMLValue) Extend(values []interface{}) error {
	if arr, ok := yv.data.([]interface{}); ok {
		yv.data = append(arr, values...)
		return nil
	}
	return fmt.Errorf("cannot extend non-array type")
}

// Update merges another object into this one
func (yv *YAMLValue) Update(other *YAMLValue) error {
	switch obj := yv.data.(type) {
	case map[string]interface{}:
		switch otherObj := other.data.(type) {
		case map[string]interface{}:
			for k, v := range otherObj {
				obj[k] = v
			}
			return nil
		case map[interface{}]interface{}:
			for k, v := range otherObj {
				if keyStr, ok := k.(string); ok {
					obj[keyStr] = v
				}
			}
			return nil
		}
		return fmt.Errorf("can only update with another object")
	case map[interface{}]interface{}:
		switch otherObj := other.data.(type) {
		case map[string]interface{}:
			for k, v := range otherObj {
				obj[k] = v
			}
			return nil
		case map[interface{}]interface{}:
			for k, v := range otherObj {
				obj[k] = v
			}
			return nil
		}
		return fmt.Errorf("can only update with another object")
	}
	return fmt.Errorf("cannot update non-object type")
}

// Clone creates a deep copy of the YAMLValue
func (yv *YAMLValue) Clone() *YAMLValue {
	bytes, err := yaml.Marshal(yv.data)
	if err != nil {
		return &YAMLValue{data: nil}
	}

	var cloned interface{}
	if err := yaml.Unmarshal(bytes, &cloned); err != nil {
		return &YAMLValue{data: nil}
	}

	return &YAMLValue{data: cloned}
}

// Path retrieves a nested value using a dot-separated path
func (yv *YAMLValue) Path(path string) *YAMLValue {
	parts := strings.Split(path, ".")
	current := yv

	for _, part := range parts {
		if part == "" {
			continue
		}

		if index, err := strconv.Atoi(part); err == nil {
			current = current.Get(index)
		} else {
			current = current.Get(part)
		}

		if current.IsNull() {
			break
		}
	}

	return current
}

// SetPath sets a nested value using a dot-separated path
func (yv *YAMLValue) SetPath(path string, value interface{}) error {
	parts := strings.Split(path, ".")
	if len(parts) == 0 {
		return fmt.Errorf("empty path")
	}

	current := yv
	for i, part := range parts[:len(parts)-1] {
		if part == "" {
			continue
		}

		var next *YAMLValue
		if index, err := strconv.Atoi(part); err == nil {
			next = current.Get(index)
		} else {
			next = current.Get(part)
		}

		if next.IsNull() {
			if i+1 < len(parts)-1 {
				if _, err := strconv.Atoi(parts[i+1]); err == nil {
					newArray := make([]interface{}, 0)
					current.Set(part, newArray)
				} else {
					newObj := make(map[interface{}]interface{})
					current.Set(part, newObj)
				}
			} else {
				newObj := make(map[interface{}]interface{})
				current.Set(part, newObj)
			}

			if index, err := strconv.Atoi(part); err == nil {
				next = current.Get(index)
			} else {
				next = current.Get(part)
			}
		}

		current = next
	}

	lastPart := parts[len(parts)-1]
	if index, err := strconv.Atoi(lastPart); err == nil {
		return current.Set(index, value)
	} else {
		return current.Set(lastPart, value)
	}
}

// NewObject creates a new YAMLValue representing an empty object
func NewObject() *YAMLValue {
	return &YAMLValue{data: make(map[interface{}]interface{})}
}

// NewArray creates a new YAMLValue representing an empty array
func NewArray() *YAMLValue {
	return &YAMLValue{data: make([]interface{}, 0)}
}

// NewArrayFrom creates a new YAMLValue array from a slice
func NewArrayFrom(items []interface{}) *YAMLValue {
	return &YAMLValue{data: items}
}

// NewObjectFrom creates a new YAMLValue object from a map
func NewObjectFrom(obj map[interface{}]interface{}) *YAMLValue {
	return &YAMLValue{data: obj}
}

// FromJSON converts an easyjson.JSONValue to a YAMLValue
func FromJSON(jsonValue *easyjson.JSONValue) (*YAMLValue, error) {
	jsonBytes, err := jsonValue.Dump()
	if err != nil {
		return nil, fmt.Errorf("failed to dump JSON: %w", err)
	}

	var yamlData interface{}
	if err := yaml.Unmarshal(jsonBytes, &yamlData); err != nil {
		return nil, fmt.Errorf("failed to parse as YAML: %w", err)
	}

	return &YAMLValue{data: yamlData}, nil
}

// ToJSON converts a YAMLValue to an easyjson.JSONValue
func (yv *YAMLValue) ToJSON() (*easyjson.JSONValue, error) {
	// Convert the raw data directly to JSON using easyjson
	return easyjson.New(yv.data), nil
}