package shardmap

import (
	"testing"
)

func TestMap_Clear(t *testing.T) {
	m := New(100)
	m.Set("a", 1)
	m.Clear()
	if m.Len() != 0 {
		t.Errorf("test clear, expect %d, get %d", 0, m.Len())
	}
}

func TestMap_Delete(t *testing.T) {
	m := New(100)
	m.Set("a", 1)
	m.Delete("a")
	if v, ok := m.Get("a"); ok {
		t.Errorf("test delete, expect %v, get %v", nil, v)
	}
}

func TestMap_DeleteAccept(t *testing.T) {
	m := New(100)
	m.Set("a", 1)
	m.DeleteAccept("a", func(prev interface{}, deleted bool) bool {
		return prev.(int) != 1
	})
	v, ok := m.Get("a")
	if !ok {
		t.Errorf("test delete accept, except %v, get %v", 1, v)
	}
}

func TestMap_Get(t *testing.T) {
	m := New(100)
	m.Set("a", 1)
	v, ok := m.Get("a")
	if !ok || v.(int) != 1{
		t.Errorf("test get, except %v, get %v", 1, v)
	}
}

func TestMap_Len(t *testing.T) {
	capacity := 100
	m := New(capacity)
	for capacity > 0 {
		m.Set("a", 1)
		capacity --
	}
	m.Set("b", 2)
	if m.Len() != 2 {
		t.Errorf("test len, except %v, get %v", 2, m.Len())
	}
}

func TestMap_Range(t *testing.T) {
	m := New(100)
	m.Set("a", 1)
	m.Range(func(key string, value interface{}) bool {
		if key != "a" || value.(int) != 1 {
			t.Errorf("test range, except %v, get %v", 1, value)
		}
		return true
	})
}

func TestMap_Set(t *testing.T) {
	m := New(100)
	m.Set("a", 1)
	prev, replaced := m.Set("a", 2)
	if prev.(int) != 1 || !replaced {
		t.Errorf("test set, except %v, get %v", 1, prev)
	}
}

func TestMap_SetAccept(t *testing.T) {
	m := New(100)
	m.SetAccept("a", 1, func(prev interface{}, replaced bool) bool {
		if !replaced {
			return false
		}
		return true
	})
	v, ok := m.Get("a")
	if ok {
		t.Errorf("test set accept, except %v, get %v", nil, v)
	}
}

