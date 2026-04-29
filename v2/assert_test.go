package assert

import (
	"testing"
)

func TestEqual(t *testing.T) {
	subTest := new(testing.T)
	assert := New(subTest)
	assert.Equal(1, 1)
	if subTest.Failed() {
		t.Errorf("expected integer equality failed")
	}
	assert.Equal("abc", "abc")
	if subTest.Failed() {
		t.Errorf("expected string equality failed")
	}
	assert.Equal("abc", "def")
	if !subTest.Failed() {
		t.Errorf("expected string inegality passed")
	}
	assert.Equal("abc", 1)
	if !subTest.Failed() {
		t.Errorf("expected inegality due to types failed %v != %v", "abc", 1)
	}
}

func TestNil(t *testing.T) {
	subt := new(testing.T)
	assert := New(subt)
	assert.Nil(nil)
	if subt.Failed() {
		t.Errorf("expected nil check failed")
	}

	assert.Nil("not nil")
	if !subt.Failed() {
		t.Errorf("expected failed nil check failed")
	}
}

func TestReportedLocation(t *testing.T) {
	a := New(t)
	a.True(false, "where is it reported")
}
