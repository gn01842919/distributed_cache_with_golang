package cache

import (
	"reflect"
	"testing"
)

func TestBasic(t *testing.T) {
	c := New("inmemory")

	stat := c.GetStat()
	if stat.Count != 0 || stat.KeySize != 0 || stat.ValueSize != 0 {
		t.Error("State is incorrect")
	}
	key1 := "key1"
	value1 := []byte{'v', 'a', 'l', '1'}

	key2 := "key2"
	value2 := []byte{'v', 'a', 'l', '2'}

	// Test get a key that does not exist
	val, e := c.Get(key1)
	if e != nil {
		t.Error("Fail to get key1")
	} else if len(val) != 0 {
		t.Error("incorrect value1")
	} else {
		t.Log("Test get non-existing key passed")
	}

	// Test Set
	e = c.Set(key1, value1)
	if e != nil {
		t.Error("Fail to set key1")
	} else {
		t.Log("Test set key1 passed")
	}

	stat = c.GetStat()
	if stat.Count != 1 || stat.KeySize != int64(len(key1)) || stat.ValueSize != int64(len(value1)) {
		t.Error("State is incorrect")
	}

	e = c.Set(key2, value2)
	if e != nil {
		t.Error("Fail to set key2")
	} else {
		t.Log("Test set key2 passed")
	}

	stat = c.GetStat()
	if stat.Count != 2 || stat.KeySize != int64(len(key1)+len(key2)) || stat.ValueSize != int64(len(value1)+len(value2)) {
		t.Error("State is incorrect")
	}

	// Test Get
	val, e = c.Get(key1)
	if e != nil {
		t.Error("Fail to get key1")
	} else if !reflect.DeepEqual(val, value1) {
		t.Error("incorrect value1")
	} else {
		t.Log("Test get non-existing key1 passed")
	}

	val, e = c.Get(key2)
	if e != nil {
		t.Error("Fail to get key2")
	} else if !reflect.DeepEqual(val, value2) {
		t.Error("incorrect value2")
	} else {
		t.Log("Test get non-existing key2 passed")
	}

	// Test overwrite
	e = c.Set(key1, value1)
	if e != nil {
		t.Error("Fail to set key1")
	} else {
		t.Log("Test set key1 passed")
	}

	stat = c.GetStat()
	if stat.Count != 2 || stat.KeySize != int64(len(key1)+len(key2)) || stat.ValueSize != int64(len(value1)+len(value2)) {
		t.Error("State is incorrect")
	}

	// Test Del
	e = c.Del(key1)
	if e != nil {
		t.Error("Fail to del key1")
	} else {
		t.Log("Test del key1 passed")
	}

	// Test get a deleted key1
	val, e = c.Get(key1)
	if e != nil {
		t.Error("Fail to get key1")
	} else if len(val) != 0 {
		t.Error("incorrect value1")
	} else {
		t.Log("Test get non-existing key passed")
	}

	stat = c.GetStat()
	if stat.Count != 1 || stat.KeySize != int64(len(key2)) || stat.ValueSize != int64(len(value2)) {
		t.Error("State is incorrect")
	}

}
