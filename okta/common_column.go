package okta

import (
	"context"
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
	cfg := GetConfig(d.Connection)

	domain := cfg.Domain

	if domain == nil {
		d := os.Getenv("OKTA_CLIENT_ORGURL")
		domain = &d
	}

	domainName := strings.Split(*domain, "https://")[1]

	return domainName, nil
}
