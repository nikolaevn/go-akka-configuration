package configuration

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

func TestParseKeyOrder(t *testing.T) {

	wg := &sync.WaitGroup{}

	fn := func() {

		defer func() {
			wg.Done()
		}()

		for i := 0; i < 100000; i++ {
			conf, _ := LoadConfig("tests/configs.conf")
			for g := 1; g < 3; g++ {
				for i := 1; i < 4; i++ {
					key := fmt.Sprintf("test.out.a.b.c.d.groups.g%d.o%d.order", g, i)
					order, err := ValueAt(conf, key)
					//fmt.Println("printing key value here : ", key, order, err)

					if order != int32(i) || err != nil {
						//fmt.Println(conf)
						t.Fatalf("order not match,group %d, except: %d, real order: %d", g, i, order)
						return
					}
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

// func TestValueAt(t *testing.T) {
// 	config, _ := LoadConfig("tests/configs.conf")

// 	testCases := []struct {
// 		path        string
// 		expected    interface{}
// 		expectedErr error
// 	}{
// 		{
// 			path:        "",
// 			expected:    config,
// 			expectedErr: nil,
// 		},
// 		{
// 			path:        "test.booleanValue",
// 			expected:    true,
// 			expectedErr: nil,
// 		},
// 		{
// 			path:        "test.stringValue",
// 			expected:    "test value",
// 			expectedErr: nil,
// 		},
// 		{
// 			path:        "test.nested.value1",
// 			expected:    1,
// 			expectedErr: nil,
// 		},
// 		{
// 			path:        "test.nested.value2",
// 			expected:    2,
// 			expectedErr: nil,
// 		},
// 		{
// 			path:        "test.array[0]",
// 			expected:    "first",
// 			expectedErr: nil,
// 		},
// 		{
// 			path:        "test.array[1]",
// 			expected:    "second",
// 			expectedErr: nil,
// 		},
// 		{
// 			path:        "test.array[2].nestedArray[0]",
// 			expected:    10,
// 			expectedErr: nil,
// 		},
// 		{
// 			path:        "test.invalidKey",
// 			expected:    nil,
// 			expectedErr: fmt.Errorf("element not found for path: test.invalidKey"),
// 		},
// 		{
// 			path:        "test.array[3]",
// 			expected:    nil,
// 			expectedErr: fmt.Errorf("element not found for path: test.array[3]"),
// 		},
// 		{
// 			path:        "test.array[2].nestedArray[2]",
// 			expected:    nil,
// 			expectedErr: fmt.Errorf("element not found for path: test.array[2].nestedArray[2]"),
// 		},
// 	}

// 	for _, tc := range testCases {
// 		actual, err := ValueAt(config, tc.path)
// 		if !reflect.DeepEqual(actual, tc.expected) {
// 			t.Errorf("Expected %v, but got %v for path %s", tc.expected, actual, tc.path)
// 		}
// 		if !reflect.DeepEqual(err, tc.expectedErr) {
// 			t.Errorf("Expected error %v, but got %v for path %s", tc.expectedErr, err, tc.path)
// 		}
// 	}
// }
