package scalingo

import (
	"testing"
	"time"
)

type isDeprecatedTest struct {
	date       time.Time
	deprecated bool
}

var isDeprecatedTests = map[string]isDeprecatedTest{
	"test isDeprecated is false when deprecation date is null|nil|zero": isDeprecatedTest{time.Time{}, false},
	"test isDeprecated is true when deprecation date is today's date":   isDeprecatedTest{time.Now(), true},
	"test isDeprecated is true when deprecation date is in the past":    isDeprecatedTest{time.Now().AddDate(0, 0, -1), true},
	"test isDeprecated is false when deprecation date is in the future": isDeprecatedTest{time.Now().AddDate(0, 0, 1), false},
}

func TestStackIsDeprecated(t *testing.T) {
	for message, test := range isDeprecatedTests {
		t.Run(message, func(t *testing.T) {
			stack := Stack{DeprecatedAt: test.date}
			if stack.IsDeprecated() != test.deprecated {
				t.Errorf("IsDeprecated expected to be %t, got %t", test.deprecated, stack.IsDeprecated())
			}
		})
	}
}
