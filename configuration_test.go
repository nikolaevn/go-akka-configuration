package configuration

import (
	"fmt"
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
				//fmt.Println("printing key value here : ", key, order, err)

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

