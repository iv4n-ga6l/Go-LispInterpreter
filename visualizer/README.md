# Lisp Interpreter Visualizer (Next.js)

A modern, React-based web visualizer for the Go Lisp Interpreter. Built with Next.js, TypeScript, and shadcn/ui components.

## Features

- 🎨 **Beautiful UI**: Built with shadcn/ui and Tailwind CSS
- 🌙 **Dark Mode**: Elegant dark theme optimized for code visualization
- ⚡ **Real-time Visualization**: See tokenization, parsing, and evaluation steps
- 🎭 **Animated Nodes**: Smooth animations showing the interpretation flow
- 📱 **Responsive Design**: Works on desktop and mobile
- 🔍 **Type Safe**: Full TypeScript support
- 🚀 **Fast**: Optimized React components with Next.js

## Prerequisites

- Node.js 18+ (or 20+)
- npm, pnpm, or yarn
- Go backend server running (see main README)

## Installation

```bash
# Install dependencies
npm install
# or
pnpm install
# or
yarn install
```

## Configuration

Create a `.env.local` file in the visualizer directory:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

You can copy from the example:
```bash
cp .env.example .env.local
```

## Running

### Development Mode

```bash
# Make sure the Go backend is running first
# In the root directory: go run . server

# Then start the Next.js dev server
npm run dev
# or
pnpm dev
# or
yarn dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

### Production Build

```bash
# Build the application for static export
npm run build
# or
pnpm build

# The output will be in the 'out' directory
# The Go server will automatically serve this directory
```

To deploy with the Go server:
```bash
# 1. Build Next.js
cd visualizer
npm run build

# 2. Start Go server (it will serve the Next.js build)
cd ..
go run . server

# 3. Open http://localhost:8080
```

## Usage

1. **Start the Go Backend**: 
   ```bash
   cd ..
   go run . server
   ```

2. **Start the Next.js Frontend**:
   ```bash
   npm run dev
   ```

3. **Use the Visualizer**:
   - Enter a Lisp expression in the textarea (left panel)
   - Click "Visualize" to see the interpretation process
   - Watch the animated nodes showing tokenization → parsing → evaluation → result
   - Try the example buttons for quick testing

## Example Expressions

```lisp
(+ 1 2)
(* (+ 2 3) (- 10 5))
(if (> 5 3) "yes" "no")
(defun square (x) (* x x))
(let ((a 10) (b 20)) (+ a b))
```

## Project Structure

```
visualizer/
├── app/                  # Next.js app directory
│   ├── page.tsx         # Main page with input panel
│   └── layout.tsx       # Root layout with providers
├── components/
│   ├── lisp-visualizer.tsx    # Main visualization component
│   ├── theme-provider.tsx     # Dark mode support
│   └── ui/              # shadcn/ui components
├── lib/                 # Utility functions
├── styles/              # Global styles
├── .env.example         # Environment variable template
└── package.json         # Dependencies and scripts
```

## Key Components

### LispVisualizer

The main visualization component that:
- Fetches data from the Go backend API
- Renders animated nodes for each interpretation stage
- Draws SVG connections between nodes
- Handles loading and error states

### Page

The main page component that:
- Provides the code input interface
- Manages visualization state
- Includes example expressions
- Handles user interactions

## API Integration

The visualizer communicates with the Go backend via REST API:

```typescript
POST http://localhost:8080/api/visualize
Content-Type: application/json

{
  "code": "(+ 1 2)"
}
```

Response includes:
- Tokenization data
- Abstract Syntax Tree (AST)
- Evaluation steps
- Final result

## Technologies

- **Framework**: Next.js 16.0
- **Language**: TypeScript
- **UI Components**: shadcn/ui (Radix UI + Tailwind)
- **Styling**: Tailwind CSS
- **Icons**: Lucide React
- **Theme**: next-themes

## Development

### Adding New Components

```bash
# Use the shadcn CLI to add components
npx shadcn@latest add [component-name]
```

### Linting

```bash
npm run lint
```

## Troubleshooting

### CORS Issues

If you encounter CORS errors, make sure:
1. The Go backend is running with `go run . server`
2. The backend has CORS enabled (it's enabled by default in visualizer.go)
3. Your `.env.local` has the correct API URL

### Port Conflicts

If port 3000 is already in use:
```bash
PORT=3001 npm run dev
```

### Connection Refused

Make sure the Go backend is running on port 8080:
```bash
# In the root directory
go run . server
```

## Contributing

This visualizer is part of the Go Lisp Interpreter project. See the main README for contribution guidelines.

## License

Same license as the main project.
