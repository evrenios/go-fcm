package fcm

import (
	"encoding/json"
	"errors"
)

var (
	// ErrMissingRegistration occurs if registration token is not set.
	ErrMissingRegistration = errors.New("MissingRegistration")

	// ErrInvalidRegistration occurs if registration token is invalid.
	ErrInvalidRegistration = errors.New("InvalidRegistration")

	// ErrNotRegistered occurs when application was deleted from device and
	// token is not registered in FCM.
	ErrNotRegistered = errors.New("NotRegistered")

	// ErrInvalidPackageName occurs if package name in message is invalid.
	ErrInvalidPackageName = errors.New("InvalidPackageName")

	// ErrMismatchSenderID occurs when application has a new registration token.
	ErrMismatchSenderID = errors.New("MismatchSenderId")

	// ErrMessageTooBig occurs when message is too big.
	ErrMessageTooBig = errors.New("MessageTooBig")

	// ErrInvalidDataKey occurs if data key is invalid.
	ErrInvalidDataKey = errors.New("InvalidDataKey")

	// ErrInvalidTTL occurs when message has invalid TTL.
	ErrInvalidTTL = errors.New("InvalidTTL")

	// ErrUnavailable occurs when FCM service is unavailable. It makes sense
	// to retry after this error.
	ErrUnavailable = connectionError("Unavailable")

	// ErrInternalServerError is internal FCM error. It makes sense to retry
	// after this error.
	ErrInternalServerError = serverError("InternalServerError")

	// ErrDeviceMessageRateExceeded occurs when client sent to many requests to
	// the device.
	ErrDeviceMessageRateExceeded = errors.New("DeviceMessageRateExceeded")

	// ErrTopicsMessageRateExceeded occurs when client sent to many requests to
	// the topics.
	ErrTopicsMessageRateExceeded = errors.New("TopicsMessageRateExceeded")

	// ErrInvalidParameters occurs when provided parameters have the right name and type
	ErrInvalidParameters = errors.New("InvalidParameters")

	// ErrUnknown for unknown error type
	ErrUnknown = errors.New("Unknown")

	// ErrInvalidApnsCredential for Invalid APNs credentials
	ErrInvalidApnsCredential = errors.New("InvalidApnsCredential")
)

var (
	errMap = map[string]error{
		"MissingRegistration":       ErrMissingRegistration,
		"InvalidRegistration":       ErrInvalidRegistration,
		"NotRegistered":             ErrNotRegistered,
		"InvalidPackageName":        ErrInvalidPackageName,
		"MismatchSenderId":          ErrMismatchSenderID,
		"MessageTooBig":             ErrMessageTooBig,
		"InvalidDataKey":            ErrInvalidDataKey,
		"InvalidTtl":                ErrInvalidTTL,
		"Unavailable":               ErrUnavailable,
		"InternalServerError":       ErrInternalServerError,
		"DeviceMessageRateExceeded": ErrDeviceMessageRateExceeded,
		"TopicsMessageRateExceeded": ErrTopicsMessageRateExceeded,
		"InvalidParameters":         ErrInvalidParameters,
		"InvalidApnsCredential":     ErrInvalidApnsCredential,
	}
)

// connectionError represents connection errors such as timeout error, etc.
// Implements `net.Error` interface.
type connectionError string

func (err connectionError) Error() string {
	return string(err)
}

func (err connectionError) Temporary() bool {
	return true
}

func (err connectionError) Timeout() bool {
	return true
}

// serverError represents internal server errors.
// Implements `net.Error` interface.
type serverError string

func (err serverError) Error() string {
	return string(err)
}

func (serverError) Temporary() bool {
	return true
}

func (serverError) Timeout() bool {
	return false
}

// Response represents the FCM server's response to the application
// server's sent message.
type Response struct {
	MulticastID  int64    `json:"multicast_id"`
	Success      int      `json:"success"`
	Failure      int      `json:"failure"`
	CanonicalIDs int      `json:"canonical_ids"`
	Results      []Result `json:"results"`

	// Device Group HTTP Response
	FailedRegistrationIDs []string `json:"failed_registration_ids"`

	// Topic HTTP response
	MessageID int64 `json:"message_id"`
	Error     error `json:"error"`
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (r *Response) UnmarshalJSON(data []byte) error {
	var response struct {
		MulticastID  int64    `json:"multicast_id"`
		Success      int      `json:"success"`
		Failure      int      `json:"failure"`
		CanonicalIDs int      `json:"canonical_ids"`
		Results      []Result `json:"results"`

		// Device Group HTTP Response
		FailedRegistrationIDs []string `json:"failed_registration_ids"`

		// Topic HTTP response
		MessageID int64  `json:"message_id"`
		Error     string `json:"error"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}

	r.MulticastID = response.MulticastID
	r.Success = response.Success
	r.Failure = response.Failure
	r.CanonicalIDs = response.CanonicalIDs
	r.Results = response.Results
	r.Success = response.Success
	r.FailedRegistrationIDs = response.FailedRegistrationIDs
	r.MessageID = response.MessageID
	if response.Error != "" {
		if val, ok := errMap[response.Error]; ok {
			r.Error = val
		} else {
			r.Error = ErrUnknown
		}
	}

	return nil
}

// Result represents the status of a processed message.
type Result struct {
	MessageID      string `json:"message_id"`
	RegistrationID string `json:"registration_id"`
	Error          error  `json:"error"`
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (r *Result) UnmarshalJSON(data []byte) error {
	var result struct {
		MessageID      string `json:"message_id"`
		RegistrationID string `json:"registration_id"`
		Error          string `json:"error"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	r.MessageID = result.MessageID
	r.RegistrationID = result.RegistrationID
	if result.Error != "" {
		if val, ok := errMap[result.Error]; ok {
			r.Error = val
		} else {
			r.Error = ErrUnknown
		}
	}

	return nil
}

// Unregistered checks if the device token is unregistered,
// according to response from FCM server. Useful to determine
// if app is uninstalled.
func (r Result) Unregistered() bool {
	switch r.Error {
	case ErrNotRegistered, ErrMismatchSenderID, ErrMissingRegistration, ErrInvalidRegistration:
		return true

	default:
		return false
	}
}
