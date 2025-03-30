package main

import (
	"fmt"
	"strconv"
	"strings"
)

// LispValue represents a value
type LispValue interface {
	String() string
}

// LispAtom represents an atomic value (symbol)
type LispAtom struct {
	Value string
}

// String returns the string representation of the atom
func (a *LispAtom) String() string {
	return a.Value
}

// LispNumber represents a numeric value
type LispNumber struct {
	Value int
}

// String returns the string representation of the number
func (n *LispNumber) String() string {
	return strconv.Itoa(n.Value)
}

// LispFloat represents a float value
type LispFloat struct {
	Value float64
}

// String returns the string representation of the float
func (f *LispFloat) String() string {
	return strconv.FormatFloat(f.Value, 'f', -1, 64)
}

// LispString represents a string value
type LispString struct {
	Value string
}

// String returns the string representation of the string
func (s *LispString) String() string {
	return "\"" + s.Value + "\""
}

// LispList represents a list of Lisp values
type LispList struct {
	Elements []LispValue
}

// String returns the string representation of the list
func (l *LispList) String() string {
	var sb strings.Builder
	sb.WriteString(string(OPEN_BRACKET))
	for i, elem := range l.Elements {
		sb.WriteString(elem.String())
		if i < len(l.Elements)-1 {
			sb.WriteString(EMPTY_STRING)
		}
	}
	sb.WriteString(string(CLOSE_BRACKET))
	return sb.String()
}

// LispFunction represents a user-defined function
type LispFunction struct {
	Name   *LispAtom
	Params []LispValue
	Body   LispValue
	Env    Environment
}

// String returns the string representation of the function
func (f *LispFunction) String() string {
	if f.Name != nil {
		return strings.ToUpper(f.Name.Value)
	}
	return FUNCTION
}

// LispBoolean represents a boolean value
type LispBoolean struct {
	Value bool
}

// String returns the string representation of the boolean
func (b *LispBoolean) String() string {
	if b.Value {
		return TRUE
	}
	return FALSE
}

// LispNil represents a nil/null value
type LispNil struct{}

// String returns the string representation of the nil
func (n *LispNil) String() string {
	return NIL
}

// Environment represents the mapping of symbols to their values
type Environment map[string]LispValue

// LispError represents an error with line and column information
type LispError struct {
	Message string
	Line    int
	Column  int
}

// Error returns the error message
func (e *LispError) Error() string {
	return fmt.Sprintf("Error at line %d, column %d: %s", e.Line, e.Column, e.Message)
}
