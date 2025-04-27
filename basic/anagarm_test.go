package main_test

import (
	"sort"
	"testing"
)

// Mem huge
func Anagram(s string, t string) bool {
	var arrOne = make([]rune, len(s))
	var arrTwo = make([]rune, len(t))
	for _, r := range s {
		arrOne = append(arrOne, r)
	}
	for _, r := range t {
		arrTwo = append(arrTwo, r)
	}
	sort.Slice(arrOne, func(i, j int) bool {
		return arrOne[i] < arrOne[j]
	})
	sort.Slice(arrTwo, func(i, j int) bool {
		return arrTwo[i] < arrTwo[j]
	})
	return string(arrOne) == string(arrTwo)
}
func Anagram2(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}
	var mapping = make(map[rune]int, len(s))
	for _, r := range s {
		mapping[r]++
	}
	for _, r := range t {
		mapping[r]--
		if mapping[r] < 0 {
			return false
		}
	}
	return true
}
func TestAnagram(t *testing.T) {
	testCases := []struct {
		s string
		t string
		e bool
	}{
		{"anagram", "nagaram", true},
		{"rat", "car", false},
	}
	for _, tc := range testCases {
		if Anagram(tc.s, tc.t) != tc.e {
			t.Errorf("Anagram(%q, %q) = %v, want %v", tc.s, tc.t, Anagram(tc.s, tc.t), tc.e)
		}
	}
	for _, tc := range testCases {
		if Anagram2(tc.s, tc.t) != tc.e {
			t.Errorf("Anagram2(%q, %q) = %v, want %v", tc.s, tc.t, Anagram2(tc.s, tc.t), tc.e)
		}
	}
}

func BenchmarkAnagram(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Anagram("anagram", "nagaram")
	}
}

func BenchmarkAnagram2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Anagram2("anagram", "nagaram")
	}
}
