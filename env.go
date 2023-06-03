package env

import (
	"errors"
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
//	if err := Parse(&config); err != nil {
//		// handle error
//	}
//
//	fmt.Println(config.Port)
//
// If the environment variable is not present, an error is returned.
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
			if err := parse(value.Addr().Interface(), field.Tag.Get(DefaultTag)); err != nil {
				return err
			}
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
// If the environment variable is not present, an error is returned.
// If the environment variable is present, but the field cannot be set, an error
// is returned.
func setField(value reflect.Value, env string) error {
	if !value.CanSet() {
		return errors.New("cannot set field value")
	}

	if !value.IsValid() {
		return errors.New("invalid field value")
	}

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.String:
		if s, ok := GetString(env); ok {
			value.SetString(s)
		} else {
			return errors.New("environment variable not found: " + env)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if i, ok := GetInt64(env); ok {
			value.SetInt(i)
		} else {
			return errors.New("environment variable not found: " + env)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if u, ok := GetUint64(env); ok {
			value.SetUint(u)
		} else {
			return errors.New("environment variable not found: " + env)
		}
	case reflect.Bool:
		if b, ok := GetBool(env); ok {
			value.SetBool(b)
		} else {
			return errors.New("environment variable not found: " + env)
		}
	case reflect.Float32, reflect.Float64:
		if f, ok := GetFloat64(env); ok {
			value.SetFloat(f)
		} else {
			return errors.New("environment variable not found: " + env)
		}
	}

	return nil
}

// GetString returns the value of the environment variable named by the key.
// If the variable is not present in the environment, an empty string and false are returned.
func GetString(key string) (string, bool) {
	value, ok := os.LookupEnv(key)
	return value, ok
}

// GetInt returns the value of the environment variable named by the key.
// If the variable is not present in the environment, 0 and false are returned.
func GetInt(key string) (int, bool) {
	i, ok := GetInt64(key)
	return int(i), ok
}

// GetInt64 returns the value of the environment variable named by the key.
// If the variable is not present in the environment, 0 and false are returned.
func GetInt64(key string) (int64, bool) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return 0, false
	}

	i, err := parseInt(value)
	if err != nil {
		return 0, false
	}

	return i, true
}

// GetUint returns the value of the environment variable named by the key.
// If the variable is not present in the environment, 0 and false are returned.
func GetUint(key string) (uint, bool) {
	u, ok := GetUint64(key)
	return uint(u), ok
}

// GetUint64 returns the value of the environment variable named by the key.
// If the variable is not present in the environment, 0 and false are returned.
func GetUint64(key string) (uint64, bool) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return 0, false
	}

	u, err := parseUint(value)
	if err != nil {
		return 0, false
	}

	return u, true
}

// GetBool returns the value of the environment variable named by the key.
// If the variable is not present in the environment, false and false are returned.
func GetBool(key string) (bool, bool) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return false, false
	}

	b, err := parseBool(value)
	if err != nil {
		return false, false
	}

	return b, true
}

// GetFloat64 returns the value of the environment variable named by the key.
// If the variable is not present in the environment, 0 and false are returned.
func GetFloat64(key string) (float64, bool) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return 0, false
	}

	f, err := parseFloat(value)
	if err != nil {
		return 0, false
	}

	return f, true
}

// ParseInt parses the string value into an int64.
// If the string is empty or parsing fails, 0 and false are returned.
func ParseInt(value string) (int64, bool) {
	if value == "" {
		return 0, false
	}

	i, err := parseInt(value)
	if err != nil {
		return 0, false
	}

	return i, true
}

// ParseUint parses the string value into an uint64.
// If the string is empty or parsing fails, 0 and false are returned.
func ParseUint(value string) (uint64, bool) {
	if value == "" {
		return 0, false
	}

	u, err := parseUint(value)
	if err != nil {
		return 0, false
	}

	return u, true
}

// ParseBool parses the string value into a bool.
// If the string is empty or parsing fails, false and false are returned.
func ParseBool(value string) (bool, bool) {
	if value == "" {
		return false, false
	}

	b, err := parseBool(value)
	if err != nil {
		return false, false
	}

	return b, true
}

// ParseFloat parses the string value into a float64.
// If the string is empty or parsing fails, 0 and false are returned.
func ParseFloat(value string) (float64, bool) {
	if value == "" {
		return 0, false
	}

	f, err := parseFloat(value)
	if err != nil {
		return 0, false
	}

	return f, true
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
