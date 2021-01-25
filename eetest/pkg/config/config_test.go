package config

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_RetrievesConfigFromFlags(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"cmd",
		"--host=test-host",
		"--port=1001",
	}

	conf, err := NewConfig()
	require.Nil(t, err)
	// Test that flags are appropriately set.
	require.Equal(t, "test-host", conf.Host)
	require.Equal(t, 1001, conf.Port)
	// Test that defaults are appropriately set.
	require.Equal(t, 15, conf.WriteTimeout)
}

func Test_RetrievesConfigFromFlagsWithError(t *testing.T) {
	// Arrange
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"cmd",
		"--port=AAAA",
	}
	_, err := NewConfig()
	require.NotNil(t, err)
}

func Test_RetrievesConfigFromEnvironmentVariable(t *testing.T) {
	// Arrange
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Setenv("HOST", "foo")
	os.Setenv("PORT", "1001")

	os.Args = []string{
		"cmd",
	}

	conf, err := NewConfig()

	require.Nil(t, err)
	// Test that enviroment is pullling
	require.Equal(t, "foo", conf.Host)
	require.Equal(t, 1001, conf.Port)
	// Test that defaults are appropriately set.
	require.Equal(t, 15, conf.WriteTimeout)
}
