package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorCodes(t *testing.T) {
	type test struct {
		code     ErrorCode
		expected string
	}

	var tests = []test{
		{Required, "required"},
		{NotUnique, "not_unique"},
		{Malformed, "malformed"},
		{NotFound, "not_found"},
	}

	for _, test := range tests {
		if test.code != test.expected {
			t.Errorf("ErrorCode %v to not equal expected %v", test.code, test.expected)
		}
	}
}

func TestNewSuccess(t *testing.T) {
	assert.Equal(t, &Result{Success: true}, NewSuccess())
}

func TestNewFieldsFailure(t *testing.T) {
	got := NewFieldsFailure(FieldError{Field: "name"}, FieldError{Field: "email"})
	expected := &Result{
		Success:   false,
		FieldErrs: []FieldError{{Field: "name"}, {Field: "email"}},
	}

	assert.Equal(t, expected, got)
}
