package main

import (
	"fmt"
	"os"
	"time"

	"github.com/c-bata/go-prompt"
)

// evalMultipleExpressions evaluates multiple expressions and returns the results
func evalMultipleExpressions(env Environment, expressions []LispValue) ([]LispValue, error) {
	results := make([]LispValue, 0, len(expressions))
	for _, expr := range expressions {
		result, err := Eval(env, expr)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

// Environment represents a symbol table
var env Environment

// completer returns suggestions for the prompt
func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}

	for key, value := range builtins {
		s = append(s, prompt.Suggest{Text: key, Description: value})
	}

	// Add defined symbols from the environment
	for symbol := range env {
		s = append(s, prompt.Suggest{Text: symbol, Description: "Defined symbol"})
	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

// executor reads the input, tokenizes it, parses it, and evaluates it
func executor(input string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	tokens := Tokenize(input)
	expr, _, err := Parse(tokens)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if list, ok := expr.(*LispList); ok {
		results, err := evalMultipleExpressions(env, list.Elements)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			for _, result := range results {
				fmt.Println(result)
			}
		}
	} else {
		result, err := Eval(env, expr)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println(result)
		}
	}
}

// initEnvironment initializes the environment with predefined symbols
func initEnvironment() Environment {
	env := make(Environment)
	env[T] = &LispBoolean{Value: true}
	env[NIL] = &LispNil{}
	env[TRUE] = &LispBoolean{Value: true}
	env[FALSE] = &LispBoolean{Value: false}
	return env
}

// readFile reads the content of a file and returns it as a string
func readFile(filepath string) (string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func main() {
	env = initEnvironment()

	if len(os.Args) > 1 {
		// File execution mode
		filepath := os.Args[1]
		content, err := readFile(filepath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		start := time.Now()
		tokens := Tokenize(content)
		expr, _, err := Parse(tokens)
		if err != nil {
			fmt.Println("Error parsing file:", err)
			return
		}
		results, err := evalMultipleExpressions(env, expr.(*LispList).Elements)
		if err != nil {
			fmt.Println("Error evaluating file:", err)
			return
		}
		elapsed := time.Since(start)

		for _, result := range results {
			fmt.Println(result)
		}
		fmt.Printf("\n")
		fmt.Printf("Execution time: %v\n", elapsed)
	} else {
		// REPL mode
		p := prompt.New(
			func(input string) {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("Recovered from panic:", r)
					}
				}()
				executor(input)
			},
			completer,
			prompt.OptionPrefix("cclisp> "),
			prompt.OptionTitle("CCLisp REPL"),
			prompt.OptionAddKeyBind(prompt.KeyBind{
				Key: prompt.ControlC,
				Fn: func(buf *prompt.Buffer) {
					fmt.Println("Exiting REPL...")
					os.Exit(0)
				},
			}),
		)
		p.Run()
	}
}
