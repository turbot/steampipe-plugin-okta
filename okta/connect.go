package okta

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/okta/okta-sdk-golang/v2/okta"
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

	// The default value has been set as per the API doc: https://github.com/okta/okta-sdk-golang?tab=readme-ov-file#environment-variables
	// SDK supported environment variables: https://github.com/okta/okta-sdk-golang/blob/master/okta/config.go#L33-L70
	var requestTimeout, maxBackoff int64 = 30, 30
	var maxRetries int32 = 5

	if oktaConfig.MaxBackoff != nil {
		maxBackoff = *oktaConfig.MaxBackoff
	} else {
		maxBackoffIntValue, err := strconv.ParseInt(os.Getenv("OKTA_CLIENT_RATE_LIMIT_MAX_BACKOFF"), 10, 64)
		if err != nil {
			// handle the error in case of invalid string
			return nil, fmt.Errorf("Error converting max backoff string type to int64:", err)
		}
		maxBackoff = maxBackoffIntValue
	}

	if oktaConfig.RequestTimeout != nil {
		requestTimeout = *oktaConfig.RequestTimeout
	} else {
		requestTimeoutIntValue, err := strconv.ParseInt(os.Getenv("OKTA_CLIENT_REQUEST_TIMEOUT"), 10, 64)
		if err != nil {
			// handle the error in case of invalid string
			return nil, fmt.Errorf("Error converting request timeout string type to int64:", err)
		}
		requestTimeout = requestTimeoutIntValue
	}

	if oktaConfig.MaxRetries != nil {
		maxRetries = *oktaConfig.MaxRetries
	} else {
		maxRetriesIntValue, err := strconv.ParseInt(os.Getenv("OKTA_CLIENT_RATE_LIMIT_MAX_RETRIES"), 10, 32)
		if err != nil {
			// handle the error in case of invalid string
			return nil, fmt.Errorf("Error converting max retries string type to int32:", err)
		}
		maxRetries = int32(maxRetriesIntValue)
	}

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
		_, client, err := okta.NewClient(ctx, okta.WithOrgUrl(domain), okta.WithToken(token), okta.WithRequestTimeout(requestTimeout), okta.WithRateLimitMaxRetries(maxRetries), okta.WithRateLimitMaxBackOff(maxBackoff))
		if err != nil {
			return nil, err
		}
		d.ConnectionManager.Cache.Set(sessionCacheKey, client)
		client.GetConfig()
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
		_, client, err := okta.NewClient(ctx, okta.WithOrgUrl(domain), okta.WithAuthorizationMode("PrivateKey"), okta.WithClientId(clientID), okta.WithPrivateKey(privateKey), okta.WithScopes(scopes), okta.WithRequestTimeout(requestTimeout), okta.WithRateLimitMaxRetries(maxRetries), okta.WithRateLimitMaxBackOff(maxBackoff))
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
	_, client, err := okta.NewClient(ctx, okta.WithRequestTimeout(requestTimeout), okta.WithRateLimitMaxRetries(maxRetries), okta.WithRateLimitMaxBackOff(maxBackoff))
	if err != nil {
		return nil, err
	}

	// Save session into cache
	d.ConnectionManager.Cache.Set(sessionCacheKey, client)

	return client, err
}
