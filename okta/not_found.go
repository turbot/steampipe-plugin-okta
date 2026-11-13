package okta

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

// function which returns an ErrorPredicate for Okta API calls
func isNotFoundError(notFoundErrors []string) plugin.ErrorPredicateWithContext {
	return func(_ context.Context, _ *plugin.QueryData, _ *plugin.HydrateData, err error) bool {
		if err != nil {
			for _, item := range notFoundErrors {
				if strings.Contains(err.Error(), item) {
					return true
				}
			}
		}
		return false
	}
}
