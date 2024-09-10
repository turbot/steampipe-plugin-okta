package okta

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/memoize"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func commonColumns(c []*plugin.Column) []*plugin.Column {
	return append([]*plugin.Column{
		{
			Name:        "domain",
			Description: "The okta domain name.",
			Hydrate:     getOktaDomainName,
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromValue(),
		},
	}, c...)
}

// if the caching is required other than per connection, build a cache key for the call and use it in Memoize.
var getOktaDomainNameMemoized = plugin.HydrateFunc(getOktaDomainNameUncached).Memoize(memoize.WithCacheKeyFunction(getOktaDomainNameCacheKey))

// declare a wrapper hydrate function to call the memoized function
// - this is required when a memoized function is used for a column definition
func getOktaDomainName(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return getOktaDomainNameMemoized(ctx, d, h)
}

// Build a cache key for the call to getOktaDomainNameCacheKey.
func getOktaDomainNameCacheKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	key := "getOktaDomainName"
	return key, nil
}

func getOktaDomainNameUncached(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Retrieve the Okta configuration
	config := GetConfig(d.Connection)

	// Retrieve the domain from the config
	domain := config.Domain

	// If the domain is not set in the config, fall back to the environment variable
	if domain == nil {
		envDomain := os.Getenv("OKTA_CLIENT_ORGURL")
		domain = &envDomain
	}

	// Extract the domain name by removing the "https://" prefix
	splitDomain := strings.SplitN(*domain, "https://", 2)
	if len(splitDomain) != 2 {
		return nil, fmt.Errorf("invalid okta domain format: %s", *domain)
	}
	domainName := splitDomain[1]

	return domainName, nil
}
