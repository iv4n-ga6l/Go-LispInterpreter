"use client"

import { useCallback, useEffect, useState } from "react"
import { Loader } from "lucide-react"
import {
  ReactFlow,
  Node,
  Edge,
  Background,
  Controls,
  MiniMap,
  useNodesState,
  useEdgesState,
  MarkerType,
  Handle,
  Position,
} from "@xyflow/react"
import "@xyflow/react/dist/style.css"

interface TokenVisualization {
  type: string
  value: string
  position: number
  line: number
  column: number
}

interface ASTNodeVisualization {
  type: string
  value: string
  children?: ASTNodeVisualization[]
  nodeId: string
}

interface EvaluationStepVisualization {
  expression: string
  result: string
  nodeId: string
  parentId: string
  type: string
}

interface VisualizationResponse {
  success: boolean
  input: string
  tokens: TokenVisualization[]
  ast: ASTNodeVisualization
  evaluationSteps: EvaluationStepVisualization[]
  finalResult: string
  error?: string
}

interface LispVisualizerProps {
  code: string
  isVisualizing: boolean
  onVisualize?: () => void
}

const getAPIBase = () => {
  if (process.env.NEXT_PUBLIC_API_URL) {
    return process.env.NEXT_PUBLIC_API_URL
  }
  return "http://localhost:8080"
}

const API_BASE = getAPIBase()

// Custom node component
const CustomNode = ({ data }: { data: any }) => {
  return (
    <div className="relative">
      <Handle type="target" position={Position.Top} className="w-3 h-3 !bg-blue-400" />
      <div className="px-4 py-3 shadow-lg rounded-lg bg-card border-2 border-border min-w-[250px] max-w-[450px]">
        <div className="flex items-center gap-3 mb-3 pb-3 border-b border-border">
          <div className="w-10 h-10 bg-accent rounded-lg flex items-center justify-center text-xl">
            {data.icon}
          </div>
          <div className="flex-1 font-semibold text-sm text-foreground">{data.label}</div>
          <div className="w-2.5 h-2.5 rounded-full bg-green-500 animate-pulse" />
        </div>
        
        {data.description && (
          <div className="text-xs text-muted-foreground mb-3">{data.description}</div>
        )}
        
        {/* Token badges */}
        {data.tokenList && (
          <div className="flex flex-wrap gap-2 p-2 bg-accent/5 rounded-lg max-h-[250px] overflow-auto">
            {data.tokenList.map((token: TokenVisualization, i: number) => (
              <span
                key={i}
                className="inline-block px-3 py-1.5 rounded-md text-xs font-mono bg-accent border border-accent/60 text-foreground shadow-sm hover:shadow-md transition-shadow"
                title={`Type: ${token.type}\nLine: ${token.line}, Column: ${token.column}`}
              >
                {token.value}
              </span>
            ))}
          </div>
        )}
        
        {/* AST Tree */}
        {data.astTree && (
          <div className="p-3 bg-accent/5 rounded-lg text-xs font-mono max-h-[350px] overflow-auto">
            {renderASTTree(data.astTree, 0)}
          </div>
        )}
        
        {/* Evaluation step */}
        {data.expression && (
          <div className="p-3 bg-accent/5 rounded-lg border-l-4 border-accent">
            <div className="text-xs font-mono text-foreground mb-2">
              <span className="text-muted-foreground">Expression: </span>
              {data.expression}
            </div>
            {data.result && (
              <div className="text-sm font-mono text-green-400 font-bold">
                → {data.result}
              </div>
            )}
          </div>
        )}
        
        {/* Final result */}
        {data.finalResult && (
          <div className="p-4 bg-green-500/10 rounded-lg border-2 border-green-500/30">
            <div className="text-xl font-bold text-green-400 text-center font-mono">
              {data.finalResult}
            </div>
          </div>
        )}
      </div>
      <Handle type="source" position={Position.Bottom} className="w-3 h-3 !bg-blue-400" />
    </div>
  )
}

// Helper function to render AST tree with proper styling
const renderASTTree = (node: ASTNodeVisualization, level: number): React.ReactElement => {
  return (
    <div key={node.nodeId} className={level > 0 ? "ml-5 mt-2 pl-3 border-l-2 border-border/40" : "mb-2"}>
      <div className="inline-block px-2 py-1 bg-card border border-border rounded text-foreground mb-1">
        <span className="text-muted-foreground text-[10px] uppercase">{node.type}</span>
        <span className="mx-1">:</span>
        <span className="font-semibold">{node.value}</span>
      </div>
      {node.children && node.children.map((child) => renderASTTree(child, level + 1))}
    </div>
  )
}

const nodeTypes = {
  custom: CustomNode,
}

export default function LispVisualizer({ code, isVisualizing, onVisualize }: LispVisualizerProps) {
  const [nodes, setNodes, onNodesChange] = useNodesState<Node>([] as Node[])
  const [edges, setEdges, onEdgesChange] = useEdgesState<Edge>([] as Edge[])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!code) {
      setNodes([])
      setEdges([])
      setError(null)
    }
  }, [code])

  useEffect(() => {
    if (isVisualizing && code) {
      visualize()
    }
  }, [isVisualizing, code])

  const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms))

  const createTokenizationNode = (tokens: TokenVisualization[]): Node => {
    return {
      id: "token-node",
      type: "custom",
      position: { x: 50, y: 50 },
      data: {
        icon: "�",
        label: "Tokenization",
        description: `Found ${tokens.length} tokens`,
        tokenList: tokens,
        status: true,
      },
    }
  }

  const createParsingNode = (ast: ASTNodeVisualization): Node => {
    return {
      id: "ast-node",
      type: "custom",
      position: { x: 50, y: 350 },
      data: {
        icon: "🌳",
        label: "Abstract Syntax Tree",
        description: "Hierarchical structure built",
        astTree: ast,
        status: true,
      },
    }
  }

  const createEvaluationNodes = (steps: EvaluationStepVisualization[]): Node[] => {
    return steps.map((step, i) => ({
      id: step.nodeId,
      type: "custom",
      position: {
        x: 50,
        y: 650 + i * 200,
      },
      data: {
        icon: getIconForStepType(step.type),
        label: formatStepType(step.type),
        expression: step.expression,
        result: step.result,
        status: true,
      },
    }))
  }

  const createResultNode = (result: string): Node => {
    return {
      id: "result-node",
      type: "custom",
      position: { x: 50, y: 1000 }, // Will be adjusted based on eval nodes
      data: {
        icon: "✅",
        label: "Final Result",
        finalResult: result,
        status: true,
      },
    }
  }

  // Helper functions for evaluation node styling
  const getIconForStepType = (type: string): string => {
    const iconMap: Record<string, string> = {
      "function-call": "⚙️",
      "atom": "📌",
      "number": "🔢",
      "string": "📝",
      "list": "📋",
      "builtin": "🔧",
    }
    return iconMap[type] || "⚡"
  }

  const formatStepType = (type: string): string => {
    return type
      .split("-")
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
      .join(" ")
  }

  const visualize = async () => {
    if (!code.trim()) return

    setLoading(true)
    setError(null)
    setNodes([])
    setEdges([])

    try {
      const response = await fetch(`${API_BASE}/api/visualize`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ code }),
      })

      const data: VisualizationResponse = await response.json()

      if (!data.success) {
        setError(data.error || "Unknown error occurred")
        setLoading(false)
        return
      }

      // Create all nodes
      const newNodes: Node[] = []
      const newEdges: Edge[] = []

      // Stage 1: Tokenization
      await sleep(300)
      const tokenNode = createTokenizationNode(data.tokens)
      newNodes.push(tokenNode)
      setNodes([...newNodes])

      // Stage 2: Parsing (AST)
      await sleep(800)
      const astNode = createParsingNode(data.ast)
      newNodes.push(astNode)
      
      // Connection: Token -> AST
      newEdges.push({
        id: "e-token-ast",
        source: "token-node",
        target: "ast-node",
        animated: true,
        type: "smoothstep",
        markerEnd: { type: MarkerType.ArrowClosed, color: "#60a5fa" },
        style: { stroke: "#60a5fa", strokeWidth: 2 },
      })
      
      setNodes([...newNodes])
      setEdges([...newEdges])

      // Stage 3: Evaluation
      await sleep(800)
      const evalNodes = createEvaluationNodes(data.evaluationSteps)
      newNodes.push(...evalNodes)
      
      // Connection: AST -> First Evaluation
      if (evalNodes.length > 0) {
        newEdges.push({
          id: "e-ast-eval0",
          source: "ast-node",
          target: evalNodes[0].id,
          animated: true,
          type: "smoothstep",
          markerEnd: { type: MarkerType.ArrowClosed, color: "#60a5fa" },
          style: { stroke: "#60a5fa", strokeWidth: 2 },
        })

        // Connections between evaluation steps
        for (let i = 0; i < evalNodes.length - 1; i++) {
          newEdges.push({
            id: `e-eval${i}-eval${i + 1}`,
            source: evalNodes[i].id,
            target: evalNodes[i + 1].id,
            animated: true,
            type: "smoothstep",
            markerEnd: { type: MarkerType.ArrowClosed, color: "#60a5fa" },
            style: { stroke: "#60a5fa", strokeWidth: 2 },
          })
        }
      }
      
      setNodes([...newNodes])
      setEdges([...newEdges])

      // Final Result
      await sleep(800)
      
      // Adjust result node position based on last eval node
      const lastEvalY = evalNodes.length > 0 ? 650 + (evalNodes.length - 1) * 200 + 200 : 650
      const resultNode = createResultNode(data.finalResult)
      resultNode.position.y = lastEvalY
      newNodes.push(resultNode)
      
      // Connection: Last Evaluation -> Result (or AST -> Result if no eval steps)
      if (evalNodes.length > 0) {
        newEdges.push({
          id: "e-eval-result",
          source: evalNodes[evalNodes.length - 1].id,
          target: "result-node",
          animated: true,
          type: "smoothstep",
          markerEnd: { type: MarkerType.ArrowClosed, color: "#4ade80" },
          style: { stroke: "#4ade80", strokeWidth: 2 },
        })
      } else {
        newEdges.push({
          id: "e-ast-result",
          source: "ast-node",
          target: "result-node",
          animated: true,
          type: "smoothstep",
          markerEnd: { type: MarkerType.ArrowClosed, color: "#4ade80" },
          style: { stroke: "#4ade80", strokeWidth: 2 },
        })
      }

      setNodes([...newNodes])
      setEdges([...newEdges])

      setLoading(false)
      if (onVisualize) onVisualize()
    } catch (err) {
      setError(`Failed to connect to server: ${err instanceof Error ? err.message : String(err)}`)
      setLoading(false)
    }
  }

  const renderWelcome = () => (
    <div className="h-full flex items-center justify-center">
      <div className="text-center">
        <div className="w-20 h-20 mx-auto mb-6 rounded-lg bg-secondary flex items-center justify-center">
          <svg className="w-10 h-10 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={1.5}
              d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
            />
          </svg>
        </div>
        <h3 className="text-lg font-bold mb-2">Ready to Visualize</h3>
        <p className="text-sm text-muted-foreground">
          Enter a Lisp expression and click Visualize
          <br />
          to see the evaluation process
        </p>
      </div>
    </div>
  )

  const renderLoading = () => (
    <div className="h-full flex items-center justify-center">
      <div className="text-center">
        <Loader className="w-8 h-8 animate-spin mx-auto mb-4 text-accent" />
        <p className="text-sm text-muted-foreground">Processing your code...</p>
      </div>
    </div>
  )

  const renderError = () => (
    <div className="h-full flex items-center justify-center p-6">
      <div className="bg-destructive/10 border border-destructive/30 rounded-lg p-6 max-w-md text-center">
        <svg className="w-8 h-8 mx-auto mb-3 text-destructive" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
          />
        </svg>
        <p className="font-semibold text-destructive mb-2">Error</p>
        <p className="text-sm text-foreground font-mono">{error}</p>
      </div>
    </div>
  )

  if (loading) {
    return renderLoading()
  }

  if (error) {
    return renderError()
  }

  if (!code && nodes.length === 0) {
    return renderWelcome()
  }

  return (
    <div className="flex-1 bg-background relative h-full">
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        nodeTypes={nodeTypes}
        defaultEdgeOptions={{
          type: "smoothstep",
          animated: true,
          style: { stroke: "#60a5fa", strokeWidth: 3 },
        }}
        fitView
        minZoom={0.5}
        maxZoom={1.5}
        defaultViewport={{ x: 0, y: 0, zoom: 0.8 }}
        className="bg-background"
      >
        <Background color="#333" gap={16} />
        <Controls />
        <MiniMap nodeColor="#4a5568" maskColor="rgb(12, 15, 10, 0.8)" />
      </ReactFlow>
    </div>
  )
}
