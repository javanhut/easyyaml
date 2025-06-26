# easyYaml

A powerful and intuitive YAML library for Go that brings the simplicity of Python's PyYAML to the Go ecosystem. Built with the same philosophy as [easyJson](https://github.com/javanhut/easyjson), easyYaml provides a fluent interface for working with YAML data and seamless conversion between YAML and JSON formats.

## Features

- 🚀 **PyYAML-like Interface** - Familiar `Loads()`, `Dumps()`, `LoadFile()`, and `DumpFile()` functions
- 🔍 **Fluent Query API** - Chain queries with `Q()` or use dot notation with `Path()`
- 🔄 **JSON Integration** - Seamless conversion between easyYaml and easyJson
- 📁 **File Operations** - Easy reading and writing of YAML files
- 🎯 **Type Safety** - Built-in type checking and conversion methods
- 🛠️ **Full CRUD Operations** - Create, read, update, and delete operations on YAML structures
- 🌳 **Deep Nesting Support** - Handle complex nested structures with ease

## Installation

```bash
go get github.com/javanhut/easyyaml
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "github.com/javanhut/easyyaml"
)

func main() {
    // Parse YAML from string
    yamlStr := `
name: John Doe
age: 30
address:
  city: New York
  zip: 10001
hobbies:
  - reading
  - coding
`
    
    data, err := easyyaml.Loads(yamlStr)
    if err != nil {
        log.Fatal(err)
    }
    
    // Access data fluently
    fmt.Println(data.Get("name").AsString())           // John Doe
    fmt.Println(data.Q("address", "city").AsString())  // New York
    fmt.Println(data.Path("address.zip").AsInt())      // 10001
    fmt.Println(data.Q("hobbies", 0).AsString())       // reading
}
```

## API Reference

### Parsing and Serialization

#### Loading YAML

```go
// From string
data, err := easyyaml.Loads(yamlString)

// From bytes
data, err := easyyaml.Load(yamlBytes)

// From file
data, err := easyyaml.LoadFile("config.yaml")
```

#### Dumping YAML

```go
// To string
yamlStr, err := data.Dumps()

// To bytes
yamlBytes, err := data.Dump()

// To file
err := data.DumpFile("output.yaml")
```

### Data Access

#### Basic Access

```go
// Get value by key
value := data.Get("key")

// Get nested value using Q (query)
value := data.Q("parent", "child", "grandchild")

// Get value using dot notation
value := data.Path("parent.child.grandchild")

// Array access
item := data.Get("items").Get(0)  // First item
item := data.Q("items", 2)         // Third item
```

#### Type Conversion

```go
// Convert to basic types
str := data.Get("name").AsString()
num := data.Get("age").AsInt()
flt := data.Get("price").AsFloat()
bool := data.Get("active").AsBool()

// Get as collections
arr := data.Get("items").AsArray()
obj := data.Get("config").AsObject()

// Get raw value
raw := data.Get("data").Raw()
```

#### Type Checking

```go
if data.Get("field").IsString() { /* ... */ }
if data.Get("field").IsNumber() { /* ... */ }
if data.Get("field").IsObject() { /* ... */ }
if data.Get("field").IsArray() { /* ... */ }
if data.Get("field").IsBool() { /* ... */ }
if data.Get("field").IsNull() { /* ... */ }
```

### Data Manipulation

#### Setting Values

```go
// Set single value
data.Set("key", "value")

// Set nested value
data.SetPath("config.server.port", 8080)

// Update multiple values
updates := easyyaml.NewObject()
updates.Set("version", "2.0")
updates.Set("updated", true)
data.Update(updates)
```

#### Working with Arrays

```go
// Create array
arr := easyyaml.NewArray()

// Append values
arr.Append("first")
arr.Append("second")

// Extend with multiple values
arr.Extend([]interface{}{"third", "fourth"})
```

#### Working with Objects

```go
// Create object
obj := easyyaml.NewObject()
obj.Set("name", "Example")
obj.Set("version", "1.0.0")

// Get all keys
keys := obj.Keys()

// Get all values
values := obj.Values()

// Get key-value pairs
items := obj.Items()
```

### JSON Integration

```go
// Convert easyJson to easyYaml
jsonValue, _ := easyjson.Loads(`{"name": "John", "age": 30}`)
yamlValue, err := easyyaml.FromJSON(jsonValue)

// Convert easyYaml to easyJson
jsonValue, err := yamlValue.ToJSON()
```

## Advanced Examples

### Building Complex YAML Structures

```go
// Create a configuration file structure
config := easyyaml.NewObject()
config.Set("app_name", "MyService")
config.Set("version", "1.0.0")

// Add server configuration
server := easyyaml.NewObject()
server.Set("host", "0.0.0.0")
server.Set("port", 8080)
server.Set("workers", 4)
config.Set("server", server.Raw())

// Add database configuration
db := easyyaml.NewObject()
db.Set("host", "localhost")
db.Set("port", 5432)
db.Set("name", "myapp")

credentials := easyyaml.NewObject()
credentials.Set("username", "admin")
credentials.Set("password", "secret")
db.Set("credentials", credentials.Raw())
config.Set("database", db.Raw())

// Add features list
features := easyyaml.NewArray()
features.Append("authentication")
features.Append("caching")
features.Append("monitoring")
config.Set("features", features.Raw())

// Save to file
err := config.DumpFile("config.yaml")
```

### Working with Existing YAML Files

```go
// Load configuration
config, err := easyyaml.LoadFile("config.yaml")
if err != nil {
    log.Fatal(err)
}

// Modify values
config.SetPath("server.port", 9000)
config.Q("features").Append("logging")

// Check if a key exists
if config.Has("database") {
    // Update database settings
    config.SetPath("database.pool_size", 10)
}

// Delete a key
config.Delete("debug_mode")

// Save changes back
err = config.DumpFile("config.yaml")
```

### Cloning and Merging

```go
// Deep clone
original, _ := easyyaml.LoadFile("config.yaml")
backup := original.Clone()

// Modify original
original.Set("modified", true)

// Backup remains unchanged
fmt.Println(backup.Get("modified").IsNull()) // true

// Merge configurations
base, _ := easyyaml.LoadFile("base.yaml")
override, _ := easyyaml.LoadFile("override.yaml")
base.Update(override) // Merge override into base
```

### Error Handling

```go
// Safe access with null checking
data, err := easyyaml.LoadFile("config.yaml")
if err != nil {
    log.Fatal("Failed to load config:", err)
}

// Check if path exists before accessing
if !data.Path("optional.setting").IsNull() {
    value := data.Path("optional.setting").AsString()
    fmt.Println("Setting:", value)
}

// Type-safe access
age := data.Get("age")
if age.IsNumber() {
    fmt.Printf("Age: %d\n", age.AsInt())
} else {
    fmt.Println("Age is not a valid number")
}
```

## PyYAML Compatibility

easyYaml is designed to feel familiar to Python developers who use PyYAML:

| PyYAML (Python) | easyYaml (Go) |
|-----------------|---------------|
| `yaml.load(string)` | `easyyaml.Loads(string)` |
| `yaml.load(file)` | `easyyaml.LoadFile(filename)` |
| `yaml.dump(data)` | `data.Dumps()` |
| `yaml.dump(data, file)` | `data.DumpFile(filename)` |
| `data['key']` | `data.Get("key")` |
| `data['key'] = value` | `data.Set("key", value)` |

## Best Practices

1. **Always check for errors** when loading YAML:
   ```go
   data, err := easyyaml.LoadFile("config.yaml")
   if err != nil {
       log.Fatal(err)
   }
   ```

2. **Use type checking** before type conversion:
   ```go
   if data.Get("port").IsNumber() {
       port := data.Get("port").AsInt()
   }
   ```

3. **Prefer Path() for deeply nested access**:
   ```go
   // Instead of: data.Get("a").Get("b").Get("c")
   value := data.Path("a.b.c")
   ```

4. **Use Q() for mixed key types**:
   ```go
   // When accessing arrays within objects
   item := data.Q("users", 0, "name")
   ```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Related Projects

- [easyJson](https://github.com/javanhut/easyjson) - Sister library for JSON manipulation
- [easyHttp](https://github.com/javanhut/easyhttp) - HTTP client with easyJson integration