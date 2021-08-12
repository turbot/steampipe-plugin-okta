package okta

import (
	"fmt"

	"github.com/ettle/strcase"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
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
	// TODO - policy types creating issue "OAUTH_AUTHORIZATION_POLICY", "ACCESS_POLICY", "PROFILE_ENROLLMENT"
	policyTypes = []string{"OKTA_SIGN_ON", "PASSWORD", "MFA_ENROLL", "IDP_DISCOVERY"}
)

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
