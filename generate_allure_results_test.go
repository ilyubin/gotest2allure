package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateAllureResults(t *testing.T) {
	t.Run("test 1 passed", func(t *testing.T) {
		t.Log("step a")
		t.Log("step b")
		assert.True(t, true)
	})

	t.Run("test 2 failed", func(t *testing.T) {
		t.Log("step c")
		assert.True(t, true)
	})

	t.Run("test 3 passed", func(t *testing.T) {
		t.Log("step d")
		t.Log("step e")
		assert.True(t, true)
	})

	t.Run("test 4 passed", func(t *testing.T) {
		t.Log("step f")
		assert.True(t, true)
	})
	t.Log("")
}
