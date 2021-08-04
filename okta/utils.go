package okta

import (
	"fmt"

	"github.com/ettle/strcase"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

const filterTimeFormat = "2006-01-02T15:04:05.000Z"

// Filters sympol - comparison operator map for okta
var operatorsMap = map[string]string{
	"=":  "eq",
	">=": "ge",
	">":  "gt",
	"<=": "le",
	"<":  "lt",
	"<>": "ne",
}

//// other useful functions

func buildQueryFilter(equalQuals plugin.KeyColumnEqualsQualMap) []string {
	filters := []string{}

	for k, v := range equalQuals {
		if v != nil {
			filters = append(filters, fmt.Sprintf("%s eq \"%s\"", strcase.ToCamel(k), v.GetStringValue()))
		}
	}

	return filters
}
