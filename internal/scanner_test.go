package scanner

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewScanner(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "basic configuration",
			config: Config{
				Hostname:  "localhost",
				StartPort: 1,
				EndPort:   100,
				Timeout:   "1s",
				Workers:   10,
				Format:    "text",
				Verbose:   false,
			},
		},
		{
			name: "json output configuration",
			config: Config{
				Hostname:  "example.com",
				StartPort: 80,
				EndPort:   443,
				Timeout:   "2s",
				Workers:   50,
				Format:    "json",
				Verbose:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewScanner(tt.config)

			assert.NotNil(t, scanner)
			assert.Equal(t, tt.config.Hostname, scanner.config.Hostname)
			assert.Equal(t, tt.config.StartPort, scanner.config.StartPort)
			assert.Equal(t, tt.config.EndPort, scanner.config.EndPort)
			assert.Equal(t, tt.config.Workers, scanner.config.Workers)
			assert.NotNil(t, scanner.results)
			assert.Len(t, scanner.results, 0)
		})
	}
}

func TestIdentifyService(t *testing.T) {
	scanner := NewScanner(Config{})

	tests := []struct {
		name     string
		port     int
		banner   string
		expected string
	}{
		{
			name:     "SSH port with SSH banner",
			port:     22,
			banner:   "SSH-2.0-OpenSSH_8.0",
			expected: "SSH",
		},
		{
			name:     "HTTP port with HTTP banner",
			port:     80,
			banner:   "HTTP/1.1 200 OK",
			expected: "HTTP",
		},
		{
			name:     "Known port without banner",
			port:     443,
			banner:   "",
			expected: "HTTPS",
		},
		{
			name:     "Unknown port with MySQL banner",
			port:     3307,
			banner:   "mysql_native_password",
			expected: "MySQL",
		},
		{
			name:     "Unknown port with unknown banner",
			port:     9999,
			banner:   "custom-service",
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scanner.identifyService(tt.port, tt.banner)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAddResult(t *testing.T) {
	scanner := NewScanner(Config{})

	result1 := Result{
		Hostname: "localhost",
		Port:     80,
		Status:   "open",
		Service:  "HTTP",
	}

	result2 := Result{
		Hostname: "localhost",
		Port:     443,
		Status:   "closed",
	}

	scanner.addResult(result1)
	scanner.addResult(result2)

	assert.Len(t, scanner.results, 2)
	assert.Equal(t, result1, scanner.results[0])
	assert.Equal(t, result2, scanner.results[1])
}

func TestValidateTimeout(t *testing.T) {
	tests := []struct {
		name        string
		timeout     string
		shouldError bool
	}{
		{
			name:        "valid duration seconds",
			timeout:     "1s",
			shouldError: false,
		},
		{
			name:        "valid duration milliseconds",
			timeout:     "500ms",
			shouldError: false,
		},
		{
			name:        "valid duration minutes",
			timeout:     "1m",
			shouldError: false,
		},
		{
			name:        "invalid duration",
			timeout:     "invalid",
			shouldError: true,
		},
		{
			name:        "empty duration",
			timeout:     "",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := time.ParseDuration(tt.timeout)
			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRun(t *testing.T) {
	t.Run("invalid timeout format", func(t *testing.T) {
		config := Config{
			Hostname:  "localhost",
			StartPort: 80,
			EndPort:   80,
			Timeout:   "invalid",
			Workers:   1,
			Format:    "text",
		}

		err := Run(config)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid timeout format")
	})
}
