package okta

import (
	"fmt"

	"github.com/ettle/strcase"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
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
		if v != nil && helpers.StringSliceContains(filterKeys, k) {
			filters = append(filters, fmt.Sprintf("%s eq \"%s\"", strcase.ToCamel(k), v.GetStringValue()))
		}
	}

	return filters
}
