package okta

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

// Handle rate limit errors  
func shouldRetryError(retryErrors []string) plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {

		if strings.Contains(err.Error(), "429") {
			plugin.Logger(ctx).Debug("okta_errors.shouldRetryError", "rate_limit_error", err)
			return true
		}
		return false
	}
}