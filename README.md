# restful-toolkit
This is a Go library that provides common, non-domain specific, utilities for the building RESTFul APIs.

## Packages

### render
Contains shortcut methods for rendering common HTTP responses.

## validation
The validation package is used for schema validation, supporting both field level validation and schema level validation.
The package defines methods for binding incoming JSON request bodies and URL query parameters to structs and validating
those structs using the [validator package](https://github.com/go-playground/validator), producing a standardized Result schema.
