package jsonconf

import (
	"reflect"
	"testing"
)

func Test(t *testing.T) {
	str := `
	{
		"int": 10,
		"str_array": ["one", "two", "three"],
		"int_array": [0, 1, 2, 3],
		"nested": {
			"int": 11,
			"str": "tvelwe"
		}
	}
	`
	c, err := ReadString(str)
	if err != nil {
		t.Error("ReadString", err)
		return
	}
	t.Logf("int:%+v\n", c("int", 110))
	if val, ok := c("int", 110).(int64); !ok || val != 10 {
		t.Error("invalid 'int'")
	}
	t.Logf("int_not_existing type:%v\n", reflect.TypeOf(c("int_not_existing", 110)))
	if val, ok := c("int_not_existing", 110).(int64); !ok || val != 110 {
		t.Error("invalid 'int_not_existing'")
	}
	if val, ok := c("nested.int", 111).(int64); !ok || val != 11 {
		t.Error("invalid 'nested.int'")
	}
	t.Logf("str_array type:%v\n", reflect.TypeOf(c("str_array", 110)))
	if s, ok := c("str_array", []string{"foo", "bar"}).([]string); !ok || len(s) != 3 || s[2] != "three" {
		t.Errorf("invalid 'str_array' '%s'\n", s[2])
	}
	if s, ok := c("str_array_not_existing", []string{"foo", "bar"}).([]string); !ok || len(s) != 2 || s[1] != "bar" {
		t.Error("invalid 'str_array_not_existing'")
	}
	if s, ok := c("int_array", []int64{100, 101}).([]int64); !ok || len(s) != 4 || s[2] != 2 {
		t.Errorf("invalid 'str_array' '%s'\n", s[2])
	}
}

func TestInvalidConf(t *testing.T) {
	c, err := ReadString("invalid json")
	if err == nil {
		t.Error("ReadString must fail")
		return
	}
	if val, ok := c("foo", "bar").(string); !ok || val != "bar" {
		t.Fail()
	}
	if val, ok := c("foo.bar", 10).(int64); !ok || val != 10 {
		t.Fail()
	}
}

