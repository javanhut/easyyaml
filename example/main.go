package main

import (
	"fmt"
	"log"
	"os"

	"github.com/javanhut/easyjson"
	"github.com/javanhut/easyyaml"
)

func main() {
	// Example YAML data
	yamlData := `
name: John Doe
age: 30
contact:
  email: john@example.com
  phone: "+1-555-0123"
address:
  street: 123 Main St
  city: New York
  state: NY
  zip: 10001
hobbies:
  - reading
  - swimming
  - coding
preferences:
  theme: dark
  notifications: true
  max_items: 100
`

	fmt.Println("=== easyYAML Demo ===\n")

	// Parse YAML
	yv, err := easyyaml.Loads(yamlData)
	if err != nil {
		log.Fatal("Failed to parse YAML:", err)
	}

	// Basic access
	fmt.Printf("Name: %s\n", yv.Get("name").AsString())
	fmt.Printf("Age: %d\n", yv.Get("age").AsInt())

	// Fluent query interface
	fmt.Printf("Email: %s\n", yv.Q("contact", "email").AsString())
	fmt.Printf("City: %s\n", yv.Q("address", "city").AsString())
	fmt.Printf("First hobby: %s\n", yv.Q("hobbies", 0).AsString())

	// Path-based access
	fmt.Printf("Phone: %s\n", yv.Path("contact.phone").AsString())
	fmt.Printf("ZIP: %d\n", yv.Path("address.zip").AsInt())
	fmt.Printf("Theme: %s\n", yv.Path("preferences.theme").AsString())

	// Type checking
	fmt.Printf("Is age a number? %t\n", yv.Get("age").IsNumber())
	fmt.Printf("Is address an object? %t\n", yv.Get("address").IsObject())
	fmt.Printf("Is hobbies an array? %t\n", yv.Get("hobbies").IsArray())

	// Working with arrays
	hobbies := yv.Get("hobbies")
	fmt.Printf("Number of hobbies: %d\n", hobbies.Len())
	for i := 0; i < hobbies.Len(); i++ {
		fmt.Printf("  - %s\n", hobbies.Get(i).AsString())
	}

	// Modifying data
	fmt.Println("\n=== Modifying Data ===")
	yv.Set("age", 31)
	yv.Path("address.city").AsString() // This won't modify, let's use SetPath
	yv.SetPath("address.city", "Boston")
	yv.Get("hobbies").Append("photography")

	fmt.Printf("Updated age: %d\n", yv.Get("age").AsInt())
	fmt.Printf("Updated city: %s\n", yv.Path("address.city").AsString())
	fmt.Printf("Updated hobbies count: %d\n", yv.Get("hobbies").Len())

	// Convert to YAML string
	fmt.Println("\n=== YAML Output ===")
	yamlOutput, err := yv.Dumps()
	if err != nil {
		log.Fatal("Failed to serialize YAML:", err)
	}
	fmt.Println(yamlOutput)

	// Creating new YAML structures
	fmt.Println("\n=== Creating New Structures ===")
	
	// Create a new object
	person := easyyaml.NewObject()
	person.Set("name", "Alice Smith")
	person.Set("age", 28)
	
	// Create nested object
	contact := easyyaml.NewObject()
	contact.Set("email", "alice@example.com")
	contact.Set("phone", "+1-555-0456")
	person.Set("contact", contact.Raw())
	
	// Create array
	skills := easyyaml.NewArray()
	skills.Append("Python")
	skills.Append("Go")
	skills.Append("JavaScript")
	person.Set("skills", skills.Raw())

	personYAML, _ := person.Dumps()
	fmt.Println("New person YAML:")
	fmt.Println(personYAML)

	// JSON to YAML conversion
	fmt.Println("\n=== JSON ↔ YAML Conversion ===")
	
	// Create JSON data
	jsonStr := `{"product": "laptop", "price": 999.99, "specs": {"cpu": "Intel i7", "ram": "16GB"}}`
	jsonValue, err := easyjson.Loads(jsonStr)
	if err != nil {
		log.Fatal("Failed to parse JSON:", err)
	}
	
	// Convert JSON to YAML
	yamlFromJson, err := easyyaml.FromJSON(jsonValue)
	if err != nil {
		log.Fatal("Failed to convert JSON to YAML:", err)
	}
	
	yamlStr, _ := yamlFromJson.Dumps()
	fmt.Println("JSON converted to YAML:")
	fmt.Println(yamlStr)
	
	// Convert back to JSON
	jsonFromYaml, err := yamlFromJson.ToJSON()
	if err != nil {
		log.Fatal("Failed to convert YAML to JSON:", err)
	}
	
	jsonStr2, _ := jsonFromYaml.Dumps()
	fmt.Println("YAML converted back to JSON:")
	fmt.Println(jsonStr2)

	// Demonstrate PyYAML-like features
	fmt.Println("\n=== PyYAML-like Features ===")
	
	// Load from string (like yaml.load)
	config, _ := easyyaml.Loads(`
database:
  host: localhost
  port: 5432
  name: myapp
  credentials:
    username: admin
    password: secret
`)
	
	// Access nested values easily
	dbHost := config.Q("database", "host").AsString()
	dbPort := config.Q("database", "port").AsInt()
	dbUser := config.Q("database", "credentials", "username").AsString()
	
	fmt.Printf("Database: %s:%d (user: %s)\n", dbHost, dbPort, dbUser)
	
	// Dump back (like yaml.dump)
	configYAML, _ := config.Dumps()
	fmt.Println("\nConfig as YAML:")
	fmt.Println(configYAML)

	// File operations
	fmt.Println("\n=== File Operations ===")
	
	// Save to file (like yaml.dump(data, file))
	err = config.DumpFile("config.yaml")
	if err != nil {
		log.Fatal("Failed to save YAML file:", err)
	}
	fmt.Println("Saved config to config.yaml")
	
	// Load from file (like yaml.load(file))
	loadedConfig, err := easyyaml.LoadFile("config.yaml")
	if err != nil {
		log.Fatal("Failed to load YAML file:", err)
	}
	
	// Verify loaded data
	fmt.Printf("Loaded from file - Database name: %s\n", 
		loadedConfig.Q("database", "name").AsString())
	
	// Create a more complex YAML file
	appConfig := easyyaml.NewObject()
	appConfig.Set("app_name", "MyApp")
	appConfig.Set("version", "2.1.0")
	appConfig.Set("environment", "production")
	
	server := easyyaml.NewObject()
	server.Set("host", "0.0.0.0")
	server.Set("port", 8080)
	server.Set("workers", 4)
	appConfig.Set("server", server.Raw())
	
	features := easyyaml.NewArray()
	features.Append("authentication")
	features.Append("caching")
	features.Append("monitoring")
	appConfig.Set("features", features.Raw())
	
	// Save app config to file
	err = appConfig.DumpFile("app_config.yaml")
	if err != nil {
		log.Fatal("Failed to save app config:", err)
	}
	fmt.Println("\nSaved app_config.yaml with:")
	appYAML, _ := appConfig.Dumps()
	fmt.Println(appYAML)
	
	// Clean up example files
	os.Remove("config.yaml")
	os.Remove("app_config.yaml")
}