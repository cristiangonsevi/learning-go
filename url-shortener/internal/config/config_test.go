package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name      string
		env       map[string]string
		wantError bool
	}{{
		name: "Valid config",
		env: map[string]string{
			"DB_HOST": "localhost",
			"DB_NAME": "url_shortener_db",
			"DB_USER": "postgres",
			"DB_PASS": "postgres",
		},
		wantError: false,
	},
		{
			name: "invalid config",
			env: map[string]string{
				"DB_HOST": "localhost",
				"DB_NAME": "url_shortener_db",
			},
			wantError: true,
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			_, err := Load()

			if (err != nil) != tt.wantError {
				t.Errorf("Erorr al probar el test %v, %v", tt.name, err)
			}
		})
	}
}
