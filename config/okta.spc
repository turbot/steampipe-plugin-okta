connection "okta" {
  plugin = "okta"

  # Get your API token from Okta https://developer.okta.com/docs/guides/create-an-api-token/create-the-token/
  # Can also be set with the OKTA_CLIENT_ORGURL environment variable.
  # domain = "https://<your_okta_domain>.okta.com"

  # Okta API token. Can also be set with the OKTA_CLIENT_TOKEN environment variable.
  # token  = "02d0YZgNSJwlNew6lZG-6qGThisisatest-token"

  # Or use an Okta application and the client credentials flow for authenticating: https://developer.okta.com/docs/guides/implement-oauth-for-okta-serviceapp/overview/
  # Can also be set with the OKTA_CLIENT_ORGURL environment variable.
  # domain      = "https://<your_okta_domain>.okta.com"

  # Okta App client id, used with PrivateKey OAuth auth mode. Can also be set with the OKTA_CLIENT_CLIENTID environment variable.
  # client_id   = "0oa10zpa2bo6tAm9Test"

  # Private key value. Can also be set with the OKTA_CLIENT_PRIVATEKEY environment variable. Can also be set with the OKTA_CLIENT_PRIVATEKEY environment variable.
  # private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAK..."

  # The maximum number of attempts (including the initial call) Steampipe will
  # make for failing API calls. Can also be set with the OKTA_CLIENT_RATE_LIMIT_MAX_RETRIES environment variable.
  # Defaults to 5 and must be greater than or equal to 1.
  # max_retries = 5

  # The maximum amount of time to wait on request back off. Can also be set with the OKTA_CLIENT_RATE_LIMIT_MAX_BACKOFF environment variable.
  # Defaults to 30 and must be greater than or equal to 1.
  # max_backoff = 30

  # HTTP request time out in seconds. Can also be set with the OKTA_CLIENT_REQUEST_TIMEOUT environment variable.
  # Defaults to 30 and must be greater than or equal to 1.
  # request_timeout = 30
}
