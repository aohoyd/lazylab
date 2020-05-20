package gitlab

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	testToken := "deadbeef"
	config := ClientConfig{defaultBaseURL, testToken}
	c, err := NewClient(config)
	if err != nil {
		t.Errorf("NewClient() error = %v", err)
	}

	expectedURL := defaultBaseURL + defaultVersionedAPIPath

	if c.baseURL.String() != expectedURL {
		t.Errorf("NewClient() URLPath = %s, want = %s", c.baseURL.Path, expectedURL)
	}
	if c.UserAgent != defaultUserAgent {
		t.Errorf("NewClient() UserAgent = %s, want = %s", c.UserAgent, defaultUserAgent)
	}
	if c.token != testToken {
		t.Errorf("NewClient() token = %s, want = %s", c.token, testToken)
	}
}