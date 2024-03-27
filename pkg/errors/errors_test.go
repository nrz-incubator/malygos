package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NotFoundError(t *testing.T) {
	err := NewNotFoundError("thing", "test")
	assert.True(t, IsNotFound(err))
	assert.False(t, IsConflict(err))
	assert.False(t, IsNotFound(nil))
	assert.False(t, IsNotFound(fmt.Errorf("test")))
	assert.Equal(t, "thing test not found", err.Error())
}

func Test_ConflictError(t *testing.T) {
	err := NewConflictError("thing", "test")
	assert.True(t, IsConflict(err))
	assert.False(t, IsNotFound(err))
	assert.False(t, IsConflict(nil))
	assert.False(t, IsConflict(fmt.Errorf("test")))
	assert.Equal(t, "conflict: thing test already exists", err.Error())
}
