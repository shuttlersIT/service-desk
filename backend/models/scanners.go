package models

import (
	"fmt"
	"go/types"
	"reflect"
	"strings"
	"time"

	"golang.org/x/tools/go/packages"
)

// ImplementInterfacesForModels generates and implements driver.Valuer and driver.Scanner interfaces

// GetModels retrieves all models/struct names in the models package.
func GetModels() ([]string, error) {
	var models []string

	// Get the package name dynamically
	packageName := reflect.TypeOf(DeviceRegistration{}).PkgPath()

	// Find all types in the models package
	pkg, err := loadPackage(packageName)
	if err != nil {
		return nil, fmt.Errorf("error loading package %s: %v", packageName, err)
	}

	// Iterate through the exported types and add struct names to the models slice
	for _, typeName := range pkg.Types.Scope().Names() {
		obj := pkg.Types.Scope().Lookup(typeName)
		if typ, ok := obj.Type().(*types.Named); ok {
			if _, isStruct := typ.Underlying().(*types.Struct); isStruct {
				models = append(models, typeName)
			}
		}
	}

	return models, nil
}

// AssetModelNames retrieves the names of all models in the models package.
func AssetModelNames() []string {
	return []string{
		"DeviceRegistration",
		"Assets",
		"InventoryItem",
		// Add more model names as needed
	}
}

// loadPackage dynamically loads a package by its name.
func loadPackage(packageName string) (*packages.Package, error) {
	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedImports | packages.NeedSyntax}
	pkgs, err := packages.Load(cfg, packageName)
	if err != nil {
		return nil, err
	}

	if packages.PrintErrors(pkgs) > 0 {
		return nil, fmt.Errorf("error loading package %s", packageName)
	}

	return pkgs[0], nil
}

// getModelType gets the reflect.Type of a model struct by its full name.
func getModelType(fullModelName string) (reflect.Type, error) {
	// Assuming fullModelName is in the format "packagePath.TypeName"
	parts := strings.Split(fullModelName, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid fullModelName format: %s", fullModelName)
	}

	packageName := parts[0]
	typeName := parts[1]

	pkg, err := loadPackage(packageName)
	if err != nil {
		return nil, fmt.Errorf("error loading package %s: %v", packageName, err)
	}

	// Look up the named type within the package
	obj := pkg.Types.Scope().Lookup(typeName)
	if obj == nil {
		return nil, fmt.Errorf("%s not found in package %s", typeName, packageName)
	}

	// Return the reflect type of the named type
	return reflect.TypeOf(obj.Type()), nil
}

// addSliceValuerScanner generates Valuer and Scanner implementations for string slices.
func addSliceValuerScanner(modelName, fieldName string) {
	// Implement the Valuer interface
	eval := fmt.Sprintf(`
func (m %s) %s() (driver.Value, error) {
	return strings.Join(m.%s, ","), nil
}
`, modelName, getValuerMethodName(fieldName), fieldName)

	// Implement the Scanner interface
	escan := fmt.Sprintf(`
func (m *%s) %s(src interface{}) error {
	str, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("Invalid type for %s: expected []byte, got %T", src)
	}
	m.%s = strings.Split(string(str), ",")
	return nil
}
`, modelName, getScannerMethodName(fieldName), modelName, fieldName)

	fmt.Printf("%s\n%s\n", eval, escan)
}

// addStringValuerScanner generates Valuer and Scanner implementations for string fields.
func addStringValuerScanner(modelName, fieldName string) {
	// Implement the Valuer interface
	eval := fmt.Sprintf(`
func (m %s) %s() (driver.Value, error) {
	return m.%s, nil
}
`, modelName, getValuerMethodName(fieldName), fieldName)

	// Implement the Scanner interface
	escan := fmt.Sprintf(`
func (m *%s) %s(src interface{}) error {
	str, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("Invalid type for %s: expected []byte, got %T", src)
	}
	m.%s = string(str)
	return nil
}
`, modelName, getScannerMethodName(fieldName), modelName, fieldName)

	fmt.Printf("%s\n%s\n", eval, escan)
}

// addIntValuerScanner generates Valuer and Scanner implementations for int and int64 fields.
func addIntValuerScanner(modelName, fieldName string) {
	// Implement the Valuer interface
	eval := fmt.Sprintf(`
func (m %s) %s() (driver.Value, error) {
	return int64(m.%s), nil
}
`, modelName, getValuerMethodName(fieldName), fieldName)

	// Implement the Scanner interface
	escan := fmt.Sprintf(`
func (m *%s) %s(src interface{}) error {
	switch v := src.(type) {
	case int64:
		m.%s = int(v)
	case []byte:
		val, err := strconv.ParseInt(string(v), 10, 0)
		if err != nil {
			return err
		}
		m.%s = int(val)
	default:
		return fmt.Errorf("Invalid type for %s: expected int64 or []byte, got %T", src)
	}
	return nil
}
`, modelName, getScannerMethodName(fieldName), modelName, fieldName, modelName)

	fmt.Printf("%s\n%s\n", eval, escan)
}

// addFloat64ValuerScanner generates Valuer and Scanner implementations for float64 fields.
func addFloat64ValuerScanner(modelName, fieldName string) {
	// Implement the Valuer interface
	eval := fmt.Sprintf(`
func (m %s) %s() (driver.Value, error) {
	return m.%s, nil
}
`, modelName, getValuerMethodName(fieldName), fieldName)

	// Implement the Scanner interface
	escan := fmt.Sprintf(`
func (m *%s) %s(src interface{}) error {
	switch v := src.(type) {
	case float64:
		m.%s = v
	case []byte:
		val, err := strconv.ParseFloat(string(v), 64)
		if err != nil {
			return err
		}
		m.%s = val
	default:
		return fmt.Errorf("Invalid type for %s: expected float64 or []byte, got %T", src)
	}
	return nil
}
`, modelName, getScannerMethodName(fieldName), modelName, fieldName, modelName)

	fmt.Printf("%s\n%s\n", eval, escan)
}

// getValuerMethodName generates the method name for the Valuer interface.
func getValuerMethodName(fieldName string) string {
	return fmt.Sprintf("Value%s", fieldName)
}

// getScannerMethodName generates the method name for the Scanner interface.
func getScannerMethodName(fieldName string) string {
	return fmt.Sprintf("Scan%s", fieldName)
}

// addTimeValuerScanner generates Valuer and Scanner implementations for time.Time fields.
func addTimeValuerScanner(modelName, fieldName string) {
	// Implement the Valuer interface
	eval := fmt.Sprintf(`
func (m %s) %s() (driver.Value, error) {
	return m.%s.UTC().Format(time.RFC3339), nil
}
`, modelName, getValuerMethodName(fieldName), fieldName)

	// Implement the Scanner interface
	escan := fmt.Sprintf(`
func (m *%s) %s(src interface{}) error {
	str, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("Invalid type for %s: expected []byte, got %T", src)
	}
	t, err := time.Parse(time.RFC3339, string(str))
	if err != nil {
		return err
	}
	m.%s = t
	return nil
}
`, modelName, getScannerMethodName(fieldName), modelName, fieldName)

	fmt.Printf("%s\n%s\n", eval, escan)
}

// addBoolValuerScanner generates Valuer and Scanner implementations for bool fields.
func addBoolValuerScanner(modelName, fieldName string) {
	// Implement the Valuer interface
	eval := fmt.Sprintf(`
func (m %s) %s() (driver.Value, error) {
	return strconv.FormatBool(m.%s), nil
}
`, modelName, getValuerMethodName(fieldName), fieldName)

	// Implement the Scanner interface
	escan := fmt.Sprintf(`
func (m *%s) %s(src interface{}) error {
	b, err := strconv.ParseBool(string(src.([]byte)))
	if err != nil {
		return err
	}
	m.%s = b
	return nil
}
`, modelName, getScannerMethodName(fieldName), modelName, fieldName)

	fmt.Printf("%s\n%s\n", eval, escan)
}

// addCustomTypeValuerScanner generates Valuer and Scanner implementations for custom types.
func addCustomTypeValuerScanner(modelName, fieldName string) {
	// Implement the Valuer interface
	eval := fmt.Sprintf(`
func (m %s) %s() (driver.Value, error) {
	// Implement conversion logic for %s.%s to string representation
	return customTypeToString(m.%s), nil
}
`, modelName, getValuerMethodName(fieldName), modelName, fieldName, fieldName)

	// Implement the Scanner interface
	escan := fmt.Sprintf(`
func (m *%s) %s(src interface{}) error {
	str, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("Invalid type for %s: expected []byte, got %T", src)
	}
	// Implement conversion logic for string to %s.%s
	m.%s = stringToCustomType(string(str))
	return nil
}
`, modelName, getScannerMethodName(fieldName), modelName, fieldName, fieldName)

	fmt.Printf("%s\n%s\n", eval, escan)
}

// ... (previous code)

// addTypeValuerScanner generates Valuer and Scanner implementations for a given type.
func addTypeValuerScanner(modelName, fieldName, typeName string) {
	// Implement the Valuer interface
	eval := fmt.Sprintf(`
func (m %s) %s() (driver.Value, error) {
	return %s(m.%s), nil
}
`, modelName, getValuerMethodName(fieldName), typeName, fieldName)

	// Implement the Scanner interface
	escan := fmt.Sprintf(`
func (m *%s) %s(src interface{}) error {
	switch v := src.(type) {
	case %s:
		m.%s = %s(v)
	case []byte:
		val, err := %sFromString(string(v))
		if err != nil {
			return err
		}
		m.%s = val
	default:
		return fmt.Errorf("Invalid type for %s: expected %s or []byte, got %T", src)
	}
	return nil
}
`, modelName, getScannerMethodName(fieldName), typeName, fieldName, typeName, typeName, modelName, fieldName, typeName)

	fmt.Printf("%s\n%s\n", eval, escan)
}

// addTypeConversionFunctions generates conversion functions for a given type.
func addTypeConversionFunctions(typeName string) {
	// Implement conversion logic for type to string
	evalToString := fmt.Sprintf(`
func %sToString(val %s) string {
	// Implement conversion logic from %s to string
	return fmt.Sprintf("%%v", val)
}
`, typeName, typeName, typeName)

	// Implement conversion logic for string to type
	evalFromString := fmt.Sprintf(`
func %sFromString(str string) (%s, error) {
	// Implement conversion logic from string to %s
	var val %s
	// Parse the string and set values for %s fields
	// Example: val.Field = parsedValue
	return val, nil
}
`, typeName, typeName, typeName, typeName, typeName)

	fmt.Printf("%s\n%s\n", evalToString, evalFromString)
}

// ImplementInterfacesForModels generates and implements driver.Valuer and driver.Scanner interfaces
func ImplementInterfacesForModels() error {
	models, err := GetModels()
	if err != nil {
		return fmt.Errorf("Error getting models", err)
	}

	for _, modelName := range models {
		fullModelName := fmt.Sprintf("models.%s", modelName)
		modelType, err := getModelType(fullModelName)
		if err != nil {
			return fmt.Errorf("Error getting reflect.Type for %s", fullModelName)
		}

		if err := implementScannerValuerForModel(modelType); err != nil {
			return fmt.Errorf("Error implementing Scanner and Valuer for %s", fullModelName)
		}
	}

	return nil
}

// implementScannerValuerForModel generates Scanner and Valuer implementations for the given model type.
func implementScannerValuerForModel(modelType reflect.Type) error {
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldType := field.Type

		switch fieldType.Kind() {
		case reflect.Slice:
			if fieldType.Elem().Kind() == reflect.String {
				fmt.Printf("Implementing for %s.%s (%s)\n", modelType.Name(), field.Name, fieldType.Name())
				addSliceValuerScanner(modelType.Name(), field.Name)
			}
		case reflect.String:
			fmt.Printf("Implementing for %s.%s (%s)\n", modelType.Name(), field.Name, fieldType.Name())
			addStringValuerScanner(modelType.Name(), field.Name)
		case reflect.Int, reflect.Int64:
			fmt.Printf("Implementing for %s.%s (%s)\n", modelType.Name(), field.Name, fieldType.Name())
			addIntValuerScanner(modelType.Name(), field.Name)
		case reflect.Float64:
			fmt.Printf("Implementing for %s.%s (%s)\n", modelType.Name(), field.Name, fieldType.Name())
			addFloat64ValuerScanner(modelType.Name(), field.Name)
		case reflect.Bool:
			fmt.Printf("Implementing for %s.%s (%s)\n", modelType.Name(), field.Name, fieldType.Name())
			addBoolValuerScanner(modelType.Name(), field.Name)
		case reflect.Struct:
			if fieldType == reflect.TypeOf(time.Time{}) {
				fmt.Printf("Implementing for %s.%s (%s)\n", modelType.Name(), field.Name, fieldType.Name())
				addTimeValuerScanner(modelType.Name(), field.Name)
			}
		default:
			// Handle additional types or custom types
			typeName := fieldType.String()

			switch typeName {
			case "CustomType":
				fmt.Printf("Implementing for %s.%s (%s)\n", modelType.Name(), field.Name, fieldType.Name())
				addCustomTypeValuerScanner(modelType.Name(), field.Name)
			case "AnotherType":
				fmt.Printf("Implementing for %s.%s (%s)\n", modelType.Name(), field.Name, fieldType.Name())
				addTypeValuerScanner(modelType.Name(), field.Name, typeName)
				addTypeConversionFunctions(typeName)
				// Add more cases for other types as needed
			}
		}
	}

	return nil
}

// ... (more code)
