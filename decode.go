package bencode

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

// ErrInvalidArgument describes an error which occurs when an invalid
// argument type is passed to the Unmarshal function.
type ErrInvalidArgument struct {
	t reflect.Type
}

func (e *ErrInvalidArgument) Error() string {
	return fmt.Sprintf("bencode: invalid argument (%s), want a pointer", e.t)
}

// Is satisfies errors.Is requirements.
func (e *ErrInvalidArgument) Is(err error) bool {
	_, ok := err.(*ErrInvalidArgument)
	return ok
}

// Unmarshal parses the bencoded data and stores the result in the value
// pointed by v. If v is nil or not a pointer, Unmarshal returns an
// ErrInvalidArgument.
func Unmarshal(data []byte, i interface{}) error {
	r := bytes.NewReader(data)
	d := NewDecoder(r)
	return d.Decode(i)
}

type Decoder struct {
	reader io.Reader
	parser *Parser
}

func NewDecoder(r io.Reader) *Decoder {
	p := NewParser(r)
	return &Decoder{reader: r, parser: p}
}

// Decode takes reader
func (d *Decoder) Decode(i interface{}) error {
	// i must be a pointer
	p := reflect.ValueOf(i)
	if p.Kind() != reflect.Ptr || p.IsNil() {
		return &ErrInvalidArgument{reflect.TypeOf(i)}
	}

	// parse data
	v, err := d.parser.Parse()
	if err != nil {
		return err
	}

	rv := p.Elem() // get what p points to
	return d.put(rv, v)
}

// put dispatches putting data from Value to reflect.Value depending on
// the Kind of the destination
func (d *Decoder) put(dst reflect.Value, src Value) error {
	if dst.Kind() == reflect.Interface && dst.NumMethod() == 0 {
		dst.Set(reflect.ValueOf(src.Interface()))
	}

	switch dst.Kind() {
	case reflect.Int64:
		return d.putInt(dst, src)
	case reflect.String:
		return d.putString(dst, src)
	case reflect.Slice:
		if dst.Type().Elem().Kind() == reflect.Uint8 {
			return d.putBencode(dst, src)
		}
		return d.putSlice(dst, src)
	case reflect.Map:
		return d.putMap(dst, src)
	case reflect.Struct:
		return d.putStruct(dst, src)
	}

	return nil
}

func (d *Decoder) putInt(dst reflect.Value, src Value) error {
	i, ok := src.(Int)
	if !ok {
		return fmt.Errorf("trying to put %T into int", src)
	}

	dst.SetInt(int64(i))

	return nil
}

func (d *Decoder) putString(dst reflect.Value, src Value) error {
	s, ok := src.(String)
	if !ok {
		return fmt.Errorf("trying to put %T into string", src)
	}

	dst.SetString(string(s))

	return nil
}

func (d *Decoder) putSlice(dst reflect.Value, src Value) error {
	l, ok := src.(List)
	if !ok {
		return fmt.Errorf("trying to put %T into slice", src)
	}

	// extend destination slice, if needed
	// also handles allocating new slice if dst is nil slice
	dif := len(l) - dst.Len()
	if dif > 0 {
		ext := reflect.MakeSlice(dst.Type(), dif, dif)
		dst.Set(reflect.AppendSlice(dst, ext))
	}

	for i, v := range l {
		elem := dst.Index(i)
		d.put(elem, v)
	}

	return nil
}

func (d *Decoder) putMap(dst reflect.Value, src Value) error {
	// only strings allowed to be the keys in bencode
	mapKeyType := dst.Type().Key()
	if mapKeyType.Kind() != reflect.String {
		return fmt.Errorf("map keys must be of type string, not %v", mapKeyType)
	}

	dict, ok := src.(*Dict)
	if !ok {
		return fmt.Errorf("trying to put %T into map", src)
	}

	// handles allocating a new map if dst is a nil map (zero value)
	if dst.IsZero() {
		dst.Set(reflect.MakeMap(dst.Type()))
	}

	mapElemType := dst.Type().Elem()
	for k, v := range dict.m {
		key := reflect.ValueOf(string(k))
		elem := reflect.New(mapElemType).Elem()
		d.put(elem, v)
		dst.SetMapIndex(key, elem)
	}

	return nil
}

func (d *Decoder) putStruct(dst reflect.Value, src Value) error {
	dict, ok := src.(*Dict)
	if !ok {
		return fmt.Errorf("trying to put %T into struct", src)
	}

	for i := 0; i < dst.NumField(); i++ {
		field := dst.Field(i)
		if !field.CanSet() {
			return fmt.Errorf("struct field must be settable, i.e. exported")
		}
		tag := dst.Type().Field(i).Tag.Get("bencode")
		if value := dict.Get(String(tag)); value != nil {
			if err := d.put(field, value); err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *Decoder) putBencode(dst reflect.Value, src Value) error {
	dst.Set(reflect.ValueOf(src.Bencode()))
	return nil
}
