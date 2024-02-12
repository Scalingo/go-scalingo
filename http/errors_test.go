package http

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_IsOTPRequired(t *testing.T) {
	cases := map[string]struct {
		expectedValue bool
		givenError    error
	}{
		"returns OTP required when the ErrOTPRequired variable is passed": {
			givenError:    ErrOTPRequired,
			expectedValue: true,
		},
		"returns OTP required when a message 'OTP Required' is returned by API": {
			givenError:    requestFailedError("OTP Required"),
			expectedValue: true,
		},
		"returns false when API does not returns 'OPT Required'": {
			givenError:    requestFailedError("bad request"),
			expectedValue: false,
		},
		"returns false when an normal error is given": {
			givenError:    errors.New("bad request"),
			expectedValue: false,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			// When
			optRequired := IsOTPRequired(c.givenError)

			// Then
			require.Equal(t, c.expectedValue, optRequired)
		})
	}
}

func requestFailedError(message string) error {
	return &RequestFailedError{Message: message}
}
