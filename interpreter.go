package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

// Eval evaluates a Lisp expression in the given environment
func Eval(env Environment, expr LispValue) (LispValue, error) {
	switch v := expr.(type) {
	case *LispAtom:
		if val, ok := env[v.Value]; ok {
			return val, nil
		}
		return nil, &LispError{Message: fmt.Sprintf("unbound symbol: %s", v.Value), Line: 0, Column: 0}
	case *LispNumber, *LispFloat, *LispString, *LispBoolean, *LispNil:
		return v, nil
	case *LispList:
		if len(v.Elements) == 0 {
			return v, nil
		}
		fn, ok := v.Elements[0].(*LispAtom)
		if !ok {
			return nil, &LispError{Message: fmt.Sprintf("invalid function call: %v", v.Elements[0]), Line: 0, Column: 0}
		}
		args := v.Elements[1:]
		switch fn.Value {
		case FORMAT:
			return builtinFormat(env, args)
		case READ:
			return builtinRead(env, args)
		case PRINT:
			return builtinPrint(env, args)
		case PLUS:
			return builtinAdd(env, args)
		case MINUS:
			return builtinSub(env, args)
		case STAR:
			return builtinMul(env, args)
		case SLASH:
			return builtinDiv(env, args)
		case PERCENT:
			return builtinMod(env, args)
		case POW:
			return builtinPow(env, args)
		case SQRT:
			return builtinSqrt(env, args)
		case CONCAT:
			return builtinConcat(env, args)
		case SUBSTRING:
			return builtinSubstring(env, args)
		case IS_NUMBER:
			return builtinIsNumber(env, args)
		case IS_STRING:
			return builtinIsString(env, args)
		case LESS_THAN:
			return builtinLt(env, args)
		case LESS_OR_EQUAL_THAN:
			return builtinLtOrEq(env, args)
		case GREATER_THAN:
			return builtinGt(env, args)
		case GREATER_OR_EQUAL_THAN:
			return builtinGtOrEq(env, args)
		case EQUAL:
			return builtinEq(env, args)
		case IF:
			return builtinIf(env, args)
		case DEFUN:
			return builtinDefun(env, args)
		case LAMBDA:
			return builtinLambda(env, args)
		case LET:
			return builtinLet(env, args)
		case AND:
			return builtinAnd(env, args)
		case OR:
			return builtinOr(env, args)
		case NOT:
			return builtinNot(env, args)
		case LIST:
			return builtinList(args)
		case CAR:
			return builtinCar(env, args)
		case CDR:
			return builtinCdr(env, args)
		case CONS:
			return builtinCons(env, args)
		case LENGTH:
			return builtinLength(env, args)
		case APPEND:
			return builtinAppend(env, args)
		default:
			return callFunction(env, fn.Value, args)
		}
	default:
		return nil, &LispError{Message: fmt.Sprintf("unknown expression type: %T", v), Line: 0, Column: 0}
	}
}

// Helper function to convert Lisp values to Go values
func lispValueToGoValue(value LispValue) interface{} {
	switch v := value.(type) {
	case *LispNumber:
		return v.Value
	case *LispFloat:
		return v.Value
	case *LispString:
		return v.Value
	case *LispAtom:
		return v.Value
	case *LispBoolean:
		return v.Value
	case *LispNil:
		return nil
	default:
		return v
	}
}

// Built-in function implementations

// builtinFormat is the implementation of the format function
func builtinFormat(env Environment, args []LispValue) (LispValue, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("wrong number of arguments to format")
	}

	// The second argument is the format string
	formatStr, ok := args[1].(*LispString)
	if !ok {
		return nil, fmt.Errorf("invalid format string: %v", args[1])
	}

	// The remaining arguments are the values to be formatted
	formatArgs := args[2:]

	// Prepare the arguments for fmt.Sprintf
	var sprintfArgs []interface{}
	for _, arg := range formatArgs {
		evalArg, err := Eval(env, arg)
		if err != nil {
			return nil, err
		}
		sprintfArgs = append(sprintfArgs, lispValueToGoValue(evalArg))
	}

	// Perform the formatting
	formattedStr := fmt.Sprintf(formatStr.Value, sprintfArgs...)

	return &LispString{Value: formattedStr}, nil
}

// builtinRead reads input from the user
func builtinRead(_ Environment, args []LispValue) (LispValue, error) {
	scanner := bufio.NewScanner(os.Stdin)
	if len(args) > 0 {
		for _, arg := range args {
			fmt.Print(arg.String())
		}
	}
	scanner.Scan()
	input := scanner.Text()
	return &LispString{Value: input}, nil
}

// builtinPrint prints a Lisp value to the console
func builtinPrint(env Environment, args []LispValue) (LispValue, error) {
	for _, arg := range args {
		val, err := Eval(env, arg)
		if err != nil {
			return nil, err
		} else {
			return val, nil
		}
	}
	return &LispString{}, nil
}

// builtinAdd is built-in implementation of addition operation
func builtinAdd(env Environment, args []LispValue) (LispValue, error) {
	var sum float64
	for _, arg := range args {
		val, err := Eval(env, arg)
		if err != nil {
			return nil, err
		}
		switch v := val.(type) {
		case *LispNumber:
			sum += float64(v.Value)
		case *LispFloat:
			sum += v.Value
		default:
			return nil, &LispError{Message: fmt.Sprintf("invalid argument to +: %v", val), Line: 0, Column: 0}
		}
	}
	if float64(int(sum)) == sum {
		return &LispNumber{Value: int(sum)}, nil
	}
	return &LispFloat{Value: sum}, nil
}

// builtinSub is built-in implementation of subtraction operation
func builtinSub(env Environment, args []LispValue) (LispValue, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("wrong number of arguments to -")
	}
	val, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	var diff float64
	switch v := val.(type) {
	case *LispNumber:
		_, ok := val.(*LispNumber)
		if !ok {
			return nil, fmt.Errorf("invalid argument to -: %v", val)
		}
		diff = float64(v.Value)
	case *LispFloat:
		_, ok := val.(*LispFloat)
		if !ok {
			return nil, fmt.Errorf("invalid argument to -: %v", val)
		}
		diff = float64(v.Value)
	default:
		return nil, &LispError{Message: fmt.Sprintf("invalid argument to +: %v", val), Line: 0, Column: 0}
	}
	for _, arg := range args[1:] {
		val, err := Eval(env, arg)
		if err != nil {
			return nil, err
		}
		switch v := val.(type) {
		case *LispNumber:
			diff -= float64(v.Value)
		case *LispFloat:
			diff -= v.Value
		default:
			return nil, &LispError{Message: fmt.Sprintf("invalid argument to +: %v", val), Line: 0, Column: 0}
		}
	}
	if float64(int(diff)) == diff {
		return &LispNumber{Value: int(diff)}, nil
	}
	return &LispFloat{Value: diff}, nil
}

// builtinMul is built-in implementation of multiplication operation
func builtinMul(env Environment, args []LispValue) (LispValue, error) {
	var prod float64
	prod = 1
	for _, arg := range args {
		val, err := Eval(env, arg)
		if err != nil {
			return nil, err
		}
		switch v := val.(type) {
		case *LispNumber:
			prod *= float64(v.Value)
		case *LispFloat:
			prod *= v.Value
		default:
			return nil, &LispError{Message: fmt.Sprintf("invalid argument to +: %v", val), Line: 0, Column: 0}
		}
	}
	if float64(int(prod)) == prod {
		return &LispNumber{Value: int(prod)}, nil
	}
	return &LispFloat{Value: prod}, nil
}

// builtinDiv is built-in implementation of division operation
func builtinDiv(env Environment, args []LispValue) (LispValue, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("wrong number of arguments to /")
	}

	val, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}

	var quot float64
	switch v := val.(type) {
	case *LispNumber:
		_, ok := val.(*LispNumber)
		if !ok {
			return nil, fmt.Errorf("invalid argument to -: %v", val)
		}
		quot = float64(v.Value)
	case *LispFloat:
		_, ok := val.(*LispFloat)
		if !ok {
			return nil, fmt.Errorf("invalid argument to -: %v", val)
		}
		quot = float64(v.Value)
	default:
		return nil, &LispError{Message: fmt.Sprintf("invalid argument to +: %v", val), Line: 0, Column: 0}
	}

	for _, arg := range args[1:] {
		val, err := Eval(env, arg)
		if err != nil {
			return nil, err
		}

		switch v := val.(type) {
		case *LispNumber:
			if v.Value == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			quot /= float64(v.Value)
		case *LispFloat:
			if v.Value == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			quot /= v.Value
		default:
			return nil, &LispError{Message: fmt.Sprintf("invalid argument to +: %v", val), Line: 0, Column: 0}
		}
	}
	if float64(int(quot)) == quot {
		return &LispNumber{Value: int(quot)}, nil
	}
	return &LispFloat{Value: quot}, nil
}

// builtinMod is built-in implementation of modulo operation
func builtinMod(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 2 {
		return nil, &LispError{Message: "wrong number of arguments to %", Line: 0, Column: 0}
	}
	val1, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	val2, err := Eval(env, args[1])
	if err != nil {
		return nil, err
	}
	num1, ok1 := val1.(*LispNumber)
	num2, ok2 := val2.(*LispNumber)
	if !ok1 || !ok2 {
		return nil, &LispError{Message: "invalid arguments to %", Line: 0, Column: 0}
	}
	if num2.Value == 0 {
		return nil, &LispError{Message: "division by zero", Line: 0, Column: 0}
	}
	return &LispNumber{Value: num1.Value % num2.Value}, nil
}

// builtinPow is built-in implementation of pow operation
func builtinPow(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 2 {
		return nil, &LispError{Message: "wrong number of arguments to pow", Line: 0, Column: 0}
	}
	base, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	exp, err := Eval(env, args[1])
	if err != nil {
		return nil, err
	}
	baseVal, expVal := 0.0, 0.0
	switch v := base.(type) {
	case *LispNumber:
		baseVal = float64(v.Value)
	case *LispFloat:
		baseVal = v.Value
	default:
		return nil, &LispError{Message: "invalid base argument to pow", Line: 0, Column: 0}
	}
	switch v := exp.(type) {
	case *LispNumber:
		expVal = float64(v.Value)
	case *LispFloat:
		expVal = v.Value
	default:
		return nil, &LispError{Message: "invalid exponent argument to pow", Line: 0, Column: 0}
	}
	result := math.Pow(baseVal, expVal)
	if float64(int(result)) == result {
		return &LispNumber{Value: int(result)}, nil
	}
	return &LispFloat{Value: result}, nil
}

// builtinSqrt is built-in implementation of sqrt operation
func builtinSqrt(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 1 {
		return nil, &LispError{Message: "wrong number of arguments to sqrt", Line: 0, Column: 0}
	}
	val, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	var num float64
	switch v := val.(type) {
	case *LispNumber:
		num = float64(v.Value)
	case *LispFloat:
		num = v.Value
	default:
		return nil, &LispError{Message: "invalid argument to sqrt", Line: 0, Column: 0}
	}
	if num < 0 {
		return nil, &LispError{Message: "cannot take square root of negative number", Line: 0, Column: 0}
	}
	result := math.Sqrt(num)
	if float64(int(result)) == result {
		return &LispNumber{Value: int(result)}, nil
	}
	return &LispFloat{Value: result}, nil
}

// builtinConcat is built-in implementation of concat operation
func builtinConcat(env Environment, args []LispValue) (LispValue, error) {
	var result strings.Builder
	for _, arg := range args {
		val, err := Eval(env, arg)
		if err != nil {
			return nil, err
		}
		str, ok := val.(*LispString)
		if !ok {
			return nil, &LispError{Message: "invalid argument to concat", Line: 0, Column: 0}
		}
		result.WriteString(str.Value)
	}
	return &LispString{Value: result.String()}, nil
}

// builtinSubstring is built-in implementation of substring operation
func builtinSubstring(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 3 {
		return nil, &LispError{Message: "wrong number of arguments to substring", Line: 0, Column: 0}
	}
	str, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	start, err := Eval(env, args[1])
	if err != nil {
		return nil, err
	}
	end, err := Eval(env, args[2])
	if err != nil {
		return nil, err
	}
	strVal, ok := str.(*LispString)
	if !ok {
		return nil, &LispError{Message: "first argument to substring must be a string", Line: 0, Column: 0}
	}
	startVal, ok := start.(*LispNumber)
	if !ok {
		return nil, &LispError{Message: "second argument to substring must be a number", Line: 0, Column: 0}
	}
	endVal, ok := end.(*LispNumber)
	if !ok {
		return nil, &LispError{Message: "third argument to substring must be a number", Line: 0, Column: 0}
	}
	if startVal.Value < 0 || endVal.Value > len(strVal.Value) || startVal.Value > endVal.Value {
		return nil, &LispError{Message: "invalid substring range", Line: 0, Column: 0}
	}
	return &LispString{Value: strVal.Value[startVal.Value:endVal.Value]}, nil
}

// builtinIsNumber is built-in implementation of isNumber operation
func builtinIsNumber(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 1 {
		return nil, &LispError{Message: "wrong number of arguments to is-number", Line: 0, Column: 0}
	}
	val, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	_, isNum := val.(*LispNumber)
	_, isFloat := val.(*LispFloat)
	return &LispBoolean{Value: isNum || isFloat}, nil
}

// builtinIsString is built-in implementation of isString operation
func builtinIsString(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 1 {
		return nil, &LispError{Message: "wrong number of arguments to is-string", Line: 0, Column: 0}
	}
	val, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	_, isString := val.(*LispString)
	return &LispBoolean{Value: isString}, nil
}

// builtinLt is built-in implementation of less than condition
func builtinLt(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments to <")
	}
	val1, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	val2, err := Eval(env, args[1])
	if err != nil {
		return nil, err
	}
	num1, ok := val1.(*LispNumber)
	if !ok {
		return nil, fmt.Errorf("invalid argument to <: %v", val1)
	}
	num2, ok := val2.(*LispNumber)
	if !ok {
		return nil, fmt.Errorf("invalid argument to <: %v", val2)
	}
	if num1.Value < num2.Value {
		return &LispAtom{Value: "true"}, nil
	}
	return &LispAtom{Value: "false"}, nil
}

// builtinLtOrEq is built-in implementation of less or equal than condition
func builtinLtOrEq(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments to <")
	}
	val1, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	val2, err := Eval(env, args[1])
	if err != nil {
		return nil, err
	}
	num1, ok := val1.(*LispNumber)
	if !ok {
		return nil, fmt.Errorf("invalid argument to <: %v", val1)
	}
	num2, ok := val2.(*LispNumber)
	if !ok {
		return nil, fmt.Errorf("invalid argument to <: %v", val2)
	}
	if num1.Value <= num2.Value {
		return &LispAtom{Value: "true"}, nil
	}
	return &LispAtom{Value: "false"}, nil
}

// builtinGt is built-in implementation of greater than condition
func builtinGt(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments to >")
	}
	val1, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	val2, err := Eval(env, args[1])
	if err != nil {
		return nil, err
	}
	num1, ok := val1.(*LispNumber)
	if !ok {
		return nil, fmt.Errorf("invalid argument to >: %v", val1)
	}
	num2, ok := val2.(*LispNumber)
	if !ok {
		return nil, fmt.Errorf("invalid argument to >: %v", val2)
	}
	if num1.Value > num2.Value {
		return &LispAtom{Value: "true"}, nil
	}
	return &LispAtom{Value: "false"}, nil
}

// builtinGtOrEq is built-in implementation of greater or equal than condition
func builtinGtOrEq(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments to >")
	}
	val1, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	val2, err := Eval(env, args[1])
	if err != nil {
		return nil, err
	}
	num1, ok := val1.(*LispNumber)
	if !ok {
		return nil, fmt.Errorf("invalid argument to >: %v", val1)
	}
	num2, ok := val2.(*LispNumber)
	if !ok {
		return nil, fmt.Errorf("invalid argument to >: %v", val2)
	}
	if num1.Value >= num2.Value {
		return &LispAtom{Value: "true"}, nil
	}
	return &LispAtom{Value: "false"}, nil
}

// builtinEq is built-in implementation of equal to condition
func builtinEq(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments to =")
	}
	val1, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	val2, err := Eval(env, args[1])
	if err != nil {
		return nil, err
	}
	num1, ok1 := val1.(*LispNumber)
	num2, ok2 := val2.(*LispNumber)
	if ok1 && ok2 {
		if num1.Value == num2.Value {
			return &LispAtom{Value: "true"}, nil
		}
		return &LispAtom{Value: "false"}, nil
	}
	if val1.String() == val2.String() {
		return &LispAtom{Value: "true"}, nil
	}
	return &LispAtom{Value: "false"}, nil
}

// builtinIf is built-in implementation of if conditional struct
func builtinIf(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("wrong number of arguments to if")
	}
	cond, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	if atom, ok := cond.(*LispAtom); ok && atom.Value == "true" {
		return Eval(env, args[1])
	}
	return Eval(env, args[2])
}

// builtinDefun is built-in implementation of function definition
func builtinDefun(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("wrong number of arguments to defun")
	}
	name, ok := args[0].(*LispAtom)
	if !ok {
		return nil, fmt.Errorf("invalid function name: %v", args[0])
	}
	params, ok := args[1].(*LispList)
	if !ok {
		return nil, fmt.Errorf("invalid function parameters: %v", args[1])
	}
	fn := &LispFunction{Name: name, Params: params.Elements, Body: args[2], Env: env}
	env[name.Value] = fn
	return fn, nil
}

// builtinLambda is built-in implementation of lambda function definition
func builtinLambda(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments to lambda")
	}
	params, ok := args[0].(*LispList)
	if !ok {
		return nil, fmt.Errorf("invalid lambda parameters: %v", args[0])
	}
	return &LispFunction{Params: params.Elements, Body: args[1], Env: env}, nil
}

// builtinLet is built-in implementation of let local variable definition
func builtinLet(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments to let")
	}
	bindings, ok := args[0].(*LispList)
	if !ok {
		return nil, fmt.Errorf("invalid let bindings: %v", args[0])
	}
	localEnv := make(Environment)
	for key, value := range env {
		localEnv[key] = value
	}
	for _, binding := range bindings.Elements {
		bindList, ok := binding.(*LispList)
		if !ok || len(bindList.Elements) != 2 {
			return nil, fmt.Errorf("invalid let binding: %v", binding)
		}
		key, ok := bindList.Elements[0].(*LispAtom)
		if !ok {
			return nil, fmt.Errorf("invalid let binding key: %v", bindList.Elements[0])
		}
		val, err := Eval(localEnv, bindList.Elements[1])
		if err != nil {
			return nil, err
		}
		localEnv[key.Value] = val
	}
	return Eval(localEnv, args[1])
}

// builtinAnd is built-in implementation of and logical operation
func builtinAnd(env Environment, args []LispValue) (LispValue, error) {
	for _, arg := range args {
		val, err := Eval(env, arg)
		if err != nil {
			return nil, err
		}
		if boolean, ok := val.(*LispBoolean); ok && !boolean.Value {
			return &LispBoolean{Value: false}, nil
		}
	}
	return &LispBoolean{Value: true}, nil
}

// builtinOr is built-in implementation of or logical operation
func builtinOr(env Environment, args []LispValue) (LispValue, error) {
	for _, arg := range args {
		val, err := Eval(env, arg)
		if err != nil {
			return nil, err
		}
		if boolean, ok := val.(*LispBoolean); ok && boolean.Value {
			return &LispBoolean{Value: true}, nil
		}
	}
	return &LispBoolean{Value: false}, nil
}

// builtinNot is built-in implementation of not logical operation
func builtinNot(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 1 {
		return nil, &LispError{Message: "wrong number of arguments to not", Line: 0, Column: 0}
	}
	val, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	if boolean, ok := val.(*LispBoolean); ok {
		return &LispBoolean{Value: !boolean.Value}, nil
	}
	return &LispBoolean{Value: false}, nil
}

// builtinList is built-in implementation of list definition
func builtinList(args []LispValue) (LispValue, error) {
	return &LispList{Elements: args}, nil
}

// builtinCar is built-in implementation of car list operation. It retrieves first element of a list.
func builtinCar(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of arguments to car")
	}
	val, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	list, ok := val.(*LispList)
	if !ok {
		return nil, fmt.Errorf("invalid argument to car: %v", val)
	}
	if len(list.Elements) == 0 {
		return nil, fmt.Errorf("car of empty list")
	}
	return list.Elements[0], nil
}

// builtinCdr is built-in implementation of cdr list operation. It retrieves the rest elements of a list.
func builtinCdr(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of arguments to cdr")
	}
	val, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	list, ok := val.(*LispList)
	if !ok {
		return nil, fmt.Errorf("invalid argument to cdr: %v", val)
	}
	if len(list.Elements) == 0 {
		return &LispList{Elements: []LispValue{}}, nil
	}
	return &LispList{Elements: list.Elements[1:]}, nil
}

// builtinCons is built-in implementation of cons list operation. It add element to a list.
func builtinCons(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments to cons")
	}
	elem, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	val, err := Eval(env, args[1])
	if err != nil {
		return nil, err
	}
	list, ok := val.(*LispList)
	if !ok {
		return nil, fmt.Errorf("invalid argument to cons: %v", val)
	}
	return &LispList{Elements: append([]LispValue{elem}, list.Elements...)}, nil
}

// builtinLength is built-in implementation of length list operation. It retrieves the length of a list.
func builtinLength(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of arguments to length")
	}
	val, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	list, ok := val.(*LispList)
	if !ok {
		return nil, fmt.Errorf("invalid argument to length: %v", val)
	}
	return &LispNumber{Value: len(list.Elements)}, nil
}

// builtinAppend is built-in implementation of append list operation. It add a list to another list.
func builtinAppend(env Environment, args []LispValue) (LispValue, error) {
	var result []LispValue
	for _, arg := range args {
		val, err := Eval(env, arg)
		if err != nil {
			return nil, err
		}
		list, ok := val.(*LispList)
		if !ok {
			return nil, fmt.Errorf("invalid argument to append: %v", val)
		}
		result = append(result, list.Elements...)
	}
	return &LispList{Elements: result}, nil
}

// callFunction calls a user-defined function
func callFunction(env Environment, name string, args []LispValue) (LispValue, error) {
	fn, ok := env[name]
	if !ok {
		return nil, fmt.Errorf("undefined function: %s", name)
	}
	lambda, ok := fn.(*LispFunction)
	if !ok {
		return nil, fmt.Errorf("invalid function: %s", name)
	}
	if len(lambda.Params) != len(args) {
		return nil, fmt.Errorf("wrong number of arguments to %s", name)
	}
	localEnv := make(Environment)
	for key, value := range lambda.Env {
		localEnv[key] = value
	}
	for i, param := range lambda.Params {
		paramName, ok := param.(*LispAtom)
		if !ok {
			return nil, fmt.Errorf("invalid parameter name: %v", param)
		}
		argVal, err := Eval(env, args[i])
		if err != nil {
			return nil, err
		}
		localEnv[paramName.Value] = argVal
	}
	return Eval(localEnv, lambda.Body)
}
