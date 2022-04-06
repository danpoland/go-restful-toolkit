package validation

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type testSchema struct {
	Name string `json:"first_name" validate:"required"`
}

type addressSchema struct {
	Street string `json:"street" validate:"required"`
}

type nestedTestSchema struct {
	Name      string          `json:"first_name" validate:"required"`
	Addresses []addressSchema `json:"addresses" validate:"required,dive"`
}

func (t *testSchema) Validate(ctx context.Context) (*Result, error) {
	return &Result{Success: true}, nil
}

func TestValidateJsonTags(t *testing.T) {
	errs := (validate.Struct(testSchema{})).(validator.ValidationErrors)

	if errs[0].Field() != "first_name" {
		t.Errorf("Validator not using json tag for field.")
	}
}

func TestBind(t *testing.T) {
	type test struct {
		body           string
		expectedStruct testSchema
		expectedResult Result
	}

	tests := []test{
		{
			"{\"first_name\": \"test\"}",
			testSchema{Name: "test"},
			Result{Success: true},
		},
		{
			"{\"first_name\": \"test\"",
			testSchema{},
			Result{SchemaErrs: []SchemaError{{Code: Malformed}}},
		},
	}

	for _, tc := range tests {
		reader := strings.NewReader(tc.body)
		req := httptest.NewRequest("POST", "http://questboard.io/test", reader)

		var ts testSchema
		res := Bind(req, &ts)

		assert.Equal(t, tc.expectedStruct, ts)
		assert.Equal(t, &tc.expectedResult, res)
	}
}

func TestValidate(t *testing.T) {
	type test struct {
		schema         testSchema
		expectedResult Result
		expectedError  error
	}

	tests := []test{
		{
			testSchema{Name: "test"},
			Result{Success: true},
			nil,
		},
		{
			testSchema{},
			Result{
				Success: false,
				FieldErrs: []FieldError{{
					Code:  Required,
					Field: "first_name",
				}},
			},
			nil,
		},
	}

	for _, tc := range tests {
		res, err := Validate(context.TODO(), &tc.schema)

		assert.Equal(t, &tc.expectedResult, res)
		assert.Equal(t, tc.expectedError, err)
	}
}

func TestBindAndValidate(t *testing.T) {
	type test struct {
		body           string
		expectedStruct testSchema
		expectedResult Result
		expectedError  error
	}

	tests := []test{
		{
			"{\"first_name\": \"test\"",
			testSchema{},
			Result{SchemaErrs: []SchemaError{{Code: Malformed}}},
			nil,
		},
	}

	for _, tc := range tests {
		reader := strings.NewReader(tc.body)
		req := httptest.NewRequest("POST", "http://questboard.io/test", reader)

		var ts testSchema
		res, err := BindAndValidate(context.TODO(), req, &ts)

		assert.Equal(t, tc.expectedStruct, ts)
		assert.Equal(t, &tc.expectedResult, res)
		assert.Equal(t, tc.expectedError, err)
	}
}

func TestMarshalValidationErrors(t *testing.T) {
	errs := (validate.Struct(testSchema{})).(validator.ValidationErrors)
	expected := []FieldError{{
		Code:  Required,
		Field: "first_name",
	}}

	assert.Equal(t, expected, MarshalValidationErrors(errs))
}

func TestMarshalValidationErrors_Nested(t *testing.T) {
	schema := nestedTestSchema{
		Name:      "Test",
		Addresses: []addressSchema{{}},
	}
	expected := []FieldError{{
		Code:  Required,
		Field: "addresses[0].street",
	}}
	errs := (validate.Struct(schema)).(validator.ValidationErrors)

	assert.Equal(t, expected, MarshalValidationErrors(errs))
}
