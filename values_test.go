package bencode

import (
	"reflect"
	"testing"
)

func TestNewDict(t *testing.T) {
	// empty Dict
	got := NewDict()
	want := &Dict{}

	if !reflect.DeepEqual(got, want) {
		t.Error("got:", got, "want:", want)
	}

	// construct from []DictItem
	got = NewDict([]DictItem{{String("spam"), String("eggs")}}...)
	want = &Dict{keys: []String{String("spam")}, m: map[String]Value{String("spam"): String("eggs")}}

	if !reflect.DeepEqual(got, want) {
		t.Error("got:", got, "want:", want)
	}
}

func TestDictGet(t *testing.T) {
	tests := []struct {
		name string
		dict *Dict
		key  String
		want Value
	}{
		{"Exist", &Dict{keys: []String{String("spam")}, m: map[String]Value{String("spam"): String("eggs")}}, String("spam"), String("eggs")},
		{"Not exist", &Dict{keys: []String{String("spam")}, m: map[String]Value{String("spam"): String("eggs")}}, String("notexist"), nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.dict.Get(test.key)

			if got != test.want {
				t.Error("got:", got, "want:", test.want)
			}
		})

	}
}

func TestBencode(t *testing.T) {
	tests := []struct {
		name string
		val  Value
		want string
	}{
		// Int
		{"Int/Positive", Int(42), `i42e`},
		{"Int/Negative", Int(-42), `i-42e`},
		{"Int/Zero", Int(0), `i0e`},
		// String
		{"String/Simple", String("spam"), `4:spam`},
		{"String/Empty", String(""), `0:`},
		// List
		{"List/One string", List{String("spam")}, `l4:spame`},
		{"List/Two strings", List{String("spam"), String("eggs")}, `l4:spam4:eggse`},
		{"List/String and int", List{String("spam"), Int(42)}, `l4:spami42ee`},
		{"List/Nested list", List{String("spam"), List{String("spam"), String("eggs")}}, `l4:spaml4:spam4:eggsee`},
		// Dict
		{"Dict/Simple", NewDict([]DictItem{{String("spam"), String("eggs")}}...), `d4:spam4:eggse`},
		{"Dict/Three items", NewDict([]DictItem{{String("spam"), String("eggs")}, {String("key"), String("val")}, {String("answer"), Int(42)}}...), `d4:spam4:eggs3:key3:val6:answeri42ee`},
		{"Dict/Nested dict", NewDict([]DictItem{{String("spam"), NewDict([]DictItem{{String("spam"), String("eggs")}}...)}}...), `d4:spamd4:spam4:eggsee`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.val.Bencode()
			sgot := string(got)
			if sgot != test.want {
				t.Error("got:", sgot, "want:", test.want)
			}
			// t.Log("got:", string(got), "want:", test.want)
		})
	}
}

// func TestIntBencode(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		in   int64
// 		want string
// 	}{
// 		{"Positive", 42, `i42e`},
// 		{"Negative", -42, `i-42e`},
// 		{"Zero", 0, `i0e`},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			intVal := Int(test.in)
// 			got := intVal.Bencode()
// 			sgot := string(got)
// 			if sgot != test.want {
// 				t.Error("got:", sgot, "want:", test.want)
// 			}
// 			// t.Log("got:", sgot, "want:", test.want)
// 		})
// 	}
// }

// func TestStringBencode(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		in   string
// 		want string
// 	}{
// 		{"Simple", "spam", `4:spam`},
// 		{"Empty", "", `0:`},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			val := String(test.in)
// 			got, err := val.Bencode()
// 			if err != nil {
// 				t.Error("unexpected error:", err)
// 			}
// 			sgot := string(got)
// 			if sgot != test.want {
// 				t.Error("got:", sgot, "want:", test.want)
// 			}
// 			// t.Log("got:", string(got), "want:", test.want)
// 		})
// 	}
// }

// func TestListBencode(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		in   []Value
// 		want string
// 	}{
// 		{"1 string", []Value{String("spam")}, `l4:spame`},
// 		{"2 strings", []Value{String("spam"), String("eggs")}, `l4:spam4:eggse`},
// 		{"String and int", []Value{String("spam"), Int(42)}, `l4:spami42ee`},
// 		{"Nested list", []Value{String("spam"), List([]Value{String("spam"), String("eggs")})}, `l4:spaml4:spam4:eggsee`},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			val := List(test.in)
// 			got, err := val.Bencode()
// 			if err != nil {
// 				t.Error("unexpected error:", err)
// 			}
// 			sgot := string(got)
// 			if sgot != test.want {
// 				t.Error("got:", sgot, "want:", test.want)
// 			}
// 			// t.Log("got:", string(got), "want:", test.want)
// 		})
// 	}
// }

// func TestDictBencode(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		val  Dict
// 		want string
// 	}{
// 		{"Simple", NewDict([]DictItem{{String("spam"), String("eggs")}}...), `d4:spam4:eggse`},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			got, err := test.val.Bencode()
// 			if err != nil {
// 				t.Error("unexpected error:", err)
// 			}
// 			sgot := string(got)
// 			if sgot != test.want {
// 				t.Error("got:", sgot, "want:", test.want)
// 			}
// 			// t.Log("got:", string(got), "want:", test.want)
// 		})
// 	}
// }
