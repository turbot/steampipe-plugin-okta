package okta

import (
	"context"
	"os"

	"github.com/okta/okta-sdk-golang/v2/okta"
	oktaV3 "github.com/okta/okta-sdk-golang/v3/okta"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func Connect(ctx context.Context, d *plugin.QueryData) (*okta.Client, error) {
	// have we already created and cached the session?
	sessionCacheKey := "OktaSession"
	if cachedData, ok := d.ConnectionManager.Cache.Get(sessionCacheKey); ok {
		return cachedData.(*okta.Client), nil
	}

	oktaConfig := GetConfig(d.Connection)

	var domain, token, clientID, privateKey string
	scopes := []string{"okta.users.read", "okta.groups.read", "okta.roles.read", "okta.apps.read", "okta.policies.read", "okta.authorizationServers.read", "okta.trustedOrigins.read", "okta.factors.read"}
	if oktaConfig.Domain != nil {
		domain = *oktaConfig.Domain
	} else {
		domain = os.Getenv("OKTA_CLIENT_ORGURL")
	}

	if oktaConfig.Token != nil {
		token = *oktaConfig.Token
	} else {
		token = os.Getenv("OKTA_CLIENT_TOKEN")
	}

	if domain != "" && token != "" {
		_, client, err := okta.NewClient(ctx, okta.WithOrgUrl(domain), okta.WithToken(token), okta.WithRequestTimeout(30), okta.WithRateLimitMaxRetries(5))
		if err != nil {
			return nil, err
		}
		d.ConnectionManager.Cache.Set(sessionCacheKey, client)
		return client, err
	}

	if oktaConfig.ClientID != nil {
		clientID = *oktaConfig.ClientID
	} else {
		clientID = os.Getenv("OKTA_CLIENT_CLIENTID")
	}

	if oktaConfig.PrivateKey != nil {
		privateKey = *oktaConfig.PrivateKey
	} else {
		privateKey = os.Getenv("OKTA_CLIENT_PRIVATEKEY")
	}

	if domain != "" && clientID != "" && privateKey != "" {
		_, client, err := okta.NewClient(ctx, okta.WithOrgUrl(domain), okta.WithAuthorizationMode("PrivateKey"), okta.WithClientId(clientID), okta.WithPrivateKey(privateKey), okta.WithScopes(scopes), okta.WithRequestTimeout(15), okta.WithRateLimitMaxRetries(5))
		if err != nil {
			return nil, err
		}
		d.ConnectionManager.Cache.Set(sessionCacheKey, client)
		return client, err
	}

	/* *
	*	Try with okta sdk default options
	*
	*	https://github.com/okta/okta-sdk-golang#configuration-reference
	*	1. An okta.yaml file in a .okta folder in the current user's home directory (~/.okta/okta.yaml or %userprofile\.okta\okta.yaml)
	* 2. An .okta.yaml file in the application or project's root directory
	* 3. Environment variables
	* 4. Configuration explicitly passed to the constructor (see the example in Getting started)
	*	*/
	_, client, err := okta.NewClient(ctx, okta.WithRequestTimeout(30), okta.WithRateLimitMaxRetries(5))
	if err != nil {
		return nil, err
	}

	// Save session into cache
	d.ConnectionManager.Cache.Set(sessionCacheKey, client)

	return client, err
}


func ConnectV3(ctx context.Context, d *plugin.QueryData) (*oktaV3.APIClient, error) {
	// have we already created and cached the session?
	sessionCacheKey := "OktaSessionV3"
	if cachedData, ok := d.ConnectionManager.Cache.Get(sessionCacheKey); ok {
		return cachedData.(*oktaV3.APIClient), nil
	}

	oktaConfig := GetConfig(d.Connection)

	var domain, token, clientID, privateKey string
	scopes := []string{"okta.users.read", "okta.groups.read", "okta.roles.read", "okta.apps.read", "okta.policies.read", "okta.authorizationServers.read", "okta.trustedOrigins.read", "okta.factors.read", "okta.devices.read"}
	if oktaConfig.Domain != nil {
		domain = *oktaConfig.Domain
	} else {
		domain = os.Getenv("OKTA_CLIENT_ORGURL")
	}

	if oktaConfig.Token != nil {
		token = *oktaConfig.Token
	} else {
		token = os.Getenv("OKTA_CLIENT_TOKEN")
	}

	if domain != "" && token != "" {
		oktaConfiguratiopn, err := oktaV3.NewConfiguration(oktaV3.WithOrgUrl(domain), oktaV3.WithToken(token), oktaV3.WithRequestTimeout(30), oktaV3.WithRateLimitMaxRetries(5))
		if err != nil {
			return nil, err
		}
		client := oktaV3.NewAPIClient(oktaConfiguratiopn)

		d.ConnectionManager.Cache.Set(sessionCacheKey, client)
		return client, err
	}

	if oktaConfig.ClientID != nil {
		clientID = *oktaConfig.ClientID
	} else {
		clientID = os.Getenv("OKTA_CLIENT_CLIENTID")
	}

	if oktaConfig.PrivateKey != nil {
		privateKey = *oktaConfig.PrivateKey
	} else {
		privateKey = os.Getenv("OKTA_CLIENT_PRIVATEKEY")
	}

	if domain != "" && clientID != "" && privateKey != "" {
		oktaConfiguratiopn, err := oktaV3.NewConfiguration(oktaV3.WithOrgUrl(domain), oktaV3.WithAuthorizationMode("PrivateKey"), oktaV3.WithClientId(clientID), oktaV3.WithPrivateKey(privateKey), oktaV3.WithScopes(scopes), oktaV3.WithRequestTimeout(15), oktaV3.WithRateLimitMaxRetries(5))
		if err != nil {
			return nil, err
		}

		client := oktaV3.NewAPIClient(oktaConfiguratiopn)

		d.ConnectionManager.Cache.Set(sessionCacheKey, client)
		return client, err
	}

	/* *
	*	Try with okta sdk default options
	*
	*	https://github.com/okta/okta-sdk-golang#configuration-reference
	*	1. An okta.yaml file in a .okta folder in the current user's home directory (~/.okta/okta.yaml or %userprofile\.okta\okta.yaml)
	* 2. An .okta.yaml file in the application or project's root directory
	* 3. Environment variables
	* 4. Configuration explicitly passed to the constructor (see the example in Getting started)
	*	*/
	oktaConfiguratiopn, err := oktaV3.NewConfiguration(oktaV3.WithRequestTimeout(30), oktaV3.WithRateLimitMaxRetries(5))
	if err != nil {
			return nil, err
		}
	client := oktaV3.NewAPIClient(oktaConfiguratiopn)

	// Save session into cache
	d.ConnectionManager.Cache.Set(sessionCacheKey, client)

	return client, err
}
