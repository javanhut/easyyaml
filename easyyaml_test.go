package easyyaml

import (
	"os"
	"testing"

	"github.com/javanhut/easyjson"
)

const testYAML = `
name: John Doe
age: 30
email: john@example.com
address:
  street: 123 Main St
  city: New York
  zip: 10001
hobbies:
  - reading
  - swimming
  - coding
settings:
  theme: dark
  notifications: true
  max_items: 100
`

func TestLoads(t *testing.T) {
	yv, err := Loads(testYAML)
	if err != nil {
		t.Fatalf("Failed to load YAML: %v", err)
	}

	if yv.Get("name").AsString() != "John Doe" {
		t.Errorf("Expected name to be 'John Doe', got %s", yv.Get("name").AsString())
	}

	if yv.Get("age").AsInt() != 30 {
		t.Errorf("Expected age to be 30, got %d", yv.Get("age").AsInt())
	}
}

func TestQ(t *testing.T) {
	yv, err := Loads(testYAML)
	if err != nil {
		t.Fatalf("Failed to load YAML: %v", err)
	}

	city := yv.Q("address", "city").AsString()
	if city != "New York" {
		t.Errorf("Expected city to be 'New York', got %s", city)
	}

	hobby := yv.Q("hobbies", 1).AsString()
	if hobby != "swimming" {
		t.Errorf("Expected second hobby to be 'swimming', got %s", hobby)
	}
}

func TestPath(t *testing.T) {
	yv, err := Loads(testYAML)
	if err != nil {
		t.Fatalf("Failed to load YAML: %v", err)
	}

	zip := yv.Path("address.zip").AsInt()
	if zip != 10001 {
		t.Errorf("Expected zip to be 10001, got %d", zip)
	}

	maxItems := yv.Path("settings.max_items").AsInt()
	if maxItems != 100 {
		t.Errorf("Expected max_items to be 100, got %d", maxItems)
	}
}

func TestSetAndGet(t *testing.T) {
	yv := NewObject()
	
	err := yv.Set("name", "Alice")
	if err != nil {
		t.Fatalf("Failed to set name: %v", err)
	}

	if yv.Get("name").AsString() != "Alice" {
		t.Errorf("Expected name to be 'Alice', got %s", yv.Get("name").AsString())
	}
}

func TestAppend(t *testing.T) {
	yv := NewArray()
	
	err := yv.Append("first")
	if err != nil {
		t.Fatalf("Failed to append: %v", err)
	}

	err = yv.Append("second")
	if err != nil {
		t.Fatalf("Failed to append: %v", err)
	}

	if yv.Len() != 2 {
		t.Errorf("Expected length to be 2, got %d", yv.Len())
	}

	if yv.Get(0).AsString() != "first" {
		t.Errorf("Expected first item to be 'first', got %s", yv.Get(0).AsString())
	}
}

func TestTypeChecking(t *testing.T) {
	yv, err := Loads(testYAML)
	if err != nil {
		t.Fatalf("Failed to load YAML: %v", err)
	}

	if !yv.Get("name").IsString() {
		t.Error("Expected name to be a string")
	}

	if !yv.Get("age").IsNumber() {
		t.Error("Expected age to be a number")
	}

	if !yv.Get("address").IsObject() {
		t.Error("Expected address to be an object")
	}

	if !yv.Get("hobbies").IsArray() {
		t.Error("Expected hobbies to be an array")
	}

	if !yv.Get("settings").Get("notifications").IsBool() {
		t.Error("Expected notifications to be a boolean")
	}
}

func TestDumps(t *testing.T) {
	yv := NewObject()
	yv.Set("name", "Bob")
	yv.Set("age", 25)

	yamlStr, err := yv.Dumps()
	if err != nil {
		t.Fatalf("Failed to dumps YAML: %v", err)
	}

	// Parse it back
	yv2, err := Loads(yamlStr)
	if err != nil {
		t.Fatalf("Failed to parse dumped YAML: %v", err)
	}

	if yv2.Get("name").AsString() != "Bob" {
		t.Errorf("Expected name to be 'Bob', got %s", yv2.Get("name").AsString())
	}

	if yv2.Get("age").AsInt() != 25 {
		t.Errorf("Expected age to be 25, got %d", yv2.Get("age").AsInt())
	}
}

func TestClone(t *testing.T) {
	yv, err := Loads(testYAML)
	if err != nil {
		t.Fatalf("Failed to load YAML: %v", err)
	}

	cloned := yv.Clone()
	
	// Modify original
	yv.Set("name", "Modified")
	
	// Clone should be unchanged
	if cloned.Get("name").AsString() != "John Doe" {
		t.Errorf("Expected cloned name to be 'John Doe', got %s", cloned.Get("name").AsString())
	}
}

func TestUpdate(t *testing.T) {
	yv1 := NewObject()
	yv1.Set("name", "Alice")
	yv1.Set("age", 30)

	yv2 := NewObject()
	yv2.Set("age", 25)
	yv2.Set("city", "Boston")

	err := yv1.Update(yv2)
	if err != nil {
		t.Fatalf("Failed to update: %v", err)
	}

	if yv1.Get("name").AsString() != "Alice" {
		t.Errorf("Expected name to remain 'Alice', got %s", yv1.Get("name").AsString())
	}

	if yv1.Get("age").AsInt() != 25 {
		t.Errorf("Expected age to be updated to 25, got %d", yv1.Get("age").AsInt())
	}

	if yv1.Get("city").AsString() != "Boston" {
		t.Errorf("Expected city to be 'Boston', got %s", yv1.Get("city").AsString())
	}
}

func TestJSONConversion(t *testing.T) {
	// Create a JSON value
	jsonStr := `{"name": "John", "age": 30, "hobbies": ["reading", "coding"]}`
	jsonValue, err := easyjson.Loads(jsonStr)
	if err != nil {
		t.Fatalf("Failed to create JSON value: %v", err)
	}

	// Convert to YAML
	yamlValue, err := FromJSON(jsonValue)
	if err != nil {
		t.Fatalf("Failed to convert JSON to YAML: %v", err)
	}

	if yamlValue.Get("name").AsString() != "John" {
		t.Errorf("Expected name to be 'John', got %s", yamlValue.Get("name").AsString())
	}

	if yamlValue.Get("age").AsInt() != 30 {
		t.Errorf("Expected age to be 30, got %d", yamlValue.Get("age").AsInt())
	}

	// Convert back to JSON
	jsonValue2, err := yamlValue.ToJSON()
	if err != nil {
		t.Fatalf("Failed to convert YAML to JSON: %v", err)
	}

	if jsonValue2.Get("name").AsString() != "John" {
		t.Errorf("Expected name to be 'John', got %s", jsonValue2.Get("name").AsString())
	}

	if jsonValue2.Get("age").AsInt() != 30 {
		t.Errorf("Expected age to be 30, got %d", jsonValue2.Get("age").AsInt())
	}
}

func TestKeys(t *testing.T) {
	yv, err := Loads(testYAML)
	if err != nil {
		t.Fatalf("Failed to load YAML: %v", err)
	}

	keys := yv.Keys()
	if len(keys) == 0 {
		t.Error("Expected keys to be non-empty")
	}

	// Check if expected keys exist
	keyMap := make(map[interface{}]bool)
	for _, key := range keys {
		keyMap[key] = true
	}

	expectedKeys := []string{"name", "age", "email", "address", "hobbies", "settings"}
	for _, expectedKey := range expectedKeys {
		if !keyMap[expectedKey] {
			t.Errorf("Expected key '%s' to exist", expectedKey)
		}
	}
}

func TestValues(t *testing.T) {
	yv, err := Loads(testYAML)
	if err != nil {
		t.Fatalf("Failed to load YAML: %v", err)
	}

	hobbies := yv.Get("hobbies")
	values := hobbies.Values()

	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}

	if values[0].AsString() != "reading" {
		t.Errorf("Expected first hobby to be 'reading', got %s", values[0].AsString())
	}
}

func TestDelete(t *testing.T) {
	yv := NewObject()
	yv.Set("name", "Alice")
	yv.Set("age", 30)

	err := yv.Delete("age")
	if err != nil {
		t.Fatalf("Failed to delete key: %v", err)
	}

	if yv.Has("age") {
		t.Error("Expected age key to be deleted")
	}

	if !yv.Has("name") {
		t.Error("Expected name key to still exist")
	}
}

func TestSetPath(t *testing.T) {
	yv := NewObject()
	
	err := yv.SetPath("user.profile.name", "Alice")
	if err != nil {
		t.Fatalf("Failed to set path: %v", err)
	}

	name := yv.Path("user.profile.name").AsString()
	if name != "Alice" {
		t.Errorf("Expected name to be 'Alice', got %s", name)
	}
}

func TestFileOperations(t *testing.T) {
	// Create test data
	yv := NewObject()
	yv.Set("name", "Test User")
	yv.Set("version", "1.0.0")
	
	config := NewObject()
	config.Set("host", "localhost")
	config.Set("port", 8080)
	yv.Set("config", config.Raw())
	
	// Test DumpFile
	filename := "test_output.yaml"
	err := yv.DumpFile(filename)
	if err != nil {
		t.Fatalf("Failed to dump to file: %v", err)
	}
	
	// Test LoadFile
	loaded, err := LoadFile(filename)
	if err != nil {
		t.Fatalf("Failed to load from file: %v", err)
	}
	
	// Verify loaded data
	if loaded.Get("name").AsString() != "Test User" {
		t.Errorf("Expected name to be 'Test User', got %s", loaded.Get("name").AsString())
	}
	
	if loaded.Get("version").AsString() != "1.0.0" {
		t.Errorf("Expected version to be '1.0.0', got %s", loaded.Get("version").AsString())
	}
	
	if loaded.Q("config", "host").AsString() != "localhost" {
		t.Errorf("Expected host to be 'localhost', got %s", loaded.Q("config", "host").AsString())
	}
	
	if loaded.Q("config", "port").AsInt() != 8080 {
		t.Errorf("Expected port to be 8080, got %d", loaded.Q("config", "port").AsInt())
	}
	
	// Clean up
	os.Remove(filename)
}