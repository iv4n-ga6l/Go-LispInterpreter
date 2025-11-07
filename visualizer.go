package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// VisualizationStep represents a step in the interpretation process
type VisualizationStep struct {
	Stage       string      `json:"stage"`       // "tokenization", "parsing", "evaluation"
	Description string      `json:"description"` // Human-readable description
	Data        interface{} `json:"data"`        // Stage-specific data
	NodeID      string      `json:"nodeId"`      // Unique identifier for this node
	ParentID    string      `json:"parentId"`    // Parent node ID for tree structure
}

// TokenVisualization represents a token with visual metadata
type TokenVisualization struct {
	Type     string `json:"type"`
	Value    string `json:"value"`
	Position int    `json:"position"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
}

// ASTNodeVisualization represents an AST node with visual metadata
type ASTNodeVisualization struct {
	Type     string                  `json:"type"`
	Value    string                  `json:"value"`
	Children []*ASTNodeVisualization `json:"children,omitempty"`
	NodeID   string                  `json:"nodeId"`
}

// EvaluationStepVisualization represents an evaluation step
type EvaluationStepVisualization struct {
	Expression string `json:"expression"`
	Result     string `json:"result"`
	NodeID     string `json:"nodeId"`
	ParentID   string `json:"parentId"`
	Type       string `json:"type"` // "function-call", "atom", "number", "string", etc.
}

// VisualizationResponse represents the complete visualization data
type VisualizationResponse struct {
	Success         bool                          `json:"success"`
	Input           string                        `json:"input"`
	Tokens          []TokenVisualization          `json:"tokens"`
	AST             *ASTNodeVisualization         `json:"ast"`
	EvaluationSteps []EvaluationStepVisualization `json:"evaluationSteps"`
	FinalResult     string                        `json:"finalResult"`
	Error           string                        `json:"error,omitempty"`
	AllSteps        []VisualizationStep           `json:"allSteps"`
}

var nodeCounter = 0

func generateNodeID(prefix string) string {
	nodeCounter++
	return fmt.Sprintf("%s-%d", prefix, nodeCounter)
}

// visualizeTokens converts tokens to visualization format
func visualizeTokens(tokens []Token) []TokenVisualization {
	result := make([]TokenVisualization, len(tokens))
	for i, token := range tokens {
		result[i] = TokenVisualization{
			Type:     token.Type,
			Value:    token.Value,
			Position: i,
			Line:     token.Line,
			Column:   token.Column,
		}
	}
	return result
}

// visualizeAST converts AST to visualization format
func visualizeAST(value LispValue, nodeID string) *ASTNodeVisualization {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case *LispAtom:
		return &ASTNodeVisualization{
			Type:   "atom",
			Value:  v.Value,
			NodeID: nodeID,
		}
	case *LispNumber:
		return &ASTNodeVisualization{
			Type:   "number",
			Value:  v.String(),
			NodeID: nodeID,
		}
	case *LispFloat:
		return &ASTNodeVisualization{
			Type:   "float",
			Value:  v.String(),
			NodeID: nodeID,
		}
	case *LispString:
		return &ASTNodeVisualization{
			Type:   "string",
			Value:  fmt.Sprintf(`"%s"`, v.Value),
			NodeID: nodeID,
		}
	case *LispBoolean:
		return &ASTNodeVisualization{
			Type:   "boolean",
			Value:  v.String(),
			NodeID: nodeID,
		}
	case *LispNil:
		return &ASTNodeVisualization{
			Type:   "nil",
			Value:  "nil",
			NodeID: nodeID,
		}
	case *LispList:
		node := &ASTNodeVisualization{
			Type:     "list",
			Value:    "",
			NodeID:   nodeID,
			Children: make([]*ASTNodeVisualization, 0),
		}
		for i, elem := range v.Elements {
			childID := fmt.Sprintf("%s-%d", nodeID, i)
			child := visualizeAST(elem, childID)
			if child != nil {
				node.Children = append(node.Children, child)
			}
		}
		return node
	case *LispFunction:
		return &ASTNodeVisualization{
			Type:   "function",
			Value:  "λ",
			NodeID: nodeID,
		}
	default:
		return &ASTNodeVisualization{
			Type:   "unknown",
			Value:  fmt.Sprintf("%v", v),
			NodeID: nodeID,
		}
	}
}

// captureEvaluation performs evaluation while capturing steps
func captureEvaluation(env Environment, expr LispValue, parentID string) ([]EvaluationStepVisualization, LispValue, error) {
	steps := make([]EvaluationStepVisualization, 0)
	nodeID := generateNodeID("eval")

	// Create step for this evaluation
	step := EvaluationStepVisualization{
		Expression: expr.String(),
		NodeID:     nodeID,
		ParentID:   parentID,
	}

	// Determine type and evaluate
	switch v := expr.(type) {
	case *LispAtom:
		step.Type = "atom-lookup"
		result, err := Eval(env, expr)
		if err != nil {
			return steps, nil, err
		}
		step.Result = result.String()
		steps = append(steps, step)
		return steps, result, nil

	case *LispNumber, *LispFloat, *LispString, *LispBoolean, *LispNil:
		step.Type = "literal"
		step.Result = v.String()
		steps = append(steps, step)
		return steps, expr, nil

	case *LispList:
		if len(v.Elements) == 0 {
			step.Type = "empty-list"
			step.Result = "nil"
			steps = append(steps, step)
			return steps, &LispNil{}, nil
		}

		step.Type = "function-call"
		steps = append(steps, step)

		// Evaluate each element
		for _, elem := range v.Elements {
			childSteps, _, err := captureEvaluation(env, elem, nodeID)
			if err == nil {
				steps = append(steps, childSteps...)
			}
		}

		// Perform actual evaluation
		result, err := Eval(env, expr)
		if err != nil {
			return steps, nil, err
		}

		// Add result step
		resultStep := EvaluationStepVisualization{
			Expression: expr.String(),
			Result:     result.String(),
			NodeID:     generateNodeID("result"),
			ParentID:   nodeID,
			Type:       "result",
		}
		steps = append(steps, resultStep)
		return steps, result, nil

	default:
		result, err := Eval(env, expr)
		if err != nil {
			return steps, nil, err
		}
		step.Result = result.String()
		steps = append(steps, step)
		return steps, result, nil
	}
}

// handleVisualize handles visualization requests
func handleVisualize(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Reset node counter
	nodeCounter = 0

	response := VisualizationResponse{
		Success: true,
		Input:   req.Code,
	}

	// Step 1: Tokenization
	tokens := Tokenize(req.Code)
	response.Tokens = visualizeTokens(tokens)

	allSteps := make([]VisualizationStep, 0)

	// Add tokenization steps
	for i, token := range response.Tokens {
		allSteps = append(allSteps, VisualizationStep{
			Stage:       "tokenization",
			Description: fmt.Sprintf("Token %d: %s '%s'", i, token.Type, token.Value),
			Data:        token,
			NodeID:      generateNodeID("token"),
			ParentID:    "",
		})
	}

	// Step 2: Parsing
	expr, _, err := Parse(tokens)
	if err != nil {
		response.Success = false
		response.Error = err.Error()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	astNodeID := generateNodeID("ast")
	response.AST = visualizeAST(expr, astNodeID)

	allSteps = append(allSteps, VisualizationStep{
		Stage:       "parsing",
		Description: "Abstract Syntax Tree constructed",
		Data:        response.AST,
		NodeID:      astNodeID,
		ParentID:    "",
	})

	// Step 3: Evaluation
	evalEnv := initEnvironment()
	evaluationSteps, result, err := captureEvaluation(evalEnv, expr, "")
	if err != nil {
		response.Success = false
		response.Error = err.Error()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response.EvaluationSteps = evaluationSteps
	response.FinalResult = result.String()

	// Add evaluation steps
	for _, evalStep := range evaluationSteps {
		allSteps = append(allSteps, VisualizationStep{
			Stage:       "evaluation",
			Description: fmt.Sprintf("Evaluating: %s", evalStep.Expression),
			Data:        evalStep,
			NodeID:      evalStep.NodeID,
			ParentID:    evalStep.ParentID,
		})
	}

	response.AllSteps = allSteps

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// enableCORS adds CORS headers to the response
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
