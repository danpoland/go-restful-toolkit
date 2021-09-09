package validation

// ErrorCode is used to assign a text code to a validation error. Readable by a client to apply custom
// handling or messaging.
type ErrorCode = string

// A set of commonly used error codes.
const (
	Required  ErrorCode = "required"
	NotUnique ErrorCode = "not_unique"
	Malformed ErrorCode = "malformed"
	NotFound  ErrorCode = "not_found"
)

// FieldError represents an error associated with a single field in a validatable struct.
type FieldError struct {
	Field   string    `json:"field"`
	Code    ErrorCode `json:"code,omitempty"`
	Message string    `json:"message,omitempty"`
}

// SchemaError represents an error associated with the entire validatable struct or a collection
// of one or more fields.
type SchemaError struct {
	Fields  []string  `json:"fields,omitempty"`
	Code    ErrorCode `json:"code,omitempty"`
	Message string    `json:"message,omitempty"`
}

// Result indicates if a call to validate was successful and contains the errors
// associated with a failed validation.
type Result struct {
	Success    bool          `json:"success"`
	FieldErrs  []FieldError  `json:"field_errors"`
	SchemaErrs []SchemaError `json:"schema_errors"`
}

// Success is shortcut for creating a new successful validation result.
func Success() *Result {
	return &Result{Success: true}
}

// FieldsFailure is a shortcut for creating a new failed validation result with field errors.
func FieldsFailure(fieldErrors ...FieldError) *Result {
	return &Result{Success: false, FieldErrs: fieldErrors}
}

// SchemaFailure is a shortcut for create a new failed validation result with schema errors.
func SchemaFailure(schemaError ...SchemaError) *Result {
	return &Result{Success: false, SchemaErrs: schemaError}
}
