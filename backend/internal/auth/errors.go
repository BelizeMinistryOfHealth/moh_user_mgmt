package auth //nolint: revive

import "fmt"

// AuthError is custom errors thrown by Auth related code
type AuthError struct { //nolint: revive
	Reason string
	Inner  error
}

func (e AuthError) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("auth error: %s: %v", e.Reason, e.Inner)
	}
	return fmt.Sprintf("auth error: %s", e.Reason)
}

func (e AuthError) Unwrap() error {
	return e.Inner
}

// UserError is a custom error emitted by the UserStore
type UserError struct {
	Reason string
	Inner  error
}

func (e UserError) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("user error: %s: %v", e.Reason, e.Inner)
	}
	return fmt.Sprintf("user error: %s", e.Reason)
}

func (e UserError) Unwrap() error {
	return e.Inner
}
