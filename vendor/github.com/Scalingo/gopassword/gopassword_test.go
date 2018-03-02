package gopassword

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGenerate(t *testing.T) {
	Convey("When we want to generate a password", t, func() {
		Convey("By default, it should be 20 characters", func() {
			So(len(Generate()), ShouldEqual, 20)
		})

		Convey("With an argument, the generated password should have its length", func() {
			So(len(Generate(42)), ShouldEqual, 42)
		})

		Convey("With several arguments, only the first should be considered", func() {
			So(len(Generate(10, 20, 30)), ShouldEqual, 10)
		})
	})

	Convey("Given a generated password", t, func() {
		passwd := Generate(20)
		Convey("The character frequency should be low", func() {
			fm := frequencyMap(passwd)
			maxFreq := max(fm)
			So(maxFreq, ShouldBeLessThanOrEqualTo, 3)
		})
	})
}

func frequencyMap(p string) map[rune]int {
	charMap := make(map[rune]int)
	for _, r := range p {
		charMap[r] += 1
	}
	return charMap
}

func max(m map[rune]int) int {
	res := 0
	for _, v := range m {
		if v > res {
			res = v
		}
	}
	return res
}
