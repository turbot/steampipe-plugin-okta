package okta

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Handle rate limit errors
func shouldRetryError(retryErrors []string) plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {

		if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "context deadline exceeded") {
			plugin.Logger(ctx).Debug("okta_errors.shouldRetryError", "rate_limit_error", err)
			return true
		}
		return false
	}
}
