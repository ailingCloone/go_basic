package httperror

import (
	"encoding/json"
	"net/http"

	validator "github.com/go-playground/validator/v10"
)

const (
	// List of status code use for our end
	EmptyFieldCode          = "SH0001" // Empty for the field
	InvalidFormatCode       = "SH0002" // Parameters is invalid format
	MinCharactersCode       = "SH0010" // Mininum characters is required
	MaxCharactersCode       = "SH0011" // Maximum characters is required
	PayloadSizeCode         = "SH0020" // Payload size limitation
	InternalServerErrorCode = "SH1000" // Some internal server error on our end
	CustomMessage           = "SH1200" // We return custom error base on specify API, some api will return different response.

	// Contains all the parameters we use for request when calling the APIs

	File        = "file"
	AppVersion  = "app_version"
	Guid        = "user_id"
	P1no        = "p1no"
	FullName    = "full_name"
	FileAction  = "file_action"
	Email       = "email"
	Phase       = "phase"
	Platform    = "platform"
	Phase3      = "phase3"
	Phase2      = "phase2"
	Icno        = "icno"
	AccessToken = "access_token"
	Mobile      = "mobile"
)

var (
	// Contain general error msg

	// 400 - 499 Client error responses
	InvalidFileExtensions = "Invalid file."
	FileTooLarge          = "File too large."
	FieldEmpty            = "Field must not be empty."
	MinCharacters         = "Mininum characters not reached."
	MaxCharacters         = "Maximum characters reached."
	InvalidFieldFormat    = "Invalid field format."
	PayloadSize           = "Payload too large."

	// >= 500  Server error responses
	InternalServerError = "Internal server error."
)

var singletonHTTPErrorInstance *SingletonError
var routesMap map[string]RouteErrorInfo

// once is used for synchronization during initialization.
var once singletonOnce

type HttpError interface {
	ErrorEmptyFieldCode()
}

// Singleton represents a singleton object for handling HTTP errors.
type SingletonError struct {
	message string
	status  int
}

// ErrorResponse holds information about an error response.
type ErrorResponse struct {
	Errors    interface{}
	Status    int
	ErrorCode string
}

// RouteErrorInfo holds information about the possible error responses for a specific route.
type RouteErrorInfo struct {
	ErrorResponses []ErrorResponse
}

type singletonOnce struct {
	done bool
}

// Do is a method of singletonOnce to ensure initialization is done only once.
func (o *singletonOnce) Do(f func()) {
	if !o.done {
		f()
		o.done = true
	}
}

// GetSingletonInstance returns the singleton instance of the Singleton class.
func GetSingletonInstance() *SingletonError {
	once.Do(func() {
		// Initialize the singletonInstance only once.

		singletonHTTPErrorInstance = &SingletonError{}

		//singletonHTTPErrorInstance.init()
	})
	return singletonHTTPErrorInstance
}

// SetError updates the error message and status of the singleton instance.
//func (s *SingletonError) BuildError(route, parameters string, status int, errors interface{}) (*SingletonError) {
//	if routeErrorInfo, ok := routesMap[route]; ok {
//		routeErrorInfo.ErrorResponses = append(routeErrorInfo.ErrorResponses, ErrorResponse{
//			Errors: errors,
//			Status: status,
//			Parameters: parameters,
//		})
//		routesMap[route] = routeErrorInfo
//		return s
//	}
//
//
//	var new RouteErrorInfo
//	new.ErrorResponses = append(new.ErrorResponses, ErrorResponse{
//		Errors: errors,
//		Status: status,
//		Parameters: parameters,
//	})
//	routesMap[route] = new
//
//	return s
//}

func (s *SingletonError) init() {

	//	errors := make(map[string]interface{})
	//	errors[EmptyFieldCode] = make(map[string]ErrorResponse)
	//	errors[InvalidFormatCode] = make(map[string]ErrorResponse)
	//	errors[MinCharactersCode] = make(map[string]ErrorResponse)
	//	errors[MaxCharactersCode] = make(map[string]ErrorResponse)
	//	errors[PayloadSize] = make(map[string]ErrorResponse)
	//	errors[InternalServerErrorCode] = make(map[string]ErrorResponse)

	routesMap = make(map[string]RouteErrorInfo)

	// Below is the custom error response for each
	routesMap["/v2/profile/update-user-info"] = RouteErrorInfo{
		ErrorResponses: []ErrorResponse{
			{Errors: "Invalid file parameters.", Status: http.StatusBadRequest, ErrorCode: CustomMessage},
			{Errors: "The file format not supported.", Status: http.StatusUnsupportedMediaType, ErrorCode: CustomMessage},
			{Errors: "File size too large.", Status: http.StatusUnsupportedMediaType, ErrorCode: CustomMessage},
		},
	}
}

func ResponseErrorWithValidationErrors(w http.ResponseWriter, err validator.ValidationErrors) {

	for _, err := range err {
		switch err.Tag() {
		case "typeValidator", "pageValidator", "registerListStatusValidator", "fromValidator":
			ErrorInvalidFormat(w)
			return
		case "required":
			ErrorEmptyField(w)
			return
		case "min":
			ErrorMinCharacters(w)
			return
		case "max":
			ErrorMaxCharacters(w)
			return
		}
	}
}

const (
	appVersionField  = "field \"" + AppVersion + "\" is not set"
	platformField    = "field \"" + Platform + "\" is not set"
	accessTokenField = "field \"" + AccessToken + "\" is not set"
	p1NoField        = "field \"" + P1no + "\" is not set"
	mobileField      = "field \"" + Mobile + "\" is not set"
	guidField        = "field \"" + Guid + "\" is not set"
)

var errorMessages = map[string]func(http.ResponseWriter){
	appVersionField:  ErrorEmptyField,
	platformField:    ErrorEmptyField,
	accessTokenField: ErrorEmptyField,
	p1NoField:        ErrorEmptyField,
	mobileField:      ErrorEmptyField,
	guidField:        ErrorEmptyField,
}

func ResponseErrorWithGoZeroErrors(w http.ResponseWriter, err error) {
	if handleError, ok := errorMessages[err.Error()]; ok {
		handleError(w)
		return
	}
}

func ErrorEmptyField(w http.ResponseWriter) {
	err := ErrorResponse{
		Errors:    FieldEmpty,
		Status:    http.StatusBadRequest,
		ErrorCode: EmptyFieldCode,
	}
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(err)
}

func ErrorInvalidFormat(w http.ResponseWriter) {
	err := ErrorResponse{
		Errors:    InvalidFieldFormat,
		Status:    http.StatusBadRequest,
		ErrorCode: InvalidFormatCode,
	}
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(err)
}
func ErrorMinCharacters(w http.ResponseWriter) {
	err := ErrorResponse{
		Errors:    MinCharacters,
		Status:    http.StatusBadRequest,
		ErrorCode: MinCharactersCode,
	}
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(err)
}
func ErrorMaxCharacters(w http.ResponseWriter) {
	err := ErrorResponse{
		Errors:    MaxCharacters,
		Status:    http.StatusBadRequest,
		ErrorCode: MaxCharactersCode,
	}
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(err)
}
func ErrorPayloadSize(w http.ResponseWriter, parameters string) {
	err := ErrorResponse{
		Errors:    PayloadSize,
		Status:    http.StatusBadRequest,
		ErrorCode: PayloadSizeCode,
	}
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(err)
}
func ErrorInternalServerError(w http.ResponseWriter) {
	err := ErrorResponse{
		Errors:    InternalServerError,
		Status:    http.StatusBadRequest,
		ErrorCode: InternalServerErrorCode,
	}
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(err)
}

func (s *SingletonError) ResponseError(w http.ResponseWriter, route string, index int) {
	var errorResponse ErrorResponse

	// Check if the route exists in the map
	if routeErrorInfo, ok := routesMap[route]; ok {
		// Check if the specified error index is within bounds
		if index >= 0 && index < len(routeErrorInfo.ErrorResponses) {
			errorResponse = routeErrorInfo.ErrorResponses[index]
			w.WriteHeader(errorResponse.Status)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(errorResponse)
}
