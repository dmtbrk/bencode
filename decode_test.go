package bencode

import (
	"errors"
	"reflect"
	"testing"
)

func TestUnmarshalErrInvalidArgument(t *testing.T) {
	input := []byte(`5:valid`)
	var output interface{}
	var want *ErrInvalidArgument

	got := Unmarshal(input, output)
	if !errors.Is(got, want) {
		t.Error("got:", got, "want:", want)
	}
}

func TestUnmarshalIntoInterface(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  interface{}
	}{
		{"Int", `i42e`, int64(42)},
		{"String", `4:spam`, "spam"},
		{"List", `l4:spame`, []interface{}{"spam"}},
		{"Dict", `d4:spam4:eggse`, map[string]interface{}{"spam": "eggs"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got interface{}

			data := []byte(test.input)
			err := Unmarshal(data, &got)

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

func TestUnmarshalIntoInt(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		input := `i42e`
		var want int64 = 42
		var got int64

		data := []byte(input)
		err := Unmarshal(data, &got)

		if err != nil {
			t.Error("unexpected error:", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("\ngot: %v \nwant: %v", got, want)
		}
	})
}

func TestUnmarshalIntoString(t *testing.T) {
	t.Run("Simple string", func(t *testing.T) {
		input := `4:spam`
		want := "spam"
		var got string

		data := []byte(input)
		err := Unmarshal(data, &got)

		if err != nil {
			t.Error("unexpected error:", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("\ngot: %v \nwant: %v", got, want)
		}
	})
}

func TestUnmarshalIntoSlice(t *testing.T) {
	t.Run("Empty slice", func(t *testing.T) {
		input := `l4:spam4:eggse`
		want := []string{"spam", "eggs"}
		var got []string

		data := []byte(input)
		err := Unmarshal(data, &got)

		if err != nil {
			t.Error("unexpected error:", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("\ngot: %v \nwant: %v", got, want)
		}
	})
}

func TestUnmarshalIntoMap(t *testing.T) {
	t.Run("Key type error", func(t *testing.T) {
		input := `d1:ai1ee`

		var got map[int]int
		data := []byte(input)
		err := Unmarshal(data, &got)

		if err == nil {
			t.Error("FAILDE with no error, expected:", err)
		} else {
			t.Log("PASSED with expected error:", err)
		}
	})

	t.Run("Simple", func(t *testing.T) {
		input := `d4:spam4:eggse`
		want := map[string]string{"spam": "eggs"}

		var got map[string]string
		data := []byte(input)
		err := Unmarshal(data, &got)

		if err != nil {
			t.Error("unexpected error:", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("\ngot: %v \nwant: %v", got, want)
		}
	})
}

func TestUnmarshalIntoStruct(t *testing.T) {
	t.Run("Exported error", func(t *testing.T) {
		type TestStruct struct {
			a int64 `bencode:"a"`
		}
		input := `d1:ai1ee`

		got := &TestStruct{}
		data := []byte(input)
		err := Unmarshal(data, got)

		if err == nil {
			t.Error("no error, expected:", err)
		} else {
			t.Log("PASSED with expected error:", err)
		}
	})

	t.Run("Some complex case", func(t *testing.T) {
		type TestStruct struct {
			A int64             `bencode:"a"`
			B string            `bencode:"b"`
			C []string          `bencode:"c"`
			D map[string]string `bencode:"d"`
			E struct {
				A int64 `bencode:"a"`
			} `bencode:"e"`
		}
		input := `d1:ai1e1:b4:spam1:cl4:spam4:eggse1:dd4:spam4:eggse1:ed1:ai42eee`
		want := &TestStruct{
			A: 1,
			B: "spam",
			C: []string{"spam", "eggs"},
			D: map[string]string{"spam": "eggs"},
			E: struct {
				A int64 `bencode:"a"`
			}{42},
		}

		got := &TestStruct{}
		data := []byte(input)
		err := Unmarshal(data, got)

		if err != nil {
			t.Error("unexpected error:", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("\ngot: %v \nwant: %v", got, want)
		}
	})

	t.Run("Raw bencoded field", func(t *testing.T) {
		type TestStruct struct {
			Raw []byte `bencode:"raw"`
		}
		input := `d3:rawd5:first4:item6:second4:item5:order6:really7:matters1:!ee`
		want := &TestStruct{
			Raw: []byte(`d5:first4:item6:second4:item5:order6:really7:matters1:!e`),
		}

		got := &TestStruct{}
		data := []byte(input)
		err := Unmarshal(data, got)

		if err != nil {
			t.Error("unexpected error:", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("\ngot: %s \nwant: %s", got, want)
		}
	})
}
