package apperror

// ErrorCode is a custom alias for string representing our API error enum.
type ErrorCode string

const (
	InternalError       ErrorCode = "INTERNAL_ERROR"
	DataValidationError ErrorCode = "DATA_VALIDATION_ERROR"
	UserAlreadyExists   ErrorCode = "USER_ALREADY_EXISTS"
	Unauthorized        ErrorCode = "UNAUTHORIZED"
	NotFound            ErrorCode = "NOT_FOUND"
	BadRequest          ErrorCode = "BAD_REQUEST"
)
