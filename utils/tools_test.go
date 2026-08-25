package utils

import (
	"github.com/tidwall/gjson"
	"testing"
)

func TestEnsureBool(t *testing.T) {
	tests := []struct {
		input      interface{}
		defaultVal bool
		expected   bool
	}{
		{true, false, true},
		{false, true, false},
		{"true", false, true},
		{"TRUE", false, true},
		{"yes", false, true},
		{"1", false, true},
		{"false", true, false},
		{"FALSE", true, false},
		{"no", true, false},
		{"0", true, false},
		{"invalid", true, true},
		{"invalid", false, false},
		{gjson.Parse(`{"val": true}`).Get("val"), false, true},
		{gjson.Parse(`{"val": false}`).Get("val"), true, false},
		{gjson.Parse(`{"val": "yes"}`).Get("val"), false, true},
		{gjson.Parse(`{}`).Get("val"), true, true},
	}

	for _, tt := range tests {
		got := EnsureBool(tt.input, tt.defaultVal)
		if got != tt.expected {
			t.Errorf("EnsureBool(%v, %v) = %v; want %v", tt.input, tt.defaultVal, got, tt.expected)
		}
	}
}

func TestInArrayString(t *testing.T) {
	arr := []string{"apple", "banana", "cherry"}
	if !InArrayString("banana", arr) {
		t.Errorf("expected 'banana' in array")
	}
	if InArrayString("orange", arr) {
		t.Errorf("expected 'orange' not in array")
	}
}

func TestInArryInt(t *testing.T) {
	arr := []int{1, 5, 10, 20}
	if !InArryInt(10, arr) {
		t.Errorf("expected 10 in array")
	}
	if InArryInt(7, arr) {
		t.Errorf("expected 7 not in array")
	}
}
