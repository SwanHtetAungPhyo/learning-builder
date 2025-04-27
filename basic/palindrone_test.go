package main_test

import (
	"strings"
	"testing"
	"unicode"
)

func IsPalindrome(x string) bool {
	if x == " " {
		return true
	}
	var builder strings.Builder
	for _, r := range x {
		if r != ' ' {
			builder.WriteRune(r)
		}
	}
	input := builder.String()
	reverse := ""
	for r := len(input) - 1; r >= 0; r-- {
		reverse += string(input[r])
	}
	return input == reverse
}
func IsPalindrome2(x string) bool {
	if x == " " {
		return true
	}
	var builder strings.Builder
	for _, r := range x {
		if r != ' ' {
			builder.WriteRune(r)
		}
	}
	str := builder.String()
	left, right := 0, len(str)-1
	for left < right {
		if !unicode.IsLetter(rune(str[left])) && !unicode.IsLetter(rune(str[right])) {
			left++
			continue
		}
		if unicode.IsDigit(rune(str[left])) && unicode.IsDigit(rune(str[right])) {
			right--
			return false
		}
		if unicode.ToLower(rune(str[left])) != unicode.ToLower(rune(str[right])) {
			return false
		}
		left++
		right--
	}
	return true
}
func TestIsPalindrome(t *testing.T) {
	testCases := []struct {
		s string
		e bool
	}{
		{"", true},                            // empty string should be a palindrome
		{"a", true},                           // single character is a palindrome
		{"ab", false},                         // "ab" is not a palindrome
		{"racecar", true},                     // "racecar" is a palindrome
		{"hello", false},                      // "hello" is not a palindrome
		{"a man a plan a canal panama", true}, // palindrome with spaces
		{"", true},                            // edge case: empty string
	}

	for _, tc := range testCases {
		t.Run(tc.s, func(t *testing.T) {
			if got := IsPalindrome(tc.s); got != tc.e {
				t.Errorf("IsPalindrome(%q) = %v, want %v", tc.s, got, tc.e)
			}
		})
		t.Run(tc.s, func(t *testing.T) {
			if got := IsPalindrome2(tc.s); got != tc.e {
				t.Errorf("IsPalindrome2(%q) = %v, want %v", tc.s, got, tc.e)
			}
		})
	}
}

func BenchmarkIsPalindrome(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = IsPalindrome("racecar")
	}

}
func BenchmarkIsPalindrome2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = IsPalindrome2("racecar")
	}
}
