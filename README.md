Hocon Repository used to parse hocon config file in confcheck

this branch is only for testing and temp changes

[summary of some changes]

### config.go

- provides a Config type, which stores configuration data in a hierarchical structure and provides functions to access the values of specific nodes in that structure.
- provides a number of methods to extract values from the Config object, such as GetBoolean(), GetByteSize(), GetInt32(), GetString(), and GetTimeDuration(), among others. These methods access the values of specific nodes in the configuration hierarchy and return those values in the appropriate data type. 
- provides a NewConfigFromRoot() function to create a new Config object from a HoconRoot object, and a NewConfigFromConfig() function to create a new Config object from an existing Config object.

### parser.go
- constants define the types of values that can be found in the configuration file, such as STRING, NUMBER, BOOLEAN, NULL, INTEGER, and UNKNOWN.

- The TraverseTree function takes a HoconRoot object, traverses the HoconValue tree structure it contains, and returns a nested interface{} structure representing the configuration data. It also returns a map containing the positions of each value in the configuration file.

- The traverseHoconValueTree function is a recursive function that traverses the HoconValue tree structure and extracts the values it contains. It returns the extracted value as an interface{} object.

### configuration_test.go
- The ValueAt function takes an object and a path string, and recursively traverses the object to find the value at the specified path. The path is assumed to be a period-separated string of keys, where each key represents either a map key or an index into an array. The function returns the value at the specified path, or an error if the path is invalid or the value is not found.

- The ParseString function takes a string in the HOCON format and an optional include callback function, and returns the parsed configuration as a nested map of strings to interface{}. The include callback function is used to handle include directives in the configuration file.

- The LoadConfig function takes a filename and uses os.ReadFile to read the file contents into a string. It then calls ParseString to parse the string into a configuration map.

