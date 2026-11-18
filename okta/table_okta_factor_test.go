package okta

import "testing"

func TestParseFactor(t *testing.T) {
	cases := []struct {
		name, raw, factorType, provider, status string
		hasProfile                              bool
	}{
		{
			// Okta Verify with FastPass: not modelled by the SDK's factor union.
			name:       "signed_nonce",
			raw:        `{"id":"uft1","factorType":"signed_nonce","provider":"OKTA","vendorName":"OKTA","status":"ACTIVE","created":"2024-01-02T03:04:05.000Z","profile":{"name":"Laptop","platform":"MACOS"},"_links":{"self":{"href":"https://example.okta.com/x"}}}`,
			factorType: "signed_nonce", provider: "OKTA", status: "ACTIVE", hasProfile: true,
		},
		{
			name:       "push",
			raw:        `{"id":"opf1","factorType":"push","provider":"OKTA","status":"ACTIVE","profile":{"deviceType":"SmartPhone_IPhone"}}`,
			factorType: "push", provider: "OKTA", status: "ACTIVE", hasProfile: true,
		},
		{
			name:       "no profile",
			raw:        `{"id":"x","factorType":"token:hardware","provider":"YUBICO","status":"PENDING_ACTIVATION"}`,
			factorType: "token:hardware", provider: "YUBICO", status: "PENDING_ACTIVATION",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f, err := parseFactor([]byte(c.raw))
			if err != nil {
				t.Fatal(err)
			}
			if got := f.GetFactorType(); got != c.factorType {
				t.Errorf("factorType = %q, want %q", got, c.factorType)
			}
			if got := f.GetProvider(); got != c.provider {
				t.Errorf("provider = %q, want %q", got, c.provider)
			}
			if got := f.GetStatus(); got != c.status {
				t.Errorf("status = %q, want %q", got, c.status)
			}
			if (f.Profile != nil) != c.hasProfile {
				t.Errorf("profile set = %v, want %v", f.Profile != nil, c.hasProfile)
			}
		})
	}
}
