package common

import "testing"

func TestIsEmpty(t *testing.T) {
	cases := []struct{ input string; want bool }{
		{"", true},
		{"   ", true},
		{"hello", false},
		{"0", false},
	}
	for _, tc := range cases {
		if got := IsEmpty(tc.input); got != tc.want {
			t.Errorf("IsEmpty(%q) = %v; want %v", tc.input, got, tc.want)
		}
	}
}

func TestIsNotEmpty(t *testing.T) {
	if got := IsNotEmpty("hello"); !got {
		t.Error("IsNotEmpty(\"hello\") = false; want true")
	}
	if got := IsNotEmpty(""); got {
		t.Error("IsNotEmpty(\"\") = true; want false")
	}
}
