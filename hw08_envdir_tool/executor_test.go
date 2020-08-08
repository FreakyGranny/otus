package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	env := make(Environment)
	returnCode := RunCmd([]string{"uname -a"}, env)
	require.Equal(t, 0, returnCode)
}

func TestRunCmdNoBin(t *testing.T) {
	env := make(Environment)
	returnCode := RunCmd([]string{"/bin/bash", "xxx"}, env)
	require.Equal(t, 127, returnCode)
}

func TestRunCmdSetEnv(t *testing.T) {
	env := make(Environment)
	env["CUSTOM_VAR"] = "VALUE"
	returnCode := RunCmd([]string{"uname -a"}, env)
	require.Equal(t, 0, returnCode)
	
	val, ok := os.LookupEnv("CUSTOM_VAR")
	require.True(t, ok)
	require.Equal(t, "VALUE", val)
}

func TestRunCmdUNSetEnv(t *testing.T) {
	env := make(Environment)
	env["CUSTOM_VAR"] = ""
	os.Setenv("CUSTOM_VAR", "UNSET_ME")

	returnCode := RunCmd([]string{"uname -a"}, env)
	require.Equal(t, 0, returnCode)
	
	_, ok := os.LookupEnv("CUSTOM_VAR")
	require.False(t, ok)
}
