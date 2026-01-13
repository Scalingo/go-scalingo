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

func TestUnprocessableEntityError(t *testing.T) {
	t.Run("returns formatted errors list", func(t *testing.T) {
		err := UnprocessableEntity{
			Errors: map[string][]string{
				"plan": {"cannot be downgraded", "not allowed"},
			},
		}

		require.Equal(t, "* plan → cannot be downgraded, not allowed", err.Error())
	})

	t.Run("returns error message when no errors list is present", func(t *testing.T) {
		err := UnprocessableEntity{
			ErrMessage: "Dedicated database addons plan cannot be downgraded",
		}

		require.Equal(t, "Dedicated database addons plan cannot be downgraded", err.Error())
	})
}

func requestFailedError(message string) error {
	return &RequestFailedError{Message: message}
}
