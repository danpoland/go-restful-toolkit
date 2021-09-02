package validation

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Use the json tags as the field names in FieldErrors
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Validatable defines an interface for structs that can be validated, returning a Result.
type Validatable interface {
	Validate(ctx context.Context) (*Result, error)
}

// Bind decodes a http request's JSON body to the provided validatable interface and returns a Result
// indicating if the decode process was successful. If the decode process fails the Result will contain a SchemaError
// with the Malformed ErrorCode.
func Bind(r *http.Request, schema interface{}) *Result {
	if err := json.NewDecoder(r.Body).Decode(schema); err != nil {
		return &Result{SchemaErrs: []SchemaError{{Code: Malformed}}}
	}
	return &Result{Success: true}
}

// ValidateFields validates the field validation tags on the provided schema, returning a Result.
func ValidateFields(schema interface{}) *Result {
	if err := validate.Struct(schema); err != nil {
		fe := MarshalValidationErrors(err.(validator.ValidationErrors))
		return &Result{FieldErrs: fe}
	}
	return &Result{Success: true}
}

// Validate checks for field errors on the provided Validatable interface using validator and calls the
// Validate function on the validatable struct.
func Validate(ctx context.Context, v Validatable) (*Result, error) {
	if r := ValidateFields(v); !r.Success {
		return r, nil
	}
	return v.Validate(ctx)
}

// BindAndValidate binds a request's JSON body to the provided Validatable interface, validates
// the fields on the schema using validator, and finally calls validate on the Validatable schema.
func BindAndValidate(ctx context.Context, r *http.Request, schema Validatable) (*Result, error) {
	if res := Bind(r, schema); !res.Success {
		return res, nil
	}
	return Validate(ctx, schema)
}

// MarshalValidationErrors converts validator ValidationErrors into field errors,
// using the tag as the Code attribute.
func MarshalValidationErrors(errs validator.ValidationErrors) []FieldError {
	var fErrs []FieldError

	for _, ve := range errs {
		fErrs = append(fErrs, FieldError{
			Field: ve.Field(),
			Code:  ve.Tag(),
		})
	}

	return fErrs
}
