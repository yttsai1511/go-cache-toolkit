package cache

import (
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {
	Set("testValue", "key1", "key2")

	value, err := Get[string]("key1", "key2")
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if *value != "testValue" {
		t.Errorf("Expected value 'testValue', but got: %v", *value)
	}
}

func TestExpiration(t *testing.T) {
	SetWithTTL("shortLived", 1*time.Second, "expiringKey")

	_, err := Get[string]("expiringKey")
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	time.Sleep(1500 * time.Millisecond)

	_, err = Get[string]("expiringKey")
	if err == nil {
		t.Errorf("Expected an error for expired key, but got none")
	}
}

func TestDelete(t *testing.T) {
	Set("toBeDeleted", "deleteKey")

	Delete("deleteKey")

	_, err := Get[string]("deleteKey")
	if err == nil {
		t.Errorf("Expected an error for deleted key, but got none")
	}
}

func TestDefaultTTL(t *testing.T) {
	SetDefaultTTL(500 * time.Millisecond)

	Set("defaultTTLValue", "defaultKey")

	_, err := Get[string]("defaultKey")
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	time.Sleep(1 * time.Second)

	_, err = Get[string]("defaultKey")
	if err == nil {
		t.Errorf("Expected an error for expired key, but got none")
	}
}

func TestGenerateKey(t *testing.T) {
	key := GenerateKey("part1", 123, true)

	expected := "part1|123|true"
	if key != expected {
		t.Errorf("Expected key '%s', but got '%s'", expected, key)
	}
}

func TestCleanExpired(t *testing.T) {
	SetWithTTL("value1", 100*time.Millisecond, "key1")
	SetWithTTL("value2", 5*time.Second, "key2")

	time.Sleep(200 * time.Millisecond)

	CleanExpired(time.Now())

	_, err := Get[string]("key1")
	if err == nil {
		t.Errorf("Expected an error for expired key 'key1', but got none")
	}

	_, err = Get[string]("key2")
	if err != nil {
		t.Errorf("Expected no error for key 'key2', but got: %v", err)
	}
}
