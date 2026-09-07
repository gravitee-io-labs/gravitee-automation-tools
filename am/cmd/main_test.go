package main

import (
	"testing"

	"github.com/gravitee-io-labs/gravitee-automation-sdks/am/pkg/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommand_DefaultFlags(t *testing.T) {
	cmd := newCommand()

	port, err := cmd.Flags().GetInt("port")
	require.NoError(t, err)
	assert.Equal(t, 8080, port)

	basePath, err := cmd.Flags().GetString("base-path")
	require.NoError(t, err)
	assert.Equal(t, server.BasePath, basePath)
}

func TestCommand_ParseFlags(t *testing.T) {
	cmd := newCommand()

	err := cmd.ParseFlags([]string{"--port", "9090", "--base-path", "/api"})
	require.NoError(t, err)

	port, err := cmd.Flags().GetInt("port")
	require.NoError(t, err)
	assert.Equal(t, 9090, port)

	basePath, err := cmd.Flags().GetString("base-path")
	require.NoError(t, err)
	assert.Equal(t, "/api", basePath)
}
