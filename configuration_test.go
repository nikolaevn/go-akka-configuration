package configuration

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"sync"
	"testing"

	"github.com/tera-insights/go-akka-configuration/hocon"
)

func TestValueAt(t *testing.T) {
	wg := &sync.WaitGroup{}
	fn := func() {
		defer func() {
			wg.Done()
		}()
		conf, _ := LoadConfig("tests/configs.conf")
		for g := 1; g < 3; g++ {
			for i := 1; i < 4; i++ {
				key := fmt.Sprintf("root.test.o%d.order", i)
				order, err := ValueAt(conf, key)
				if order != int32(i) || err != nil {
					fmt.Println(conf)
					t.Fatalf("order not match,group %d, except: %d, real order: %d", g, i, order)
					return
				}
			}
			conf = nil
			runtime.Gosched()
		}
	}
	wg.Add(2)
	go fn()
	go fn()
	wg.Wait()
}

func Test2(t *testing.T) {
	config, _ := LoadConfig("tests/test.conf")
	expected := map[string]interface{}{
		"server.port":       8080,
		"server.host":       "localhost",
		"server.debug":      true,
		"database.name":     "mydb",
		"database.user":     "user",
		"database.password": "password",
	}

	for key, value := range expected {
		result, err := ValueAt(config, key)
		if err != nil {
			t.Errorf("Error getting value for key %q: %v", key, err)
		}
		if result != value {
			t.Errorf("Incorrect value for key %q: got %v, want %v", key, result, value)
		}
	}
}

func TestGetType(t *testing.T) {
	cases := []struct {
		input    string
		expected hocon.HoconType
	}{
		{"foo", hocon.String},
		{"10", hocon.Int32},
		{"true", hocon.Boolean},
		{"10.1", hocon.Double},
		{"00.122", hocon.Double},
	}

	for _, c := range cases {
		got := hocon.GetType(c.input)
		if got != c.expected {
			t.Errorf("[Error] | GetType(%q) == %v, expected %v", c.input, got, c.expected)
		}
		fmt.Printf("GetType(%q) == %v, expected %v\n", c.input, got, c.expected)
	}

}

func TestParseString(t *testing.T) {
	input := "foo.bar.baz = 42"
	expected := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": map[string]interface{}{
				"baz": 42,
			},
		},
	}
	result, _ := ParseString(input)
	fmt.Println("result =", result)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("unexpected error: %v\n", result)
	}
}

func TestLoadConfig(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "config_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	configText := "foo.bar.baz = 42"
	if _, err := tmpfile.Write([]byte(configText)); err != nil {
		t.Fatal(err)
	}
	config, _ := LoadConfig(tmpfile.Name())
	expectedConfig := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": map[string]interface{}{
				"baz": 42,
			},
		},
	}
	if !reflect.DeepEqual(config, expectedConfig) {
		t.Errorf("LoadConfig(%q) = %#v, expected %#v", tmpfile.Name(), config, expectedConfig)
	}
}

func TestTraverseTree(t *testing.T) {
	hoconStr := "foo.bar.baz = 42\n" +
		"list = [1,2,3]\n" +
		"person {\n" +
		"  name = \"Alice\"\n" +
		"  age = 30\n" +
		"  address {\n" +
		"    city = \"New York\"\n" +
		"    state = \"NY\"\n" +
		"  }\n" +
		"}"

	root := hocon.Parse(hoconStr, defaultIncludeCallback)
	result, posMap := hocon.TraverseTree(root)

	expected := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": map[string]interface{}{
				"baz": int64(42),
			},
		},
		"list": []interface{}{int64(1), int64(2), int64(3)},
		"person": map[string]interface{}{
			"name": "Alice",
			"age":  int64(30),
			"address": map[string]interface{}{
				"city":  "New York",
				"state": "NY",
			},
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("TraverseTree() = %v, expected %v", result, expected)
	}

	expectedPosMap := map[string]hocon.Position{
		"root.foo.bar.baz":          {Line: 1, Col: 15, Len: 2},
		"root.list[0]":              {Line: 2, Col: 6, Len: 1},
		"root.list[1]":              {Line: 2, Col: 8, Len: 1},
		"root.list[2]":              {Line: 2, Col: 10, Len: 1},
		"root.person.name":          {Line: 4, Col: 4, Len: 1},
		"root.person.age":           {Line: 5, Col: 4, Len: 1},
		"root.person.address.city":  {Line: 6, Col: 4, Len: 1},
		"root.person.address.state": {Line: 7, Col: 4, Len: 1},
	}

	if !reflect.DeepEqual(*posMap, expectedPosMap) {
		t.Errorf("TraverseTree() position map = %v, expected %v", *posMap, expectedPosMap)
	}
}
