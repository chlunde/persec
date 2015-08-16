package main

import (
	"bytes"
	"testing"
	"time"
)

func testSingle(t *testing.T, a, b string, elapsed time.Duration, expected string) {
	out := &bytes.Buffer{}
	didPrint := make(map[string]int)
	PerSec([]string{a}, []string{b}, time.Second, didPrint, out)
	if out.String() != expected {
		t.Errorf("PerSec([%s], [%s], %v, %v, ...) ->\nGOT:      %vEXPECTED: %v",
			a, b, elapsed, didPrint, out.String(), expected)
	}
}

func TestNoChange(t *testing.T) {
	testSingle(t, "a 1 b", "a 1 b", time.Second, "")
}

func TestSimple(t *testing.T) {
	testSingle(t, "1", "2", time.Second, "1.00/s\n")
}

func TestSuffix(t *testing.T) {
	testSingle(t, "1 b", "2 b", time.Second, "1.00/s b\n")
}

func TestPrefix(t *testing.T) {
	testSingle(t, "a 1", "a 2", time.Second, "a 1.00/s\n")
}

func TestMiddle(t *testing.T) {
	testSingle(t, "a 1 b", "a 2 b", time.Second, "a 1.00/s b\n")
}

func TestMultiDigit(t *testing.T) {
	testSingle(t, "a 1 b 0", "a 2 b 4", time.Second, "a 1.00/s b 4.00/s\n")
}

func TestCollectDigits(t *testing.T) {
	tests := []struct {
		str    string
		n      int
		length int
	}{
		//{"foo", 0, 0},
		{"2foo", 2, 1},
		{"22foo", 22, 2},
	}
	for _, test := range tests {
		n, length := collectDigits(test.str)
		if n != test.n || length != test.length {
			t.Error("collectDigits(%v) -> (%v,%v), expected (%v,%v)",
				test.str, n, length, test.n, test.length)
		}
	}
}

/*
func TestDigitsInName(t *testing.T) {
	testSingle(t, "a 1 IPv4", "a 2 IPv4", time.Second, "a 1.00/s IPv4\n")
}
*/
