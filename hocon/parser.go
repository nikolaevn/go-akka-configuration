package hocon

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type IncludeCallback func(filename string) *HoconRoot
type Parser struct {
	reader   *HoconTokenizer
	root     *HoconValue
	callback IncludeCallback

	substitutions []*HoconSubstitution
}

func Parse(text string, callback IncludeCallback) *HoconRoot {
	return new(Parser).parseText(text, callback)
}

func (p *Parser) parseText(text string, callback IncludeCallback) *HoconRoot {
	p.callback = callback
	p.root = NewHoconValue()
	p.reader = NewHoconTokenizer(text)
	p.reader.PullWhitespaceAndComments()
	p.parseObject(p.root, true, "")

	root := NewHoconRoot(p.root)

	cRoot := root.Value()

	for _, sub := range p.substitutions {
		res := getNode(cRoot, sub.Path)
		if res == nil {
			envVal, exist := os.LookupEnv(sub.OrignialPath)
			if !exist {
				if !sub.IsOptional {
					panic("Unresolved substitution:" + sub.Path)
				}
			} else {
				hv := NewHoconValue()
				hv.AppendValue(NewHoconLiteral(envVal))
				sub.ResolvedValue = hv
			}
		} else {
			sub.ResolvedValue = res
		}
	}

	return NewHoconRoot(p.root, p.substitutions...)
}

func (p *Parser) parseObject(owner *HoconValue, root bool, currentPath string) {
	if !owner.IsObject() {
		owner.NewValue(NewHoconObject())
	}

	if owner.IsObject() {
		rootObj := owner
		for rootObj.oldValue != nil {
			oldObj := rootObj.oldValue.GetObject()
			obj := rootObj.GetObject()

			if oldObj == nil || obj == nil {
				break
			}
			obj.Merge(oldObj)
			rootObj = rootObj.oldValue
		}
	}

	currentObject := owner.GetObject()

	for !p.reader.EOF() {
		t := p.reader.PullNext()

		switch t.tokenType {
		case TokenTypeInclude:
			included := p.callback(t.value)
			substitutions := included.substitutions
			for _, substitution := range substitutions {
				substitution.Path = currentPath + "." + substitution.Path
			}
			p.substitutions = append(p.substitutions, substitutions...)
			otherObj := included.value.GetObject()
			owner.GetObject().Merge(otherObj)
		case TokenTypeEoF:
		case TokenTypeKey:
			value := currentObject.GetOrCreateKey(t.value)
			nextPath := t.value
			if len(currentPath) > 0 {
				nextPath = currentPath + "." + t.value
			}
			p.parseKeyContent(value, nextPath)
			if !root {
				return
			}
		case TokenTypeObjectEnd:
			return
		}
	}
}

func (p *Parser) parseKeyContent(value *HoconValue, currentPath string) {
	for !p.reader.EOF() {
		t := p.reader.PullNext()
		switch t.tokenType {
		case TokenTypeDot:
			p.parseObject(value, false, currentPath)
			return
		case TokenTypeAssign:
			{
				if !value.IsObject() {
					value.Clear()
				}
			}
			p.ParseValue(value, false, currentPath)
			return
		case TokenTypePlusAssign:
			{
				if !value.IsObject() {
					value.Clear()
				}
			}
			p.ParseValue(value, true, currentPath)
			return
		case TokenTypeObjectStart:
			p.parseObject(value, true, currentPath)
			return
		}
	}
}

func (p *Parser) ParseValue(owner *HoconValue, isEqualPlus bool, currentPath string) {
	if p.reader.EOF() {
		panic("End of file reached while trying to read a value")
	}

	p.reader.PullWhitespaceAndComments()

	// index reading added here
	startIndex := p.reader.col
	for p.reader.isValue() {
		t := p.reader.PullValue()

		if isEqualPlus {
			sub := p.ParseSubstitution(currentPath, false)
			p.substitutions = append(p.substitutions, sub)
			owner.AppendValue(sub)
		}

		switch t.tokenType {
		case TokenTypeEoF:
		case TokenTypeLiteralValue:
			if owner.IsObject() {
				owner.Clear()
			}
			lit := NewHoconLiteral(t.value)
			owner.AppendValue(lit)
			owner.SetType(String)
		case TokenTypeLiteralValueUnquoted:
			if owner.IsObject() {
				owner.Clear()
			}
			lit := NewHoconLiteral(t.value)
			owner.AppendValue(lit)
			if len(owner.GetType()) == 0 {
				owner.SetType(Unknown)
			}
		case TokenTypeObjectStart:
			p.parseObject(owner, true, currentPath)
		case TokenTypeArrayStart:
			arr := p.ParseArray(currentPath)
			owner.AppendValue(&arr)
			owner.SetType(Array)
		case TokenTypeSubstitute:
			sub := p.ParseSubstitution(t.value, t.isOptional)
			p.substitutions = append(p.substitutions, sub)
			owner.AppendValue(sub)
		}

		if p.reader.IsSpaceOrTab() {
			p.ParseTrailingWhitespace(owner)
		}
	}
	//index ending
	endIndex := p.reader.col

	//populating position
	owner.SetPosition(Position{
		Line: p.reader.GetLine(),
		Col:  startIndex,
		Len:  endIndex - startIndex,
	})
	p.ignoreComma()
	p.ignoreNewline()
}

func (p *Parser) ParseTrailingWhitespace(owner *HoconValue) {
	ws := p.reader.PullSpaceOrTab()
	if len(ws.value) > 0 {
		wsList := NewHoconLiteral(ws.value)
		owner.AppendValue(wsList)
	}
}

func (p *Parser) ParseSubstitution(value string, isOptional bool) *HoconSubstitution {
	return NewHoconSubstitution(value, isOptional)
}

func (p *Parser) ParseArray(currentPath string) HoconArray {
	arr := NewHoconArray()
	for !p.reader.EOF() && !p.reader.IsArrayEnd() {
		v := NewHoconValue()
		p.ParseValue(v, false, currentPath)
		arr.values = append(arr.values, v)
		p.reader.PullWhitespaceAndComments()
	}
	p.reader.PullArrayEnd()
	return *arr
}

func (p *Parser) ignoreComma() {
	if p.reader.IsComma() {
		p.reader.PullComma()
	}
}

func (p *Parser) ignoreNewline() {
	if p.reader.IsNewline() {
		p.reader.PullNewline()
	}
}

func getNode(root *HoconValue, path string) *HoconValue {
	elements := splitDottedPathHonouringQuotes(path)
	currentNode := root

	if currentNode == nil {
		panic("Current node should not be null")
	}

	for _, key := range elements {
		currentNode = currentNode.GetChildObject(key)
		if currentNode == nil {
			return nil
		}
	}
	return currentNode
}

func splitDottedPathHonouringQuotes(path string) []string {
	tmp1 := strings.Split(path, "\"")
	var values []string
	for i := 0; i < len(tmp1); i++ {
		tmp2 := strings.Split(tmp1[i], ".")
		for j := 0; j < len(tmp2); j++ {
			if len(tmp2[j]) > 0 {
				values = append(values, tmp2[j])
			}
		}
	}
	return values
}

// declaring constants
const (
	STRING  = "String" //changed to camel case
	NUMBER  = "Number"
	BOOLEAN = "Boolean"
	NULL    = "Null"
	UNKNOWN = "Unknown"
	INTEGER = "Integer"
)

func TraverseTree(root *HoconRoot) (interface{}, *map[string]Position) {
	positionMap := make(map[string]Position)
	res := traverseHoconValueTree(root.value, "root", &positionMap)
	return res, &positionMap
}

// func traverseHoconValueTree(node *HoconValue, currentPath string, posMap *map[string]Position) interface{} {
// 	// If the node has a position, add it to the position map using the current path
// 	if node.pos != nil {
// 		(*posMap)[currentPath] = Position(*node.pos)
// 	}
// 	// If the node is an object, iterate over its keys and traverse each value recursively
// 	if node.IsObject() {
// 		res := make(map[string]interface{})
// 		object := node.GetObject()
// 		for key := range object.items {
// 			newPath := currentPath + "." + key
// 			val := traverseHoconValueTree(object.items[key], newPath, posMap)
// 			res[key] = val
// 		}
// 		return res
// 	} else if node.IsArray() { // If the node is an array, iterate over its elements and traverse each one recursively
// 		array := node.GetArray()
// 		res := make([]interface{}, len(array))
// 		for i, element := range array {
// 			newKey := currentPath + "[" + strconv.Itoa(i) + "]"
// 			res[i] = traverseHoconValueTree(element, newKey, posMap)
// 		}
// 		return res
// 	} else { // If the node is a literal, extract its value based on its type
// 		switch node.hoconType {
// 		case STRING:
// 			return node.GetString()
// 		case NUMBER:
// 			return node.GetInt64()
// 		case BOOLEAN:
// 			return node.GetBoolean()
// 		case UNKNOWN: // Return nil for unknown nodes (temporary fix for the issue)
// 			return nil
// 		default: // Throw a panic for unexpected value types
// 			panic(fmt.Sprintf("Unexpected value type: %v", node.hoconType))
// 		}
// 	}
// }

func traverseHoconValueTree(node *HoconValue, currentPath string, posMap *map[string]Position) interface{} {
	//handling nil case before dereferinceing
	if node.pos != nil {
		(*posMap)[currentPath] = Position(*node.pos)
	}

	if node.IsObject() {
		res := make(map[string]interface{})
		object := node.GetObject()
		for key := range object.items {
			newPath := currentPath + "." + key
			//fmt.Printf("[1] %s: %v\n", newPath, object.items[key])
			val := traverseHoconValueTree(object.items[key], newPath, posMap)
			res[key] = val
		}
		return res
	} else if node.IsArray() {
		array := node.GetArray()
		res := make([]interface{}, len(array))
		for i, element := range array {
			newKey := currentPath + "[" + strconv.Itoa(i) + "]"
			//fmt.Printf("[2] %s: %v\n", newKey, element.hoconType)
			res[i] = traverseHoconValueTree(element, newKey, posMap)
		}
		return res
	} else {
		// Extract the value of the literal based on its type
		//fmt.Println("[3] hoconType- ", string(node.hoconType))
		switch node.hoconType {
		case STRING:
			return node.GetString()
		case NUMBER:
			return node.GetInt64()
		case BOOLEAN:
			return node.GetBoolean()
		case UNKNOWN: //added to fix unknown issue (temp)
			return nil
		default:
			panic(fmt.Sprintf("Unexpected value type: %v", node.hoconType))
		}
	}
}
