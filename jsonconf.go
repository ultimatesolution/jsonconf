// A minimalistic JSON config reader.
// No nasty struct definitions, int client code, just one-line retrieving
// values by keys.
// All numeric values in JSON configs are converted to int64 type.
// JSON arrays are returned as []string and []int64.
//
// Usage:
//
// conf = jsonconf.ReadFile("config.json")
// port, ok := conf("tcp.port", 80).(int)
// ...
// host, ok := conf("tcp.host", "microsoft.com").(string)
// ...
// clusters, ok := conf("clusters", []string{}{"default"}).([]string)
// ...
// /* or the same without ", ok" if you are sure about your config, e.g.: */
// name := conf("client.name", "noname").(string)
//
// (c) 2016 Eugene Korenevsky

package jsonconf

import (
	"reflect"
	"io/ioutil"
	"github.com/Jeffail/gabs"
)

type Config func (key string, defVal interface{}) interface{}

func convertToStringSlice(s reflect.Value) interface{} {
	ret := make([]string, s.Len())
	for i := 0; i < s.Len(); i++ {
		iface := s.Index(i).Interface()
		ret[i] = reflect.ValueOf(iface).String()
	}
	return ret
}

func convertToIntSlice(s reflect.Value) interface{} {
	ret := make([]int64, s.Len())
	for i := 0; i < s.Len(); i++ {
		iface := s.Index(i).Interface()
		ret[i] = int64(reflect.ValueOf(iface).Float())
	}
	return ret
}

// Preprocess result to make it more convenient for usage in the caller's code.
// Convert all numbers to int64
// Convert []interface{} to []string or []int64 for JSON arrays.
func preprocess(v interface{}) interface{} {
	s := reflect.ValueOf(v)
	switch s.Kind()  {
	case reflect.Float32, reflect.Float64:
		return int64(s.Float())
	case reflect.Int:
		return s.Int()
	}
	if s.Kind() != reflect.Slice || s.Len() == 0 ||
				s.Index(0).Kind() != reflect.Interface {
		return v
	}
	iface := s.Index(0).Interface()
	switch reflect.ValueOf(iface).Kind() {
	case reflect.String:
		return convertToStringSlice(s)
	case reflect.Int, reflect.Float32, reflect.Float64:
		return convertToIntSlice(s)
	default:
		return v
	}
}

func read(data []byte) (c Config, err error) {
	json, err := gabs.ParseJSON(data)
	if err != nil {
		// Use default values if the config has not been read
		c = func (key string, defVal interface{}) interface{} {
			return preprocess(defVal)
		}
		return
	}
	c = func (key string, defVal interface{}) interface{} {
		if !json.ExistsP(key) {
			return preprocess(defVal)
		}
		v := json.Path(key).Data()
		return preprocess(v)
	}
	return
}

func ReadString(s string) (c Config, err error) {
	c, err = read([]byte(s))
	return
}

func ReadFile(filename string) (c Config, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	c, err = read(data)
	return
}

