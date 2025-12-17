package configuration

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
)

type Configuration struct {
	Port              string   `env:"PORT" default:"8080"`
	Address           string   `env:"ADDRESS" default:"localhost"`
	ClientDirectory   string   `env:"CLIENT_DIRECTORY" default:"../client/dist"`
	AllowedOrigins    []string `env:"ALLOWED_ORIGINS" separator:","`
	JWTSecret         string   `env:"JWT_SECRET" default:"SECRET"`
	InterserverSecret string   `env:"INTERSERVER_SECRET" default:"SECRET"`
}

func Load[T any]() (T, error) {
	var configuration T

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment.")
	}

	configurationValue := reflect.ValueOf(&configuration).Elem()
	configurationType := configurationValue.Type()

	for fieldIndex := 0; fieldIndex < configurationType.NumField(); fieldIndex++ {
		structField := configurationType.Field(fieldIndex)
		structFieldValue := configurationValue.Field(fieldIndex)

		environmentVariableName := structField.Tag.Get("env")
		if environmentVariableName == "" {
			continue
		}

		defaultValue := structField.Tag.Get("default")

		environmentValue := os.Getenv(environmentVariableName)
		if environmentValue == "" {
			environmentValue = defaultValue
		}

		if err := setFieldValueFromString(structFieldValue, environmentValue, structField.Tag); err != nil {
			return configuration, fmt.Errorf(
				"configuration field %s: %w",
				structField.Name,
				err,
			)
		}

		if environmentValue != "" {
			fmt.Printf("%s -> %s\n", environmentVariableName, environmentValue)
		}
	}

	return configuration, nil
}

func setFieldValueFromString(
	fieldValue reflect.Value,
	rawValue string,
	structTag reflect.StructTag,
) error {
	switch fieldValue.Kind() {

	case reflect.String:
		fieldValue.SetString(rawValue)

	case reflect.Slice:
		separator := structTag.Get("separator")
		if separator == "" {
			separator = ","
		}

		rawParts := strings.Split(rawValue, separator)

		sliceValue := reflect.MakeSlice(
			fieldValue.Type(),
			0,
			len(rawParts),
		)

		for _, rawPart := range rawParts {
			trimmedPart := strings.TrimSpace(rawPart)
			sliceValue = reflect.Append(
				sliceValue,
				reflect.ValueOf(trimmedPart),
			)
		}

		fieldValue.Set(sliceValue)

	default:
		return fmt.Errorf(
			"unsupported configuration field type %s",
			fieldValue.Kind(),
		)
	}

	return nil
}
