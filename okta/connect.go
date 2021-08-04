package okta

import (
	"context"
	"os"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func Connect(ctx context.Context, d *plugin.QueryData) (*okta.Client, error) {
	// have we already created and cached the session?
	sessionCacheKey := "OktaSession"
	if cachedData, ok := d.ConnectionManager.Cache.Get(sessionCacheKey); ok {
		return cachedData.(*okta.Client), nil
	}

	oktaConfig := GetConfig(d.Connection)

	var domain, token string
	if oktaConfig.Domain != nil {
		domain = *oktaConfig.Domain
	} else {
		domain = os.Getenv("OKTA_DOMAIN")
	}

	if oktaConfig.Token != nil {
		token = *oktaConfig.Token
	} else {
		token = os.Getenv("OKTA_TOKEN")
	}

	if domain == "" {
		panic("please set domain in steampipe config. Edit your connection configuration file and then restart Steampipe")
	}

	if token == "" {
		panic("please set domain in steampipe config. Edit your connection configuration file and then restart Steampipe")
	}

	_, client, err := okta.NewClient(ctx, okta.WithOrgUrl(domain), okta.WithToken(token), okta.WithRequestTimeout(15), okta.WithRateLimitMaxRetries(5))

	// Save session into cache
	d.ConnectionManager.Cache.Set(sessionCacheKey, client)

	return client, err
}
