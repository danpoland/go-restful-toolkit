package validation

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()
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

// ValidateSuccess provides a Validate function to other structs that need to implement Validatable but
// only require the basic field level validation provided by the validation attribute tags.
type ValidateSuccess struct {
}

// Validate returns a success result.
func (v *ValidateSuccess) Validate(ctx context.Context) (*Result, error) {
	return Success(), nil
}

// BindBody decodes a http request's JSON body to the provided validatable interface and returns a Result
// indicating if the decode process was successful. If the decode process fails the Result will contain a SchemaError
// with the Malformed ErrorCode.
func BindBody(r *http.Request, target interface{}) *Result {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return &Result{SchemaErrs: []SchemaError{{Code: Malformed}}}
	}
	return &Result{Success: true}
}

// BindQuery decodes a http request's URL query parameters to the provided validatable interface and returns a Result
// indicating if the decode process was successful. If the decode process fails the Result will contain FieldError
// errors with the Malformed ErrorCode.
func BindQuery(r *http.Request, target interface{}) *Result {
	err := decoder.Decode(target, r.URL.Query())
	if err == nil {
		return &Result{Success: true}
	}

	var fieldErrors []FieldError

	if mErr, ok := err.(schema.MultiError); ok {
		for _, v := range mErr {
			cErr, _ := v.(schema.ConversionError)
			fieldErrors = append(fieldErrors, FieldError{Field: cErr.Key, Code: Malformed})
		}
	}

	return FieldsFailure(fieldErrors...)
}

// ValidateFields validates the field validation tags on the provided schema, returning a Result.
func ValidateFields(target interface{}) *Result {
	if err := validate.Struct(target); err != nil {
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

// BindBodyAndValidate binds a request's JSON body to the provided Validatable interface, validates
// the fields on the schema using validator, and finally calls validate on the Validatable schema.
func BindBodyAndValidate(ctx context.Context, r *http.Request, schema Validatable) (*Result, error) {
	if res := BindBody(r, schema); !res.Success {
		return res, nil
	}
	return Validate(ctx, schema)
}

// BindQueryAndValidate binds a request's URL query parameters to the provided Validatable interface, validates
// the fields on the schema using validator, and finally calls validate on the Validatable schema.
func BindQueryAndValidate(ctx context.Context, r *http.Request, schema Validatable) (*Result, error) {
	if res := BindQuery(r, schema); !res.Success {
		return res, nil
	}
	return Validate(ctx, schema)
}

// MarshalValidationErrors converts validator ValidationErrors into field errors,
// using the tag as the Code attribute.
func MarshalValidationErrors(errs validator.ValidationErrors) []FieldError {
	var fErrs []FieldError

	for _, ve := range errs {
		parts := strings.SplitN(ve.Namespace(), ".", 2)
		fErrs = append(fErrs, FieldError{
			Field: parts[1],
			Code:  ve.Tag(),
		})
	}

	return fErrs
}
