package mockserver

import (
	cryptorand "crypto/rand"
	"encoding/hex"
	"log/slog"
	"maps"
	"math/big"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/mmadfox/swag2mcp/internal/spec"
)

const (
	randomHexLength16 = 16
	randomHexLength8  = 8
	passwordLength    = 16
	sentenceWordCount = 3
	arrayMinLength    = 1
	arrayMaxLength    = 3
	nullableChance    = 10
	maxInt32Value     = 100000
	maxInt64Value     = 100000
	maxIntValue       = 1000
	maxNumberValue    = 10000
	booleanDivisor    = 2
	numberDivisor     = 100
	hexByteDivisor    = 2

	schemaTypeString  = "string"
	schemaTypeInteger = "integer"
	schemaTypeObject  = "object"
	formatURI         = "uri"
	formatURL         = "url"
	formatPassword    = "password"
	formatInt32       = "int32"
	formatInt64       = "int64"
)

// GenerateFromSchema generates random data that conforms to the given OpenAPI schema.
// It handles all common schema types including string, integer, number, boolean,
// array, object, enum, allOf, oneOf, anyOf, and format-specific values.
// For nullable schemas, it returns nil approximately 10% of the time.
func GenerateFromSchema(schema *spec.Schema) any {
	if schema == nil {
		return nil
	}

	if schema.Nullable && cryptoRandIntN(nullableChance) == 0 {
		return nil
	}

	if len(schema.Enum) > 0 {
		index := cryptoRandIntN(len(schema.Enum))
		return schema.Enum[index]
	}

	if schema.Default != nil {
		return schema.Default
	}

	if schema.Example != nil {
		return schema.Example
	}

	if len(schema.AllOf) > 0 {
		return generateAllOf(schema.AllOf)
	}

	if len(schema.OneOf) > 0 {
		index := cryptoRandIntN(len(schema.OneOf))
		return GenerateFromSchema(schema.OneOf[index])
	}

	if len(schema.AnyOf) > 0 {
		index := cryptoRandIntN(len(schema.AnyOf))
		return GenerateFromSchema(schema.AnyOf[index])
	}

	switch schema.Type {
	case schemaTypeString:
		return generateStringValue(schema.Format)
	case schemaTypeInteger:
		return generateIntegerValue(schema.Format)
	case "number":
		return generateNumberValue()
	case "boolean":
		return cryptoRandIntN(booleanDivisor) == 0
	case "array":
		return generateArrayValue(schema.Items)
	case schemaTypeObject:
		return generateObjectValue(schema)
	default:
		if len(schema.Properties) > 0 {
			return generateObjectValue(schema)
		}
		return generateStringValue("")
	}
}

// generateStringValue returns a random string based on the given format.
// Supported formats: date, date-time, time, email, uri, url, uuid, ipv4,
// ipv6, hostname, mac, phone, byte, binary, password, int32, int64.
// For unknown formats, it picks a random word, sentence, name, or other value.
func generateStringValue(format string) string {
	switch format {
	case "date":
		return gofakeit.Date().Format("2006-01-02")
	case "date-time":
		return gofakeit.Date().Format("2006-01-02T15:04:05Z")
	case "time":
		return gofakeit.Date().Format("15:04:05")
	case "email":
		return gofakeit.Email()
	case formatURI, formatURL:
		return gofakeit.URL()
	case "uuid":
		return gofakeit.UUID()
	case "ipv4":
		return gofakeit.IPv4Address()
	case "ipv6":
		return gofakeit.IPv6Address()
	case "hostname":
		return gofakeit.DomainName()
	case "mac":
		return gofakeit.MacAddress()
	case "phone":
		return gofakeit.Phone()
	case "byte":
		return randomHexString(randomHexLength16)
	case "binary":
		return randomHexString(randomHexLength8)
	case formatPassword:
		return gofakeit.Password(true, true, true, true, false, passwordLength)
	case formatInt32, formatInt64:
		return gofakeit.Letter()
	default:
		generators := []func() string{
			gofakeit.Word,
			func() string { return gofakeit.Sentence(sentenceWordCount) },
			gofakeit.Name,
			gofakeit.City,
			gofakeit.Country,
			gofakeit.Color,
			gofakeit.Company,
			gofakeit.ProductName,
			gofakeit.CarModel,
		}
		index := cryptoRandIntN(len(generators))
		return generators[index]()
	}
}

// generateIntegerValue returns a random integer based on the given format.
// For int32 it returns int32, for int64 it returns int64, otherwise int.
func generateIntegerValue(format string) any {
	switch format {
	case "int32":
		return int32(cryptoRandIntN(maxInt32Value)) //nolint:gosec // maxInt32Value fits in int32
	case "int64":
		return cryptoRandInt64N(maxInt64Value)
	default:
		return cryptoRandIntN(maxIntValue)
	}
}

// generateNumberValue returns a random float64 value.
func generateNumberValue() any {
	return float64(cryptoRandIntN(maxNumberValue)) / numberDivisor
}

// generateArrayValue returns a slice of random values based on the items schema.
// The array length is between arrayMinLength and arrayMaxLength.
func generateArrayValue(itemsSchema *spec.Schema) any {
	if itemsSchema == nil {
		return []any{}
	}

	length := cryptoRandIntN(arrayMaxLength) + arrayMinLength
	result := make([]any, length)
	for index := range result {
		result[index] = GenerateFromSchema(itemsSchema)
	}
	return result
}

// generateObjectValue returns a map of random values based on the schema properties.
// ReadOnly and WriteOnly properties are skipped.
func generateObjectValue(schema *spec.Schema) map[string]any {
	result := make(map[string]any, len(schema.Properties))

	for propertyName, propertySchema := range schema.Properties {
		if propertySchema == nil {
			continue
		}
		if propertySchema.ReadOnly {
			continue
		}
		if propertySchema.WriteOnly {
			continue
		}
		result[propertyName] = GenerateFromSchema(propertySchema)
	}

	return result
}

// generateAllOf merges multiple schemas into a single object by generating
// data for each sub-schema and combining the results.
func generateAllOf(schemas []*spec.Schema) map[string]any {
	result := make(map[string]any)

	for _, subSchema := range schemas {
		if subSchema == nil {
			continue
		}
		generated := GenerateFromSchema(subSchema)
		if generatedMap, isMap := generated.(map[string]any); isMap {
			maps.Copy(result, generatedMap)
		}
	}

	return result
}

// randomHexString returns a random hex string of the given length.
func randomHexString(length int) string {
	byteLength := (length + hexByteDivisor - 1) / hexByteDivisor
	randomBytes := make([]byte, byteLength)
	if _, err := cryptorand.Read(randomBytes); err != nil {
		slog.Default().Warn("failed to generate random hex", "error", err)
	}
	return hex.EncodeToString(randomBytes)[:length]
}

// cryptoRandIntN returns a cryptographically secure random integer in [0, maximum).
func cryptoRandIntN(maximum int) int {
	if maximum <= 0 {
		return 0
	}
	bigMax := big.NewInt(int64(maximum))
	result, err := cryptorand.Int(cryptorand.Reader, bigMax)
	if err != nil {
		return 0
	}
	return int(result.Int64())
}

// cryptoRandInt64N returns a cryptographically secure random int64 in [0, maximum).
func cryptoRandInt64N(maximum int) int64 {
	if maximum <= 0 {
		return 0
	}
	bigMax := big.NewInt(int64(maximum))
	result, err := cryptorand.Int(cryptorand.Reader, bigMax)
	if err != nil {
		return 0
	}
	return result.Int64()
}
