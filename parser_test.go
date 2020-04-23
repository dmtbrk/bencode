package bencode

import (
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Value
	}{
		// Int
		{"Int/Positive", `i42e`, Int(42)},
		{"Int/Negative", `i-42e`, Int(-42)},
		{"Int/Zero", `i0e`, Int(0)},
		// String
		{"String/Simple", `4:spam`, String("spam")},
		{"String/Empty", `0:`, String("")},
		// List
		{"List/One string", `l4:spame`, List{String("spam")}},
		{"List/Two strings", `l4:spam4:eggse`, List{String("spam"), String("eggs")}},
		{"List/String and int", `l4:spami42ee`, List{String("spam"), Int(42)}},
		{"List/Nested list", `l4:spaml4:spam4:eggsee`, List{String("spam"), List{String("spam"), String("eggs")}}},
		// Dict
		{"Dict/Simple", `d4:spam4:eggse`, NewDict([]DictItem{{String("spam"), String("eggs")}}...)},
		{"Dict/Three items", `d4:spam4:eggs3:key3:val6:answeri42ee`, NewDict([]DictItem{{String("spam"), String("eggs")}, {String("key"), String("val")}, {String("answer"), Int(42)}}...)},
		{"Dict/Nested dict", `d4:spamd4:spam4:eggsee`, NewDict([]DictItem{{String("spam"), NewDict([]DictItem{{String("spam"), String("eggs")}}...)}}...)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := NewParser(strings.NewReader(test.input))
			got, err := parser.Parse()
			if err != nil {
				t.Error("unexpected error:", err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("\ngot: %v \nwant: %v", got, test.want)
			}
			// t.Logf("\ngot: %v \nwant: %v", got, test.want)
		})
	}
}
