package okta

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
