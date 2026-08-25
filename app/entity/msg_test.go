package entity

import (
	"encoding/json"
	"testing"
)

func TestOK(t *testing.T) {
	msg := OK("test_data")
	if msg["retcode"] != 0 {
		t.Errorf("expected retcode 0, got %v", msg["retcode"])
	}
	if msg["status"] != "ok" {
		t.Errorf("expected status ok, got %v", msg["status"])
	}
	if msg["data"] != "test_data" {
		t.Errorf("expected data 'test_data', got %v", msg["data"])
	}
}

func TestFailed(t *testing.T) {
	msg := Failed(404, "not found")
	if msg["retcode"] != 404 {
		t.Errorf("expected retcode 404, got %v", msg["retcode"])
	}
	if msg["status"] != "failed" {
		t.Errorf("expected status failed, got %v", msg["status"])
	}
	if msg["msg"] != "not found" {
		t.Errorf("expected msg 'not found', got %v", msg["msg"])
	}
}

func TestToJson(t *testing.T) {
	msg := OK("payload")
	jsonStr := msg.ToJson()
	var parsed map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &parsed)
	if err != nil {
		t.Fatalf("ToJson produced invalid json: %v", err)
	}
	if parsed["status"] != "ok" {
		t.Errorf("expected status ok in json, got %v", parsed["status"])
	}
}
