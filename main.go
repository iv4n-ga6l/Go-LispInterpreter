package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/c-bata/go-prompt"
)

// ExecuteRequest represents the request structure for code execution
type ExecuteRequest struct {
	Code string `json:"code"`
}

// ExecuteResponse represents the response structure for code execution
type ExecuteResponse struct {
	Success       bool     `json:"success"`
	Results       []string `json:"results,omitempty"`
	Error         string   `json:"error,omitempty"`
	ExecutionTime string   `json:"execution_time"`
}

// FileExecuteRequest represents the request structure for file execution
type FileExecuteRequest struct {
	Content  string `json:"content"`
	Filename string `json:"filename,omitempty"`
}

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

// executeCode executes Lisp code and returns the results
func executeCode(code string) ExecuteResponse {
	start := time.Now()

	defer func() {
		if r := recover(); r != nil {
			// Handle panic and return error response
		}
	}()

	tokens := Tokenize(code)
	expr, _, err := Parse(tokens)
	if err != nil {
		return ExecuteResponse{
			Success:       false,
			Error:         err.Error(),
			ExecutionTime: time.Since(start).String(),
		}
	}

	var results []string
	if list, ok := expr.(*LispList); ok {
		evalResults, err := evalMultipleExpressions(env, list.Elements)
		if err != nil {
			return ExecuteResponse{
				Success:       false,
				Error:         err.Error(),
				ExecutionTime: time.Since(start).String(),
			}
		}
		for _, result := range evalResults {
			results = append(results, result.String())
		}
	} else {
		result, err := Eval(env, expr)
		if err != nil {
			return ExecuteResponse{
				Success:       false,
				Error:         err.Error(),
				ExecutionTime: time.Since(start).String(),
			}
		}
		results = append(results, result.String())
	}

	return ExecuteResponse{
		Success:       true,
		Results:       results,
		ExecutionTime: time.Since(start).String(),
	}
}

// HTTP Handlers

// enableCORS adds CORS headers to the response
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// handleExecute handles code execution requests
func handleExecute(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ExecuteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	response := executeCode(req.Code)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleFileExecute handles file execution requests
func handleFileExecute(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req FileExecuteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	response := executeCode(req.Content)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleReset resets the environment
func handleReset(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	env = initEnvironment()

	response := map[string]interface{}{
		"success": true,
		"message": "Environment reset successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// serveStatic serves the HTML file
func serveStatic(w http.ResponseWriter, r *http.Request) {
	// Serve the HTML IDE file
	http.ServeFile(w, r, "ide.html")
}

func main() {
	env = initEnvironment()

	// Check if running in server mode or traditional mode
	if len(os.Args) > 1 && os.Args[1] == "server" {
		// HTTP Server mode
		http.HandleFunc("/", serveStatic)
		http.HandleFunc("/api/execute", handleExecute)
		http.HandleFunc("/api/file-execute", handleFileExecute)
		http.HandleFunc("/api/reset", handleReset)

		port := "8080"
		if len(os.Args) > 2 {
			port = os.Args[2]
		}

		fmt.Printf("CCLisp IDE Server starting on http://localhost:%s\n", port)
		fmt.Printf("API Endpoints:\n")
		fmt.Printf("  POST /api/execute - Execute Lisp code\n")
		fmt.Printf("  POST /api/file-execute - Execute Lisp file content\n")
		fmt.Printf("  POST /api/reset - Reset environment\n")
		fmt.Printf("\nPress Ctrl+C to stop the server\n\n")

		log.Fatal(http.ListenAndServe(":"+port, nil))
	} else if len(os.Args) > 1 && os.Args[1] != "server" {
		// File execution mode (original functionality)
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
		// REPL mode (original functionality)
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
