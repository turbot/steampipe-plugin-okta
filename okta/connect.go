package okta

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/okta/okta-sdk-golang/v2/okta"
	oktaV4 "github.com/okta/okta-sdk-golang/v4/okta"
	oktaV5 "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func Connect(ctx context.Context, d *plugin.QueryData) (*okta.Client, error) {
	// have we already created and cached the session?
	sessionCacheKey := "OktaSession"
	if cachedData, ok := d.ConnectionManager.Cache.Get(sessionCacheKey); ok {
		return cachedData.(*okta.Client), nil
	}

	// Get environment or steampipe config value
	domain, token, clientID, privateKey, requestTimeout, maxBackoff, maxRetries, err := getOktaConfigValues(d)
	if err != nil {
		return nil, fmt.Errorf("error in retrieving config or environment values: %v", err)
	}

	scopes := []string{"okta.users.read", "okta.groups.read", "okta.roles.read", "okta.apps.read", "okta.policies.read", "okta.authorizationServers.read", "okta.trustedOrigins.read", "okta.factors.read"}

	if domain != "" && token != "" {
		_, client, err := okta.NewClient(ctx, okta.WithOrgUrl(domain), okta.WithToken(token), okta.WithRequestTimeout(requestTimeout), okta.WithRateLimitMaxRetries(maxRetries), okta.WithRateLimitMaxBackOff(maxBackoff))
		if err != nil {
			return nil, err
		}
		d.ConnectionManager.Cache.Set(sessionCacheKey, client)
		client.GetConfig()
		return client, err
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

func ConnectV4(ctx context.Context, d *plugin.QueryData) (*oktaV4.APIClient, error) {
	// have we already created and cached the session?
	sessionCacheKey := "OktaSessionV4"
	if cachedData, ok := d.ConnectionManager.Cache.Get(sessionCacheKey); ok {
		return cachedData.(*oktaV4.APIClient), nil
	}

	// Get environment or steampipe config value
	domain, token, clientID, privateKey, requestTimeout, maxBackoff, maxRetries, err := getOktaConfigValues(d)
	if err != nil {
		return nil, fmt.Errorf("error in retrieving config or environment values: %v", err)
	}
	scopes := []string{"okta.users.read", "okta.groups.read", "okta.roles.read", "okta.apps.read", "okta.policies.read", "okta.authorizationServers.read", "okta.trustedOrigins.read", "okta.factors.read", "okta.devices.read"}

	if domain != "" && token != "" {
		oktaConfiguratiopn, err := oktaV4.NewConfiguration(oktaV4.WithOrgUrl(domain), oktaV4.WithToken(token), oktaV4.WithRequestTimeout(requestTimeout), oktaV4.WithRateLimitMaxRetries(maxRetries), oktaV4.WithRateLimitMaxBackOff(maxBackoff))
		if err != nil {
			return nil, err
		}
		client := oktaV4.NewAPIClient(oktaConfiguratiopn)

		d.ConnectionManager.Cache.Set(sessionCacheKey, client)
		return client, err
	}

	if domain != "" && clientID != "" && privateKey != "" {
		oktaConfiguratiopn, err := oktaV4.NewConfiguration(oktaV4.WithOrgUrl(domain), oktaV4.WithAuthorizationMode("PrivateKey"), oktaV4.WithClientId(clientID), oktaV4.WithPrivateKey(privateKey), oktaV4.WithScopes(scopes), oktaV4.WithRequestTimeout(requestTimeout), oktaV4.WithRateLimitMaxRetries(maxRetries), oktaV4.WithRateLimitMaxBackOff(maxBackoff))
		if err != nil {
			return nil, err
		}

		client := oktaV4.NewAPIClient(oktaConfiguratiopn)

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
	oktaConfiguratiopn, err := oktaV4.NewConfiguration(oktaV4.WithRequestTimeout(requestTimeout), oktaV4.WithRateLimitMaxRetries(maxRetries), oktaV4.WithRateLimitMaxBackOff(maxBackoff))
	if err != nil {
		return nil, err
	}
	client := oktaV4.NewAPIClient(oktaConfiguratiopn)

	// Save session into cache
	d.ConnectionManager.Cache.Set(sessionCacheKey, client)

	return client, err
}

func ConnectV5(ctx context.Context, d *plugin.QueryData) (*oktaV5.APIClient, error) {
	// have we already created and cached the session?
	sessionCacheKey := "OktaSessionV5"
	if cachedData, ok := d.ConnectionManager.Cache.Get(sessionCacheKey); ok {
		return cachedData.(*oktaV5.APIClient), nil
	}

	// Get environment or steampipe config value
	domain, token, clientID, privateKey, requestTimeout, maxBackoff, maxRetries, err := getOktaConfigValues(d)
	if err != nil {
		return nil, fmt.Errorf("error in retrieving config or environment values: %v", err)
	}

	scopes := []string{"okta.users.read", "okta.groups.read", "okta.roles.read", "okta.apps.read", "okta.policies.read", "okta.authorizationServers.read", "okta.trustedOrigins.read", "okta.factors.read", "okta.devices.read"}

	if domain != "" && token != "" {
		oktaConfiguratiopn, err := oktaV5.NewConfiguration(oktaV5.WithOrgUrl(domain), oktaV5.WithToken(token), oktaV5.WithRequestTimeout(requestTimeout), oktaV5.WithRateLimitMaxRetries(maxRetries), oktaV5.WithRateLimitMaxBackOff(maxBackoff))
		if err != nil {
			return nil, err
		}
		client := oktaV5.NewAPIClient(oktaConfiguratiopn)

		d.ConnectionManager.Cache.Set(sessionCacheKey, client)
		return client, err
	}

	if domain != "" && clientID != "" && privateKey != "" {
		oktaConfiguratiopn, err := oktaV5.NewConfiguration(oktaV5.WithOrgUrl(domain), oktaV5.WithAuthorizationMode("PrivateKey"), oktaV5.WithClientId(clientID), oktaV5.WithPrivateKey(privateKey), oktaV5.WithScopes(scopes), oktaV5.WithRequestTimeout(requestTimeout), oktaV5.WithRateLimitMaxRetries(maxRetries), oktaV5.WithRateLimitMaxBackOff(maxBackoff))
		if err != nil {
			return nil, err
		}

		client := oktaV5.NewAPIClient(oktaConfiguratiopn)

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
	oktaConfiguratiopn, err := oktaV5.NewConfiguration(oktaV5.WithRequestTimeout(requestTimeout), oktaV5.WithRateLimitMaxRetries(maxRetries), oktaV5.WithRateLimitMaxBackOff(maxBackoff))
	if err != nil {
		return nil, err
	}
	client := oktaV5.NewAPIClient(oktaConfiguratiopn)

	// Save session into cache
	d.ConnectionManager.Cache.Set(sessionCacheKey, client)

	return client, err
}

// Retrieve environment variables with default values
func getEnvVarInt64(envVar string, defaultValue int64) (int64, error) {
	if value := os.Getenv(envVar); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue, nil
		} else {
			return 0, fmt.Errorf("error in converting environment value '%s' string type to int64: %v", envVar, err)
		}
	}
	return defaultValue, nil
}

func getEnvVarInt32(envVar string, defaultValue int32) (int32, error) {
	if value := os.Getenv(envVar); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 32); err == nil {
			return int32(intValue), nil
		} else {
			return 0, fmt.Errorf("error in converting environment value '%s' string type to int32: %v", envVar, err)
		}
	}
	return defaultValue, nil
}

// Retrieve Okta configuration values
func getOktaConfigValues(d *plugin.QueryData) (domain, token, clientID, privateKey string, requestTimeout, maxBackoff int64, maxRetries int32, err error) {
	oktaConfig := GetConfig(d.Connection)

	// The default value has been set as per the API doc: https://github.com/okta/okta-sdk-golang?tab=readme-ov-file#environment-variables
	// SDK supported environment variables: https://github.com/okta/okta-sdk-golang/blob/master/okta/config.go#L33-L70
	requestTimeout, err = getEnvVarInt64("OKTA_CLIENT_REQUEST_TIMEOUT", 30)
	maxBackoff, err = getEnvVarInt64("OKTA_CLIENT_RATE_LIMIT_MAX_BACKOFF", 30)
	maxRetries, err = getEnvVarInt32("OKTA_CLIENT_RATE_LIMIT_MAX_RETRIES", 5)

	if oktaConfig.MaxBackoff != nil {
		maxBackoff = *oktaConfig.MaxBackoff
	}
	if oktaConfig.RequestTimeout != nil {
		requestTimeout = *oktaConfig.RequestTimeout
	}
	if oktaConfig.MaxRetries != nil {
		maxRetries = *oktaConfig.MaxRetries
	}

	domain = getStringValue(oktaConfig.Domain, "OKTA_CLIENT_ORGURL")
	token = getStringValue(oktaConfig.Token, "OKTA_CLIENT_TOKEN")
	clientID = getStringValue(oktaConfig.ClientID, "OKTA_CLIENT_CLIENTID")
	privateKey = getStringValue(oktaConfig.PrivateKey, "OKTA_CLIENT_PRIVATEKEY")

	return
}

func getStringValue(configValue *string, envVar string) string {
	if configValue != nil {
		return *configValue
	}
	return os.Getenv(envVar)
}
