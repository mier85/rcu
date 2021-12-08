package rcu

import (
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	m := New[string, string]()
	m.Insert("foo", "bar")
	v, lookup := m.ReadMap()["foo"]
	if !lookup {
		t.Error("value foo not in map")
	}
	if v != "bar" {
		t.Error("unexpected value for foo")
	}
}

func TestDelete(t *testing.T) {
	m := New[string, string]()
	m.Insert("foo", "bar")
	v, lookup := m.ReadMap()["foo"]
	if !lookup {
		t.Error("value foo not in map")
	}
	if v != "bar" {
		t.Error("unexpected value for foo")
	}
	m.Delete("foo")

	_, lookup = m.ReadMap()["foo"]
	if lookup {
		t.Error("value foo not in map")
	}

}

func TestParallelWrite(t *testing.T) {
	amount := 10
	wg := &sync.WaitGroup{}
	wg.Add(amount)
	m := New[int, int]()
	for i := 0; i < amount; i++ {
		go func(i int) {
			defer wg.Done()
			m.Insert(i, i)
		}(i)
	}
	wg.Wait()
	res := m.ReadMap()
	if len(res) != amount {
		t.Errorf("unexpected amount of items expected: %d, got: %d ", amount, len(res))
	}
}

func TestParallelAccess(t *testing.T) {
	amount := 100
	wg := &sync.WaitGroup{}
	wg.Add(amount)
	m := New[int, int]()
	for i := 0; i < amount; i++ {
		go func(i int) {
			defer wg.Done()
			m.Insert(i, i)
		}(i)
	}
	for i := 0; i < amount; i++ {
		go func(i int) {
			// we are not really interested in the result, we just want to make sure we can access the read model simultaneously
			_, _ = m.ReadMap()[i]
		}(i)
	}
	wg.Wait()
	res := m.ReadMap()
	if len(res) != amount {
		t.Errorf("unexpected amount of items expected: %d, got: %d ", amount, len(res))
	}
}
