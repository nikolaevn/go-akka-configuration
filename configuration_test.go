package configuration

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

func Test1(t *testing.T) {

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
					//fmt.Println(conf)
					//t.Fatalf("order not match,group %d, except: %d, real order: %d", g, i, order)
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

// func TestParseKeyOrder(t *testing.T) {

// 	wg := &sync.WaitGroup{}

// 	fn := func() {

// 		defer func() {
// 			wg.Done()
// 		}()

// 		for i := 0; i < 10; i++ {
// 			conf, _ := LoadConfig("tests/configs.conf")
// 			key := fmt.Sprintf("root.test.o%d.order", i)
// 			ValueAt(conf, key)
// 		}
// 		//conf = nil
// 		runtime.Gosched()

// 	}

// 	wg.Add(2)

// 	go fn()
// 	go fn()

// 	wg.Wait()
// }
