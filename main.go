package main

import (
	"fmt"
	"log"
	"net/http"
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

// serveVisualizer serves the Next.js visualizer or provides API info
func serveVisualizer(w http.ResponseWriter, r *http.Request) {
	// Serve static files from Next.js build if available
	nextBuildPath := "./visualizer/out"
	if _, err := os.Stat(nextBuildPath); err == nil {
		// Next.js static export exists
		fs := http.FileServer(http.Dir(nextBuildPath))
		fs.ServeHTTP(w, r)
		return
	}

	// If Next.js build not found, show API info page
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Lisp Interpreter API</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            max-width: 800px;
            margin: 50px auto;
            padding: 20px;
            background: #0C0F0A;
            color: #ffffff;
        }
        .card {
            background: #1a1f2e;
            border: 2px solid #2e3750;
            border-radius: 12px;
            padding: 30px;
            margin: 20px 0;
        }
        h1 { color: #4ade80; margin-top: 0; }
        h2 { color: #60a5fa; }
        code {
            background: #0D1D2C;
            padding: 2px 8px;
            border-radius: 4px;
            font-family: 'Courier New', monospace;
        }
        pre {
            background: #0D1D2C;
            padding: 15px;
            border-radius: 8px;
            overflow-x: auto;
        }
        .status { color: #4ade80; }
        a { color: #60a5fa; text-decoration: none; }
        a:hover { text-decoration: underline; }
    </style>
</head>
<body>
    <div class="card">
        <h1>🚀 Lisp Interpreter API Server</h1>
        <p class="status">✓ Server is running</p>
        <p>The API server is active and ready to process Lisp expressions.</p>
    </div>

    <div class="card">
        <h2>📊 Web Visualizer</h2>
        <p>To use the web visualizer, you need to run the Next.js frontend:</p>
        <pre>cd visualizer
npm install
npm run dev</pre>
        <p>Then open <a href="http://localhost:3000" target="_blank">http://localhost:3000</a></p>
    </div>

    <div class="card">
        <h2>🔌 API Endpoint</h2>
        <p><strong>POST</strong> <code>/api/visualize</code></p>
        <pre>{
  "code": "(+ 1 2)"
}</pre>
        <p><strong>Response:</strong></p>
        <pre>{
  "success": true,
  "input": "(+ 1 2)",
  "tokens": [...],
  "ast": {...},
  "evaluationSteps": [...],
  "finalResult": "3"
}</pre>
    </div>

    <div class="card">
        <h2>📖 Documentation</h2>
        <p>For more information, see the <a href="https://github.com/IvanGael/Go-LispInterpreter">GitHub repository</a></p>
    </div>
</body>
</html>`
	w.Write([]byte(html))
}

func main() {
	env = initEnvironment()

	// Check if running in server mode
	if len(os.Args) > 1 && os.Args[1] == "server" {
		// HTTP Server mode for visualizer API
		http.HandleFunc("/", serveVisualizer)
		http.HandleFunc("/api/visualize", handleVisualize)

		port := "8080"
		if len(os.Args) > 2 {
			port = os.Args[2]
		}

		fmt.Printf("🚀 Lisp Interpreter API Server starting on http://localhost:%s\n", port)
		fmt.Printf("📊 API Endpoint: POST /api/visualize\n")
		fmt.Printf("🌐 Web UI: Run 'cd visualizer && npm run dev' then open http://localhost:3000\n")
		fmt.Printf("\n💡 Press Ctrl+C to stop the server\n\n")

		log.Fatal(http.ListenAndServe(":"+port, nil))
	} else if len(os.Args) > 1 {
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
