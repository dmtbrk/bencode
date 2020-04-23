package bencode

import "strconv"

// Value is a tree node.
type Value interface {
	Interface() interface{}
	Bencode() []byte
}

// Int is a representation of bencoded integer.
type Int int64

// Interface returns an int64 put into interface{}.
func (i Int) Interface() interface{} {
	return int64(i)
}

// Bencode returns a bencoded integer.
func (i Int) Bencode() []byte {
	b := []byte{'i'}
	b = strconv.AppendInt(b, int64(i), 10)
	b = append(b, 'e')
	return b
}

// String is a representation of bencoded string.
type String string

// Interface returns a string put into interface{}.
func (s String) Interface() interface{} {
	return string(s)
}

// Bencode returns a bencoded string.
func (s String) Bencode() []byte {
	var b []byte
	b = strconv.AppendInt(b, int64(len(s)), 10)
	b = append(b, ':')
	b = append(b, s...)
	return b
}

// List is a representation of bencoded list as a slice of Value.
type List []Value

// Interface returns a []interface{} put into interface{}.
func (l List) Interface() interface{} {
	s := make([]interface{}, len(l))
	for i, v := range l {
		s[i] = v.Interface()
	}
	return s
}

// Bencode returns a bencoded list.
func (l List) Bencode() []byte {
	b := []byte{'l'}
	for _, v := range l {
		bval := v.Bencode()
		b = append(b, bval...)
	}
	b = append(b, 'e')
	return b
}

// Dict is a representation of bencoded dictionary. The need for preserving
// the order of keys has made this type a struct, not just a map[String]Value,
// so it must be used as a pointer type.
// Preserving the order of keys is needed to ensure that the order of bencoded items
// is the same as defined on creation so that it has a predictable output and
// decode/encode roundtrip produces the same result.
type Dict struct {
	keys []String
	m    map[String]Value
}

// DictItem is a helper struct which represents a dict key-value pair.
type DictItem struct {
	Key   String
	Value Value
}

// NewDict constructs a dict from key-value pairs DictItem.
func NewDict(items ...DictItem) *Dict {
	d := &Dict{}
	for _, item := range items {
		d.Set(item.Key, item.Value)
	}
	return d
}

// Get returns a value stored in Dict by provided key.
func (d *Dict) Get(key String) Value {
	return d.m[key]
}

// Set puts a key-value pair into Dict.
func (d *Dict) Set(key String, val Value) {
	if d.m == nil {
		d.m = make(map[String]Value)
	}
	_, ok := d.m[key]
	if ok { // find and delete key if present
		for i, k := range d.keys {
			if k == key {
				d.keys = append(d.keys[:i], d.keys[i+1:]...)
				break
			}
		}
	}
	d.keys = append(d.keys, key)
	d.m[key] = val
}

// Interface returns a map[string]interface{} representation of Dict put into interface{}.
func (d *Dict) Interface() interface{} {
	m := make(map[string]interface{}, len(d.m))
	for k, v := range d.m {
		m[string(k)] = v.Interface()
	}
	return m
}

// Bencode returns a bencoded dictionary. The order of key-value pairs
// would be the same as with which Dict was created or updated.
func (d *Dict) Bencode() []byte {
	b := []byte{'d'}
	for _, key := range d.keys {
		val := d.m[key]
		bkey := key.Bencode()
		bval := val.Bencode()
		b = append(b, bkey...)
		b = append(b, bval...)
	}
	b = append(b, 'e')
	return b
}
