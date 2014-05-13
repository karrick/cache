package cache

import (
	"testing"
	"time"
)

func TestSecondArgumentIsOk(t *testing.T) {
	cache := NewTTL()
	defer cache.Quit()
	_, ok := cache.Get("key")
	if ok != false {
		t.Errorf("Expected: %#v; Actual: %#v\n", false, ok)
	}
}

func TestCanStoreInformation(t *testing.T) {
	cache := NewTTL()
	defer cache.Quit()
	cache.Set("key", "value", 5*time.Second)
	actual, ok := cache.Get("key")
	if ok != true {
		t.Errorf("Expected: %#v; Actual: %#v\n", true, ok)
	}
	if actual.(string) != "value" {
		t.Errorf("Expected: %#v; Actual: %#v\n", "value", actual)
	}
}

func TestAllowsSettingWithNegativeTTL(t *testing.T) {
	cache := NewTTL()
	defer cache.Quit()
	cache.Set("key", "value", -5)
	value, ok := cache.Get("key")
	if ok != false {
		t.Errorf("Expected: %#v; Actual: %#v\n", false, ok)
	}
	if value != nil {
		t.Errorf("Expected: %#v; Actual: %#v\n", nil, value)
	}
}

func TestReturnsNothingAfterExpiration(t *testing.T) {
	cache := NewTTL()
	defer cache.Quit()
	cache.Set("key", "value", 1000)
	time.Sleep(1500)
	value, ok := cache.Get("key")
	if ok != false {
		t.Errorf("Expected: %#v; Actual: %#v\n", false, ok)
	}
	if value != nil {
		t.Errorf("Expected: %#v; Actual: %#v\n", nil, value)
	}
}

func TestPrunesExpiredEntryOnGet(t *testing.T) {
	cache := NewTTL()
	defer cache.Quit()
	cache.Set("key1", "value", -5)
	cache.Set("key2", "value", -5)
	value, ok := cache.Get("key1")
	if ok != false {
		t.Errorf("Expected: %#v; Actual: %#v\n", false, ok)
	}
	if value != nil {
		t.Errorf("Expected: %#v; Actual: %#v\n", nil, value)
	}
	if len(cache.db) != 1 {
		t.Errorf("Expected: %#v; Actual: %#v\n", 1, len(cache.db))
	}
}

func TestPruneRemovesExpiredEntries(t *testing.T) {
	cache := NewTTL()
	defer cache.Quit()
	cache.Set("key1", "should be gone", -5)
	cache.Set("key2", "should also be gone", -5)
	cache.Set("key3", "should be there", 5*time.Second)
	cache.Set("key4", "should also be there", 5*time.Second)
	cache.Set("key5", "should also be there, also", 5*time.Second)
	cache.Prune()
	cache.Get("force following lines to be done after Prune completes by serializing with a call to Get")
	if len(cache.db) != 3 {
		t.Errorf("Expected: %#v; Actual: %#v\n", 3, len(cache.db))
	}
}
