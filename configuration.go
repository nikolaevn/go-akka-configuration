package configuration

import (
	"io/ioutil"

	"github.com/tera-insights/go-akka-configuration/hocon"
)

/*
*
 */
func ValueAt(obj interface{}, path string) (interface{}, error) {
	// if the path is "", return this object as the result

	// if the path is not nil
	// split the path into "key" and "subPath" (based on the position of the first .)

	// Analyze the type "obj" is of

	// If object, recuse on ValueAt(obj[key], subPath)

	// If array, recurse on ValueAt(obj[int(key)], subPath)

	// Otherwise return an error that you did not find an element on the path

	return obj, nil
}

func ParseString(text string, includeCallback ...hocon.IncludeCallback) (interface{}, *map[string]hocon.Position) {
	var callback hocon.IncludeCallback
	if len(includeCallback) > 0 {
		callback = includeCallback[0]
	} else {
		callback = defaultIncludeCallback
	}
	root := hocon.Parse(text, callback)

	return hocon.TraverseTree(root)
}

func LoadConfig(filename string) (interface{}, *map[string]hocon.Position) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return ParseString(string(data), defaultIncludeCallback)
}

func defaultIncludeCallback(filename string) *hocon.HoconRoot {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return hocon.Parse(string(data), defaultIncludeCallback)
}
