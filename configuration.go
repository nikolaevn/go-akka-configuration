package configuration

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/tera-insights/go-akka-configuration/hocon"
)

/*
*
 */
func ValueAt(obj interface{}, path string) (interface{}, error) {
	// if the path is "", return this object as the result
	if path == "" {
		return obj, nil
	}
	// if the path is not nil
	// split the path into "key" and "subPath" (based on the position of the first .)
	keyEnd := strings.Index(path, ".")
	var key, subPath string
	if keyEnd == -1 {
		key = path
		subPath = ""
	} else {
		key = path[0:keyEnd]
		subPath = path[keyEnd+1:]
	}

	//fmt.Printf("ValueAt: key=%s, subPath=%s\n", key, subPath)

	// Analyze the type "obj" is of
	switch reflect.TypeOf(obj).Kind() {
	case reflect.Map:
		// If object, recurse on ValueAt(obj[key], subPath)
		m := obj.(map[string]interface{})
		if val, ok := m[key]; ok {
			return ValueAt(val, subPath)
		}
	case reflect.Slice:
		// If array, recurse on ValueAt(obj[int(key)], subPath)
		s := reflect.ValueOf(obj)
		index, err := strconv.Atoi(key)
		if err != nil {
			return nil, errors.New("invalid array index")
		}
		if index >= s.Len() {
			return nil, errors.New("array index out of bounds")
		}
		val := s.Index(index).Interface()
		return ValueAt(val, subPath)
	default:
		// Otherwise return an error that you did not find an element on the path
		return nil, errors.New("element not found")
	}
	return nil, errors.New("element not found")
}

func ParseString(text string, includeCallback ...hocon.IncludeCallback) (interface{}, *map[string]hocon.Position) {
	var callback hocon.IncludeCallback
	if len(includeCallback) > 0 {
		callback = includeCallback[0]
	} else {
		callback = defaultIncludeCallback
	}
	root := hocon.Parse(text, callback)
	if root == nil {
		fmt.Println("debug 6 - root is null here")
	}
	return hocon.TraverseTree(root)
}

func LoadConfig(filename string) (interface{}, *map[string]hocon.Position) {
	data, err := os.ReadFile(filename)
	if err != nil {
		//panic(err)
		fmt.Println("Error:", err)
		return nil, nil
	}
	config, positionMap := ParseString(string(data), defaultIncludeCallback)
	//fmt.Println("Stage - config:", config)
	return config, positionMap
}

func defaultIncludeCallback(filename string) *hocon.HoconRoot {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return hocon.Parse(string(data), defaultIncludeCallback)
}
