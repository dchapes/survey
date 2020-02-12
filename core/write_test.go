package core

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestWriteAnswer(t *testing.T) {
	ifMap := make(map[string]interface{})
	isMap := make(map[int]string)
	var testSet1, testSet2 testFieldSettable
	var testString1 testStringSettable
	var testSetStruct testTaggedStruct
	testPtrSetStruct := testPtrTaggedStruct{&testStringSettable{}}
	type structT struct {
		Name    string
		Age     uint
		Male    bool
		Height  float64
		Timeout time.Duration
	}
	_, parseBoolErr := strconv.ParseBool("hello")

	tests := []struct {
		name       string
		field      string
		t, v, want interface{}
		err        error
	}{
		{
			name:  "returnsErrorIfTargetNotPtr",
			t:     true,
			field: "hello",
			v:     true,
			err:   errNeedPointer,
		},
		{
			name: "canWriteToBool",
			t:    new(bool),
			v:    true,
		},
		{
			name: "canWriteString",
			t:    new(string),
			v:    "hello",
		},
		{
			name: "canWriteSlice",
			t:    new([]string),
			v:    []string{"hello", "world"},
		},
		{
			name: "recoversInvalidReflection",
			t:    new(bool),
			v:    "hello",
			err:  parseBoolErr,
		},
		{
			name: "handlesNonStructValues",
			t:    new(string),
			v:    "world",
		},
		{
			name:  "canMutateStruct",
			t:     &struct{ Name string }{},
			field: "name",
			v:     "world",
			want:  struct{ Name string }{"world"},
		},
		{
			name: "optionAnswer/writes index for ints",
			t:    new(int),
			v:    OptionAnswer{Index: 10, Value: "string value"},
			want: 10,
		},
		{
			name: "optionAnswer/writes value for strings",
			t:    new(string),
			v:    OptionAnswer{Index: 10, Value: "string value"},
			want: "string value",
		},
		{
			name: "optionAnswer/writes OptionAnswer for OptionAnswer",
			t:    new(OptionAnswer),
			v:    OptionAnswer{Index: 10, Value: "string value"},
		},
		{
			name: "optionAnswer/writes slice of indices for slice of ints",
			t:    new([]int),
			v:    []OptionAnswer{{Index: 10, Value: "string value"}},
			want: []int{10},
		},
		{
			name: "optionAnswer/writes slice of values for slice of strings",
			t:    new([]string),
			v:    []OptionAnswer{{Index: 10, Value: "string value"}},
			want: []string{"string value"},
		},
		{
			name:  "canMutateMap",
			t:     &ifMap,
			field: "name",
			v:     "world",
			want:  map[string]interface{}{"name": "world"},
		},
		{
			name:  "returnsErrorIfInvalidMapType",
			t:     &isMap,
			field: "name",
			v:     "world",
			err:   errMapType,
		},
		{
			name:  "writesStringSliceToIntSlice",
			t:     new([]int),
			field: "name",
			v:     []string{"1", "2", "3"},
			want:  []int{1, 2, 3},
		},
		{
			name:  "writesStringArrayToIntArray",
			t:     new([3]int),
			field: "name",
			v:     [3]string{"1", "2", "3"},
			want:  [3]int{1, 2, 3},
		},
		{
			name:  "returnsErrWhenFieldNotFound",
			t:     &struct{ Name string }{},
			field: "",
			v:     "world",
			err:   fieldNotMatchError{""},
		},
		// CONVERSION TESTS
		{
			name: "canStringToBool",
			t:    new(bool),
			v:    "true",
			want: true,
		},
		{
			name: "canStringToInt",
			t:    new(int),
			v:    "2",
			want: 2,
		},
		{
			name: "canStringToInt8",
			t:    new(int8),
			v:    "3",
			want: int8(3),
		},
		{
			name: "canStringToInt16",
			t:    new(int16),
			v:    "4",
			want: int16(4),
		},
		{
			name: "canStringToInt32",
			t:    new(int32),
			v:    "5",
			want: int32(5),
		},
		{
			name: "canStringToInt64",
			t:    new(int64),
			v:    "6",
			want: int64(6),
		},
		{
			name: "canStringToUint",
			t:    new(uint),
			v:    "7",
			want: uint(7),
		},
		{
			name: "canStringToUint8",
			t:    new(uint8),
			v:    "8",
			want: uint8(8),
		},
		{
			name: "canStringToUint16",
			t:    new(uint16),
			v:    "9",
			want: uint16(9),
		},
		{
			name: "canStringToUint32",
			t:    new(uint32),
			v:    "10",
			want: uint32(10),
		},
		{
			name: "canStringToUint64",
			t:    new(uint64),
			v:    "11",
			want: uint64(11),
		},
		{
			name: "canStringToFloat32",
			t:    new(float32),
			v:    "2.32",
			want: float32(2.32),
		},
		{
			name: "canStringToFloat64",
			t:    new(float64),
			v:    "2.64",
			want: float64(2.64),
		},
		{
			name:  "canConvertStructFieldType/string",
			t:     new(structT),
			field: "name",
			v:     "Bob",
			want:  structT{Name: "Bob"},
		},
		{
			name:  "canConvertStructFieldType/uint",
			t:     new(structT),
			field: "age",
			v:     "22",
			want:  structT{Age: 22},
		},
		{
			name:  "canConvertStructFieldType/bool",
			t:     new(structT),
			field: "male",
			v:     "true",
			want:  structT{Male: true},
		},
		{
			name:  "canConvertStructFieldType/float64",
			t:     new(structT),
			field: "height",
			v:     "6.2",
			want:  structT{Height: 6.2},
		},
		{
			name:  "canConvertStructFieldType/Time",
			t:     new(structT),
			field: "timeout",
			v:     "30s",
			want:  structT{Timeout: 30 * time.Second},
		},
		// WithFieldSettable
		{
			name:  "WithFieldSettable/valueMap",
			t:     &testSet1,
			field: "values",
			v:     "stringVal",
			want:  testFieldSettable{Values: map[string]string{"values": "stringVal"}},
		},
		{
			name:  "WithFieldSettable/error",
			t:     &testSet2,
			field: "values",
			v:     int64(123),
			err:   fmt.Errorf("incompatible type int64"),
		},
		{
			name:  "WithFieldSettable/StringSettable",
			t:     &testString1,
			field: "value1",
			v:     "testString1",
			want:  testStringSettable{"testString1"},
		},
		{
			name:  "WithFieldSettable/TaggedStruct",
			t:     &testSetStruct,
			field: "tagged",
			v:     "stringVal1",
			want:  testTaggedStruct{TaggedValue: testStringSettable{"stringVal1"}},
		},
		{
			name:  "WithFieldSettable/PtrTaggedStruct",
			t:     &testPtrSetStruct,
			field: "tagged",
			v:     "stringVal1",
			want:  testPtrTaggedStruct{TaggedValue: &testStringSettable{"stringVal1"}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			//t.Logf("calling WriteAnswer(%#v, %q, %v)", tc.t, tc.field, tc.v)
			// Capture the argument before WriteAnswer mucks with it
			tstr := fmt.Sprintf("%#v", tc.t)
			err := WriteAnswer(tc.t, tc.field, tc.v)
			if err != nil {
				if tc.err == nil {
					t.Fatalf("WriteAnswer(%v, %q, %v) unexpectedly failed:\n\t%v",
						tstr, tc.field, tc.v, err,
					)
				} else if !reflect.DeepEqual(err, tc.err) {
					t.Fatalf("WriteAnswer(%v, %q, %v)\n\tgave err: %#v\n\twant err: %#v",
						tstr, tc.field, tc.v, err, tc.err,
					)
					//} else {
					//	t.Logf("WriteAnswer(%v, %q, %v) failed as expected", tstr, tc.field, tc.v)
				}
				return
			} else if tc.err != nil {
				t.Fatalf("WriteAnswer(%v, %q, %v) unexpectedly succeeded\n\texpected: %v",
					tstr, tc.field, tc.v, tc.err,
				)
			}
			out := reflect.Indirect(reflect.ValueOf(tc.t)).Interface() // ick
			want := tc.want
			if want == nil {
				want = tc.v
			}
			if !reflect.DeepEqual(out, want) {
				t.Fatalf("WriteAnswer(%v, %q, %v)\n\tset  (%T)(%[4]v)\n\twant (%T)(%[5]v)",
					tstr, tc.field, tc.v, out, want,
				)
				//} else {
				//	t.Logf("correct:\n\tset  (%T)(%[1]v)\n\twant (%T)(%[2]v)", out, want)
			}
		})
	}
}

func TestFindFieldIndex(t *testing.T) {
	var value struct {
		Name     string
		Username string `survey:"tagged"`
		Other    string
		Replace  string `survey:"other"`
	}
	rv := reflect.ValueOf(value)
	_ = rv

	tests := []struct {
		field string
		i     int
		err   error
	}{
		{"name", 0, nil},   // canFindExportedField
		{"tagged", 1, nil}, // canFindTaggedField
		{"Name", 0, nil},   // canHandleCapitalAnswerNames
		{"other", 3, nil},  // tagOverwriteFieldName
		{"nosuchfield", 0, fieldNotMatchError{"nosuchfield"}},
		{"", 0, fieldNotMatchError{""}},
	}
	for _, tc := range tests {
		//t.Logf("calling findFieldIndex(%T, %q)", rv, tc.field)
		i, err := findFieldIndex(rv, tc.field)
		if err != nil {
			if tc.err == nil {
				t.Errorf("findFieldIndex(%T, %q) unexpectedly failed:\n\t%v",
					rv, tc.field, err,
				)
			} else if !reflect.DeepEqual(err, tc.err) {
				t.Errorf("findFieldIndex(%T, %q)\n\tgave err: %#v\n\twant err: %#v",
					rv, tc.field, err, tc.err,
				)
			}
			continue
		} else if tc.err != nil {
			t.Errorf("findFieldIndex(%T, %q) unexpectedly succeeded\n\texpected: %v",
				rv, tc.field, tc.err,
			)
			continue
		}
		if i != tc.i {
			t.Errorf("findFieldIndex(%T, %q) gave %d, want %d", rv, tc.field, i, tc.i)
		}
	}
}

type testFieldSettable struct {
	Values map[string]string
}

type testStringSettable struct {
	Value string `survey:"string"`
}

type testTaggedStruct struct {
	TaggedValue testStringSettable `survey:"tagged"`
}

type testPtrTaggedStruct struct {
	TaggedValue *testStringSettable `survey:"tagged"`
}

func (t *testFieldSettable) WriteAnswer(name string, value interface{}) error {
	if t.Values == nil {
		t.Values = make(map[string]string)
	}
	if v, ok := value.(string); ok {
		t.Values[name] = v
		return nil
	}
	return fmt.Errorf("incompatible type %T", value)
}

func (t *testStringSettable) WriteAnswer(_ string, value interface{}) error {
	t.Value = value.(string)
	return nil
}
