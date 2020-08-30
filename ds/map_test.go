package ds

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestNewMap(t *testing.T) {
	lm := NewLiveMap(map[string]interface{}{})

	if lm == nil {
		t.Error("should not be a nil")
	}
}

func TestRead(t *testing.T) {
	lm := &LiveMap{}

	bytes, _ := json.Marshal(map[string]interface{}{
		"hello": "world",
		"age":   29,
	})

	lm.Read(bytes)

	if assertNumKeys(lm, 2) == false {
		t.Errorf("should have %d fields, instead %v\n", 2, numKeys(lm))
	}
}

func TestGetInt(t *testing.T) {
	lm := &LiveMap{}

	bytes, _ := json.Marshal(map[string]interface{}{
		"hello": "world",
		"age":   29,
	})

	lm.Read(bytes)

	val, ok := lm.GetInt("age")
	if !ok {
		t.Errorf("should have key: %s, but cannot find it.\n", "age")
		return
	}

	if val != 29 {
		t.Errorf("should have: %v, but get %v\n", 29, val)
		return
	}
}

func TestGetIntFromString(t *testing.T) {
	lm := getExample()

	val, ok := lm.GetInt("year")
	if !ok {
		t.Errorf("should have key: %s, but cannot find it.\n", "age")
		return
	}

	if val != 2020 {
		t.Errorf("should have: %v, but get %v\n", 29, val)
		return
	}
}

func TestGetString(t *testing.T) {
	lm := getExample()

	str, ok := lm.GetString("hello")
	if !ok {
		t.Errorf("should have key %s, but does not.\n", "hello")
		return
	}

	if str != "world" {
		t.Errorf("should have val: %s, but instead: %s\n", "world", str)
		return
	}
}

func TestGetStringFromInt(t *testing.T) {
	lm := getExample()

	str, ok := lm.GetString("age")
	if !ok {
		t.Errorf("should have key %s, but does not.\n", "age")
		return
	}

	if str != "29" {
		t.Errorf("should have val: %s, but instead: %s\n", "29", str)
		return
	}
}

func TestGetStringFromFloat(t *testing.T) {
	lm := getExample()

	str, ok := lm.GetString("weight")
	if !ok {
		t.Errorf("should have key %s, but does not.\n", "weight")
		return
	}

	if str != "60.3" {
		t.Errorf("should have val: %s, but instead: %s\n", "60.3", str)
		return
	}
}

func TestGetFloat(t *testing.T) {
	lm := getExample()

	f, ok := lm.GetFloat("weight")
	if !ok {
		t.Errorf("should have key %s, but does not.\n", "weight")
		return
	}

	if f != 60.3 {
		t.Errorf("should have val: %v, but instead: %v\n", 60.3, f)
		return
	}
}

func TestGetBool(t *testing.T) {
	lm := getExample()

	b, ok := lm.GetBool("good")
	if !ok {
		t.Errorf("should have key %s, but does not.\n", "good")
		return
	}

	if b != true {
		t.Errorf("should have val: %v, but instead: %v\n", true, b)
		return
	}
}

func TestGetNestString(t *testing.T) {
	lm := NewLiveMap(map[string]interface{}{
		"user": map[string]interface{}{
			"name": "mirkowang",
		},
	})

	name, ok := lm.GetString("user.name")
	if !ok {
		t.Errorf("should have key %s, but does not.\n", "user.name")
		return
	}

	if name != "mirkowang" {
		t.Errorf("should have val: %v, but instead: %v\n", "mirko", name)
		return
	}

	_, ok = lm.GetString("user.something")
	if ok {
		t.Errorf("should have no key %s, but get false positive.\n", "user.somthing")
		return
	}
}

func TestGetLiveMap(t *testing.T) {
	lm := NewLiveMap(map[string]interface{}{
		"user": map[string]interface{}{
			"name": "mirkowang",
		},
	})

	u, ok := lm.GetLiveMap("user")
	if !ok {
		t.Errorf("should have key %s, but does not.\n", "user")
		return
	}

	if reflect.TypeOf(u) != reflect.TypeOf(lm) {
		t.Errorf("should have livemap: %v, but instead: %v\n", reflect.TypeOf(lm), reflect.TypeOf(u))
		return
	}
}

func getExample() *LiveMap {
	lm := &LiveMap{}

	bytes, _ := json.Marshal(map[string]interface{}{
		"hello":  "world",
		"age":    29,
		"year":   "2020",
		"weight": 60.3,
		"good":   true,
	})

	lm.Read(bytes)

	return lm
}

func assertNumKeys(lm *LiveMap, n int) bool {
	lv := reflect.ValueOf(lm)
	if lv.Kind() == reflect.Ptr {
		lv = lv.Elem()
	}

	return len(lv.MapKeys()) == n
}

func numKeys(lm *LiveMap) int {
	lv := reflect.ValueOf(lm)
	if lv.Kind() == reflect.Ptr {
		lv = lv.Elem()
	}

	return len(lv.MapKeys())
}
