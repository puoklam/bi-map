package bimap

import "testing"

func TestGetSetFront(t *testing.T) {
	m := New[string, string]()
	k, v := "k", "v"
	m.SetFront(k, v)
	val, ok := m.GetFront(k)
	if !ok {
		t.Error("Key not exists")
	}
	if v != val {
		t.Errorf("Values not equal, want: %s, got: %s", v, val)
	}
}

func TestGetSetBack(t *testing.T) {
	m := New[string, string]()
	k, v := "k", "v"
	m.SetBack(k, v)
	val, ok := m.GetBack(k)
	if !ok {
		t.Error("Key not exists")
	}
	if v != val {
		t.Errorf("Values not equal, want: %s, got: %s", v, val)
	}
}

func TestDeleteFront(t *testing.T) {
	m := New[string, string]()
	k, v := "k", "v"
	m.SetFront(k, v)
	m.DeleteFront(k)
	if _, ok := m.GetFront(k); ok {
		t.Error("Should be deleted")
	}
}

func TestDeleteBack(t *testing.T) {
	m := New[string, string]()
	k, v := "k", "v"
	m.SetBack(k, v)
	m.DeleteBack(k)
	if _, ok := m.GetFront(k); ok {
		t.Error("Should be deleted")
	}
}
