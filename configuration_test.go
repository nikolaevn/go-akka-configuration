package configuration

import (
	"fmt"
	"testing"
)

// func Test1(t *testing.T) {

// 	wg := &sync.WaitGroup{}

// 	fn := func() {

// 		defer func() {
// 			wg.Done()
// 		}()

// 		conf, _ := LoadConfig("tests/configs.conf")
// 		for g := 1; g < 3; g++ {
// 			for i := 1; i < 4; i++ {
// 				key := fmt.Sprintf("root.test.o%d.order", i)
// 				order, err := ValueAt(conf, key)
// 				//fmt.Println("printing key value here : ", key, order, err)

// 				if order != int32(i) || err != nil {
// 					fmt.Println(conf)
// 					t.Fatalf("order not match,group %d, except: %d, real order: %d", g, i, order)
// 					return
// 				}
// 			}
// 			conf = nil
// 			runtime.Gosched()
// 		}
// 	}

// 	wg.Add(2)
// 	go fn()
// 	go fn()
// 	wg.Wait()
// }

func Test2(t *testing.T) {
	// Parse the configuration data
	config, _ := LoadConfig("tests/test.conf")

	// Test the ValueAt function for a few sample keys
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
	fmt.Println("Printing Loaded config : ", config)
}
