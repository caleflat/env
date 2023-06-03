package env

import (
	"os"
	"reflect"
	"strconv"
)

const (
	// DefaultTag is the default tag name used for struct tags.
	DefaultTag = "env"
)

// Parse takes a struct and parses the environment variables into it.
// It uses the `env` tag on the struct fields to determine the environment
// variable name.
//
// Example:
//
//	type Config struct {
//	  Port int `env:"PORT"`
//	}
//
//	var config Config
//	env.Parse(&config)
//
//	fmt.Println(config.Port)
//
// If the environment variable is not present, the field value is not modified.
// If the environment variable is present, but the field cannot be set, an error
// is returned.
func Parse(config interface{}) error {
	return parse(config, "")
}

func parse(config interface{}, prefix string) error {
	if prefix != "" {
		prefix += "_"
	}

	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if value.Kind() == reflect.Struct {
			parse(value.Addr().Interface(), field.Tag.Get(DefaultTag))
		} else {
			env := field.Tag.Get(DefaultTag)
			if env == "" {
				continue
			}

			if err := setField(value, env); err != nil {
				return err
			}
		}
	}

	return nil
}

// setField sets the value of the field to the environment variable.
// If the environment variable is not present, the field value is not modified.
// If the environment variable is present, but the field cannot be set, an error
// is returned.
func setField(value reflect.Value, env string) error {
	if !value.CanSet() {
		return nil
	}

	if !value.IsValid() {
		return nil
	}

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.String:
		value.SetString(GetString(env, ""))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value.SetInt(GetInt64(env, 0))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value.SetUint(uint64(GetInt64(env, 0)))
	case reflect.Bool:
		value.SetBool(GetBool(env, false))
	case reflect.Float32, reflect.Float64:
		value.SetFloat(GetFloat64(env, 0))
	}

	return nil
}

// GetString returns the value of the environment variable named by the key.
// If the variable is not present in the environment, the default value is returned.
func GetString(key string, defaultValue string) string {
	value := Get(key)
	if value == "" {
		return defaultValue
	}

	return value
}

// GetInt returns the value of the environment variable named by the key.
// If the variable is not present in the environment, the default value is returned.
func GetInt(key string, defaultValue int) int {
	return int(GetInt64(key, int64(defaultValue)))
}

// GetInt64 returns the value of the environment variable named by the key.
// If the variable is not present in the environment, the default value is returned.
func GetInt64(key string, defaultValue int64) int64 {
	value := Get(key)
	if value == "" {
		return defaultValue
	}

	return ParseInt(value, defaultValue)
}

// GetUint returns the value of the environment variable named by the key.
// If the variable is not present in the environment, the default value is returned.
func GetUint(key string, defaultValue uint) uint {
	return uint(GetUint64(key, uint64(defaultValue)))
}

// GetUint64 returns the value of the environment variable named by the key.
// If the variable is not present in the environment, the default value is returned.
func GetUint64(key string, defaultValue uint64) uint64 {
	value := Get(key)
	if value == "" {
		return defaultValue
	}

	return ParseUint(value, defaultValue)
}

// GetBool returns the value of the environment variable named by the key.
// If the variable is not present in the environment, the default value is returned.
func GetBool(key string, defaultValue bool) bool {
	value := Get(key)
	if value == "" {
		return defaultValue
	}

	return ParseBool(value, defaultValue)
}

// GetFloat64 returns the value of the environment variable named by the key.
// If the variable is not present in the environment, the default value is returned.
func GetFloat64(key string, defaultValue float64) float64 {
	value := Get(key)
	if value == "" {
		return defaultValue
	}

	return ParseFloat(value, defaultValue)
}

// Get returns the value of the environment variable named by the key.
// If the variable is not present in the environment, an empty string is returned.
func Get(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return ""
	}

	return value
}

// ParseInt parses the string value into an int64.
// If the string is empty, the default value is returned.
func ParseInt(value string, defaultValue int64) int64 {
	if value == "" {
		return defaultValue
	}

	i, err := parseInt(value)
	if err != nil {
		return defaultValue
	}

	return i
}

// ParseUint parses the string value into an uint64.
// If the string is empty, the default value is returned.
func ParseUint(value string, defaultValue uint64) uint64 {
	if value == "" {
		return defaultValue
	}

	i, err := parseUint(value)
	if err != nil {
		return defaultValue
	}

	return i
}

// ParseBool parses the string value into a bool.
// If the string is empty, the default value is returned.
func ParseBool(value string, defaultValue bool) bool {
	if value == "" {
		return defaultValue
	}

	b, err := parseBool(value)
	if err != nil {
		return defaultValue
	}

	return b
}

// ParseFloat parses the string value into a float64.
// If the string is empty, the default value is returned.
func ParseFloat(value string, defaultValue float64) float64 {
	if value == "" {
		return defaultValue
	}

	f, err := parseFloat(value)
	if err != nil {
		return defaultValue
	}

	return f
}

func parseInt(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

func parseUint(value string) (uint64, error) {
	return strconv.ParseUint(value, 10, 64)
}

func parseBool(value string) (bool, error) {
	return strconv.ParseBool(value)
}

func parseFloat(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

// MustGetString returns the value of the environment variable named by the key.
// If the variable is not present in the environment, the default value is returned.
// If the variable is not present in the environment and the default value is empty, it panics.
// func MustGetString(key string, defaultValue string) string {
// 	value := GetString(key, defaultValue)
// 	if value == "" {
// 		panic(fmt.Sprintf("environment variable %s is not set", key))
// 	}

// 	return value
// }
