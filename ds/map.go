package ds

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// LiveMap live map
type LiveMap map[string]interface{}

// NewLiveMap get a new live map
func NewLiveMap(d map[string]interface{}) *LiveMap {
	lm := LiveMap(d)

	return &lm
}

// MapFrom todo
func MapFrom(str []byte) (*LiveMap, error) {
	lm := &LiveMap{}

	if err := lm.Read(str); err != nil {
		return nil, err
	}

	return lm, nil
}

// Read read a []byte
func (lm *LiveMap) Read(str []byte) error {
	return json.Unmarshal(str, lm)
}

// GetInt get integer
func (lm *LiveMap) GetInt(key string) (int, bool) {
	val, ok := (*lm)[key]
	if !ok {
		return 0, false
	}

	var v int
	if err := decodeInt(val, &v); err != nil {
		return 0, true
	}

	return v, true
}

// GetString get a string
func (lm *LiveMap) GetString(key string) (string, bool) {

	val, ok := getVal(lm, strings.Split(key, "."))
	if !ok {
		return "", false
	}

	var s string
	if err := decodeString(val, &s); err != nil {
		return "", false
	}

	return s, true
}

// GetLiveMap todo
func (lm *LiveMap) GetLiveMap(key string) (*LiveMap, bool) {
	val, ok := getVal(lm, strings.Split(key, "."))
	if !ok {
		return nil, false
	}

	v, ok := val.(map[string]interface{})
	if !ok {
		return nil, false
	}

	return NewLiveMap(v), true
}

func getVal(lm *LiveMap, keys []string) (interface{}, bool) {

	for i, str := range keys {
		val, ok := (*lm)[str]
		if !ok {
			return nil, false
		}

		if i == len(keys)-1 {
			return val, true
		}

		if isBasicType(val) {
			return nil, false
		}

		m, ok := val.(map[string]interface{})
		if !ok {
			return nil, false
		}

		return getVal(NewLiveMap(m), keys[1:])
	}

	return nil, true
}

func isBasicType(val interface{}) bool {
	return reflect.TypeOf(val).String() == reflect.ValueOf(val).Kind().String()
}

// GetFloat todo
func (lm *LiveMap) GetFloat(key string) (float64, bool) {
	val, ok := (*lm)[key]
	if !ok {
		return 0, false
	}

	var number float64
	if err := decodeFloat(val, &number); err != nil {
		return 0, false
	}

	return number, true
}

// GetBool todo
func (lm *LiveMap) GetBool(key string) (bool, bool) {
	val, ok := (*lm)[key]
	if !ok {
		return false, false
	}

	var b bool
	if err := decodeBool(val, &b); err != nil {
		return false, false
	}

	return b, true
}

func decodeInt(val interface{}, number *int) error {
	rv := reflect.ValueOf(val)
	n := reflect.ValueOf(number).Elem()

	switch getKind(rv) {
	case reflect.Int:
		n.SetInt(rv.Int())

	case reflect.Uint:
		n.SetInt(int64(rv.Uint()))

	case reflect.Float32:
		n.SetInt(int64(rv.Float()))

	case reflect.Bool:
		if rv.Bool() {
			n.SetInt(1)
			return nil
		}
		n.SetInt(0)

	case reflect.String:
		i, _ := strconv.ParseInt(rv.String(), 0, n.Type().Bits())
		n.SetInt(i)

	default:
		return fmt.Errorf("expected type '%s', got unconvertible type '%s'", n.Type(), rv.Type())
	}

	return nil
}

func decodeString(v interface{}, str *string) error {
	rv := reflect.ValueOf(v)
	vs := reflect.ValueOf(str).Elem()

	switch getKind(rv) {
	case reflect.Int:
		vs.SetString(strconv.Itoa(int(rv.Int())))

	case reflect.Uint:
		vs.SetString(strconv.Itoa(int(rv.Uint())))

	case reflect.Float32:
		vs.SetString(strconv.FormatFloat(rv.Float(), 'g', -1, 64))

	case reflect.Bool:
		if rv.Bool() {
			vs.SetString("true")
			return nil
		}
		vs.SetString("false")

	case reflect.String:
		vs.SetString(rv.String())

	default:
		return fmt.Errorf("expected type '%s', got unconvertible type '%s'", vs.Type(), rv.Type())
	}

	return nil
}

func decodeFloat(v interface{}, number *float64) error {
	rv := reflect.ValueOf(v)
	rn := reflect.ValueOf(number).Elem()

	switch getKind(rv) {
	case reflect.Int:
		rn.SetFloat(float64(rv.Int()))

	case reflect.Float32:
		rn.SetFloat(rv.Float())

	case reflect.String:
		f, err := strconv.ParseFloat(rv.String(), 64)
		if err != nil {
			return err
		}
		rn.SetFloat(f)

	case reflect.Bool:
		if rv.Bool() {
			rn.SetFloat(1)
			return nil
		}
		rn.SetFloat(0)
	default:
		return fmt.Errorf("expected type '%s', got unconvertible type '%s'", rn.Type(), rv.Type())
	}

	return nil
}

func decodeBool(v interface{}, b *bool) error {
	rv := reflect.ValueOf(v)
	rb := reflect.ValueOf(b).Elem()

	switch getKind(rv) {
	case reflect.Int:
		rb.SetBool(rv.Int() > 0)

	case reflect.Float32:
		rb.SetBool(rv.Float() > 0)

	case reflect.Bool:
		rb.SetBool(rv.Bool())

	case reflect.String:
		sv, err := strconv.ParseBool(rv.String())
		if err != nil {
			return err
		}
		rb.SetBool(sv)

	default:
		return fmt.Errorf("expected type '%s', got unconvertible type '%s'", rb.Type(), rv.Type())
	}

	return nil
}

func getKind(val reflect.Value) reflect.Kind {
	kind := val.Kind()

	switch {
	case kind >= reflect.Int && kind <= reflect.Int64:
		return reflect.Int

	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return reflect.Uint

	case kind >= reflect.Float32 && kind <= reflect.Float64:
		return reflect.Float32

	default:
		return kind
	}
}
