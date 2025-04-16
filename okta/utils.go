package okta

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/ettle/strcase"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

const (
	filterTimeFormat = "2006-01-02T15:04:05.000Z"
	titleDescription = "The title of the resource."
)

var (
	// Filters sympol - comparison operator map for okta
	operatorsMap = map[string]string{
		"=":  "eq",
		">=": "ge",
		">":  "gt",
		"<=": "le",
		"<":  "lt",
		"<>": "ne",
	}
)

func getListValues(listValue *proto.QualValueList) []*string {
	values := make([]*string, 0)
	for _, value := range listValue.Values {
		values = append(values, types.String(value.GetStringValue()))
	}
	return values
}

//// other useful functions

func buildQueryFilter(equalQuals plugin.KeyColumnEqualsQualMap, filterKeys []string) []string {
	filters := []string{}

	for k, v := range equalQuals {
		if v != nil && slices.Contains(filterKeys, k) {
			filters = append(filters, fmt.Sprintf("%s eq \"%s\"", strcase.ToCamel(k), v.GetStringValue()))
		}
	}

	return filters
}

// StructToMap converts the fields of a struct from interface{} to map[string]interface{}
func structToMap(input interface{}) (map[string]interface{}, error) {
	// Create the result map
	if input == nil {
		return nil, nil
	}
	result := make(map[string]interface{})

	// Get the value and type of the input
	value := reflect.ValueOf(input)
	typ := reflect.TypeOf(input)

	// Ensure the input is a struct or a pointer to a struct
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
		typ = typ.Elem()
	}

	if value.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct or a pointer to a struct")
	}

	// Iterate over struct fields
	for i := 0; i < value.NumField(); i++ {
		// Get field name and value
		field := value.Field(i)
		fieldType := typ.Field(i)
		fieldName := fieldType.Name

		// Only exportable fields (starting with an uppercase letter) can be accessed
		if field.CanInterface() {
			result[fieldName] = field.Interface()
		}
	}

	return result, nil
}
