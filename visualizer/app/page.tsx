"use client"

import { useState } from "react"
import { Play, RotateCcw, Zap, Code } from "lucide-react"
import LispVisualizer from "@/components/lisp-visualizer"
import { Button } from "@/components/ui/button"
import { Textarea } from "@/components/ui/textarea"

export default function Home() {
  const [code, setCode] = useState("(+ (* 2 3) (- 10 5))")
  const [isVisualizing, setIsVisualizing] = useState(false)
  const [visualizeTrigger, setVisualizeTrigger] = useState(0)

  const examples = ["(+ 1 2)", "(* (+ 2 3) (- 10 5))", '(if (> 5 3) "yes" "no")', "(defun square (x) (* x x))"]

  const handleVisualize = async () => {
    if (!code.trim()) return
    setIsVisualizing(true)
    setVisualizeTrigger((prev) => prev + 1)
  }

  const handleClear = () => {
    setCode("")
    setIsVisualizing(false)
  }

  const handleExample = () => {
    setCode(examples[Math.floor(Math.random() * examples.length)])
  }

  const onVisualizeComplete = () => {
    setIsVisualizing(false)
  }

  return (
    <div className="h-screen flex flex-col bg-background text-foreground">
      {/* Header */}
      <header className="border-b border-border bg-card shadow-sm">
        <div className="max-w-7xl mx-auto px-6 py-4 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-lg bg-accent flex items-center justify-center">
              <Code className="w-6 h-6 text-accent-foreground" />
            </div>
            <div>
              <h1 className="text-2xl font-bold tracking-tight">Lisp Viz</h1>
              <p className="text-sm text-muted-foreground">Quickly visualize Lisp expressions</p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-2 h-2 rounded-full bg-green-500 animate-pulse-ring"></div>
            <span className="text-xs text-muted-foreground font-medium">Ready</span>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <div className="flex-1 flex overflow-hidden">
        {/* Left Panel: Input */}
        <div className="w-1/3 border-r border-border flex flex-col bg-card">
          <div className="p-6 flex-1 flex flex-col gap-4 overflow-y-auto">
            <div className="space-y-2">
              <label className="text-sm font-semibold flex items-center gap-2">
                <Code className="w-4 h-4 text-muted-foreground" />
                Lisp Expression
              </label>
              <Textarea
                value={code}
                onChange={(e) => setCode(e.target.value)}
                placeholder="Enter your Lisp expression here...&#10;&#10;Example: (+ 1 2)"
                className="font-mono h-48 text-sm resize-none"
                spellCheck={false}
              />
            </div>

            <div className="flex gap-2">
              <Button
                onClick={handleVisualize}
                disabled={!code.trim()}
                className="flex-1 bg-accent hover:bg-accent/90 text-accent-foreground"
              >
                <Play className="w-4 h-4 mr-2" />
                Visualize
              </Button>
              <Button onClick={handleClear} variant="outline" size="icon" title="Clear">
                <RotateCcw className="w-4 h-4" />
              </Button>
              <Button onClick={handleExample} variant="outline" size="icon" title="Random example">
                <Zap className="w-4 h-4" />
              </Button>
            </div>

            <div className="pt-4 border-t border-border">
              <p className="text-xs font-semibold text-muted-foreground mb-3 uppercase tracking-wide">Quick Examples</p>
              <div className="space-y-2">
                {examples.map((example, i) => (
                  <button
                    key={i}
                    onClick={() => setCode(example)}
                    className="w-full text-left px-3 py-2 rounded-lg bg-secondary hover:bg-secondary/80 transition-colors text-sm font-mono text-foreground"
                  >
                    {example}
                  </button>
                ))}
              </div>
            </div>
          </div>
        </div>

        {/* Right Panel: Visualization */}
        <div className="flex-1 flex flex-col overflow-hidden">
          <LispVisualizer 
            code={code} 
            isVisualizing={isVisualizing} 
            onVisualize={onVisualizeComplete}
            key={visualizeTrigger}
          />
        </div>
      </div>
    </div>
  )
}
