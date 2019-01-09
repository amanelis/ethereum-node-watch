package main

import (
	"testing"
    "github.com/stretchr/testify/assert"
)

func TestConstants(t *testing.T) {
	assert.Equal(t, environment, "ENVIRONMENT")
}
