package config

import "testing"

func TestParseRejectsInvalidConfig(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		env  map[string]string
	}{
		{name: "missing port", env: map[string]string{}},
		{name: "blank port", env: map[string]string{"HTTP_PORT": "  "}},
		{name: "tls cert missing", env: map[string]string{"HTTP_PORT": "8081", "DEBUG": "false"}},
		{name: "tls key missing", env: map[string]string{"HTTP_PORT": "8081", "SSL_CERT_PATH": "/certs/fullchain.pem"}},
		{name: "debug unset still needs tls", env: map[string]string{"HTTP_PORT": "8081", "SSL_CERT_PATH": "/certs/fullchain.pem"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := Parse(func(key string) string { return tc.env[key] })
			if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestParseDebugSkipsTLS(t *testing.T) {
	t.Parallel()
	cfg, err := Parse(func(key string) string {
		if key == "HTTP_PORT" {
			return "8081"
		}
		if key == "DEBUG" {
			return "true"
		}
		return ""
	})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.Debug || cfg.Port != "8081" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}
