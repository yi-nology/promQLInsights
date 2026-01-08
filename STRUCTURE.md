# Project Structure Overview

```
insights/
├── frontend/                          # Vue + Vite Frontend Application
│   ├── src/
│   │   ├── components/               # Vue Components
│   │   │   ├── ASTTree.vue          # AST Tree Visualization
│   │   │   ├── ASTNode.vue          # Recursive AST Node Component
│   │   │   ├── Dependencies.vue     # Dependencies Display Component
│   │   │   ├── Issues.vue           # Issues/Warnings Display
│   │   │   └── QueryResults.vue     # Query Results (Chart/Table)
│   │   ├── views/                   # Page Views
│   │   │   └── PromQLWorkbench.vue # Main Workbench Page
│   │   ├── router/                  # Vue Router Configuration
│   │   │   └── index.ts
│   │   ├── types/                   # TypeScript Type Definitions
│   │   │   └── api.ts               # API Request/Response Types
│   │   ├── style.css                # Global Styles
│   │   ├── main.ts                  # Application Entry Point
│   │   └── App.vue                  # Root Component
│   ├── index.html                   # HTML Template
│   ├── package.json                 # NPM Dependencies
│   ├── vite.config.ts               # Vite Configuration
│   ├── tsconfig.json                # TypeScript Configuration
│   ├── tsconfig.node.json           # TypeScript Node Config
│   ├── .eslintrc.js                 # ESLint Configuration
│   ├── .prettierrc                  # Prettier Configuration
│   ├── Dockerfile                   # Docker Build File
│   └── nginx.conf                   # Nginx Configuration
│
├── backend/                         # Cloudwego Backend Application
│   ├── cmd/
│   │   └── api/
│   │       └── main.go             # Application Entry Point
│   ├── internal/
│   │   ├── handler/                 # HTTP Request Handlers
│   │   │   └── promql.go           # PromQL Handler
│   │   ├── middleware/              # HTTP Middleware
│   │   │   └── cors.go             # CORS Middleware
│   │   ├── dto/                    # Data Transfer Objects
│   │   │   └── api.go              # API DTOs
│   │   └── promql/                 # PromQL Processing
│   │       ├── parser/             # PromQL Parser Module
│   │       │   └── parser.go       # Parser Implementation
│   │       └── analyzer/           # PromQL Analyzer Module
│   │           └── analyzer.go     # Analyzer Implementation
│   ├── go.mod                      # Go Module Definition
│   ├── go.sum                      # Go Module Checksums
│   └── Dockerfile                  # Docker Build File
│
├── .gitignore                       # Git Ignore Rules
├── Makefile                         # Build Automation
├── docker-compose.yml               # Docker Compose Configuration
├── README.md                        # Project Documentation
├── ARCHITECTURE.md                  # Architecture Documentation
└── DEVELOPMENT.md                   # Development Guide
```

## Component Hierarchy

### Frontend Component Tree

```
App.vue
└── PromQLWorkbench.vue
    ├── Header
    │   ├── Title
    │   ├── Datasource Selector
    │   ├── Time Range Picker
    │   ├── Step Input
    │   └── Run Button
    ├── Editor Panel
    │   ├── Toolbar
    │   │   ├── Format Button
    │   │   └── Parse Button
    │   └── Monaco Editor
    └── Analysis Panel
        ├── Tabs
        │   ├── AST Tab
        │   ├── Dependencies Tab
        │   ├── Issues Tab
        │   └── Results Tab
        └── Tab Content
            ├── ASTTree.vue
            │   └── ASTNode.vue (recursive)
            ├── Dependencies.vue
            ├── Issues.vue
            └── QueryResults.vue
                ├── Chart View (ECharts)
                └── Table View
```

### Backend Module Hierarchy

```
main.go
├── Hertz Server
│   ├── CORS Middleware
│   └── Routes
│       ├── POST /api/promql/parse
│       └── POST /api/query/range
└── Handlers
    ├── ParsePromQL
    │   ├── Parser Module
    │   │   ├── Parse PromQL
    │   │   ├── Format Query
    │   │   └── Convert AST to JSON
    │   └── Analyzer Module
    │       ├── Extract Dependencies
    │       └── Check Issues
    └── QueryRange
        └── (Future: Proxy to Prometheus)
```

## Data Flow

### Request Flow

```
User Input
  ↓
Frontend Component
  ↓
API Call (Axios)
  ↓
Backend Handler
  ↓
Business Logic (Parser/Analyzer)
  ↓
Response
  ↓
Frontend UI Update
```

### AST Processing Flow

```
PromQL String
  ↓
Prometheus Parser
  ↓
Go AST
  ↓
AST to JSON Converter
  ↓
JSON DTO
  ↓
Frontend
  ↓
ASTTree Component
  ↓
Visualized Tree
```

## Key Files and Their Purposes

### Frontend Key Files

| File | Purpose |
|------|---------|
| [PromQLWorkbench.vue](file:///Users/zhangyi/insights/frontend/src/views/PromQLWorkbench.vue) | Main workbench page with editor and analysis panels |
| [ASTNode.vue](file:///Users/zhangyi/insights/frontend/src/components/ASTNode.vue) | Recursive component for AST tree visualization |
| [Dependencies.vue](file:///Users/zhangyi/insights/frontend/src/components/Dependencies.vue) | Displays extracted metrics, labels, functions |
| [Issues.vue](file:///Users/zhangyi/insights/frontend/src/components/Issues.vue) | Shows analysis issues and warnings |
| [QueryResults.vue](file:///Users/zhangyi/insights/frontend/src/components/QueryResults.vue) | Displays query results in chart/table format |
| [api.ts](file:///Users/zhangyi/insights/frontend/src/types/api.ts) | TypeScript type definitions for API |

### Backend Key Files

| File | Purpose |
|------|---------|
| [main.go](file:///Users/zhangyi/insights/backend/cmd/api/main.go) | Application entry point and server setup |
| [parser.go](file:///Users/zhangyi/insights/backend/internal/promql/parser/parser.go) | PromQL parsing and AST conversion |
| [analyzer.go](file:///Users/zhangyi/insights/backend/internal/promql/analyzer/analyzer.go) | Static analysis and issue detection |
| [promql.go](file:///Users/zhangyi/insights/backend/internal/handler/promql.go) | HTTP handlers for PromQL endpoints |
| [api.go](file:///Users/zhangyi/insights/backend/internal/dto/api.go) | Go DTO definitions for API |

## Configuration Files

### Frontend Configuration

- [vite.config.ts](file:///Users/zhangyi/insights/frontend/vite.config.ts) - Vite build configuration with proxy
- [tsconfig.json](file:///Users/zhangyi/insights/frontend/tsconfig.json) - TypeScript compiler options
- [.eslintrc.js](file:///Users/zhangyi/insights/frontend/.eslintrc.js) - ESLint rules
- [.prettierrc](file:///Users/zhangyi/insights/frontend/.prettierrc) - Prettier formatting rules

### Backend Configuration

- [go.mod](file:///Users/zhangyi/insights/backend/go.mod) - Go module dependencies
- [go.sum](file:///Users/zhangyi/insights/backend/go.sum) - Go module checksums

### Project Configuration

- [Makefile](file:///Users/zhangyi/insights/Makefile) - Build automation
- [docker-compose.yml](file:///Users/zhangyi/insights/docker-compose.yml) - Docker services
- [.gitignore](file:///Users/zhangyi/insights/.gitignore) - Git ignore patterns

## API Endpoints

### POST /api/promql/parse
- **Purpose**: Parse and analyze PromQL query
- **Request**: `{ query, dialect, options }`
- **Response**: `{ formatted, ast, dependencies, issues }`

### POST /api/query/range
- **Purpose**: Execute range query against Prometheus
- **Request**: `{ query, start, end, step, datasource }`
- **Response**: `{ status, data: { resultType, result } }`

## Dependencies

### Frontend Dependencies
- vue@^3.4.0 - Vue 3 framework
- vue-router@^4.2.5 - Vue Router
- pinia@^2.1.7 - State management
- axios@^1.6.2 - HTTP client
- echarts@^5.4.3 - Charting library
- monaco-editor@^0.45.0 - Code editor

### Backend Dependencies
- github.com/cloudwego/hertz@v0.8.3 - HTTP framework
- github.com/prometheus/prometheus@v0.49.1 - PromQL parser
- github.com/bytedance/sonic@v1.10.2 - JSON library

## Development Workflow

1. **Start Backend**: `make backend-dev` or `cd backend && go run cmd/api/main.go`
2. **Start Frontend**: `make frontend-dev` or `cd frontend && npm run dev`
3. **Access Application**: Open http://localhost:3000
4. **API Endpoints**: http://localhost:8888

## Build and Deployment

### Development
- Frontend: `npm run dev` (Vite dev server)
- Backend: `go run cmd/api/main.go` (Go run)

### Production
- Frontend: `npm run build` (Build to dist/)
- Backend: `go build -o bin/api cmd/api/main.go` (Build binary)

### Docker
- `docker-compose up --build` (Build and run all services)

## Testing

### Frontend Testing
- Unit tests: `npm run test`
- E2E tests: `npm run test:e2e`

### Backend Testing
- Unit tests: `go test ./...`
- Integration tests: `go test ./... -tags=integration`

## Documentation

- [README.md](file:///Users/zhangyi/insights/README.md) - Project overview and getting started
- [ARCHITECTURE.md](file:///Users/zhangyi/insights/ARCHITECTURE.md) - System architecture and design
- [DEVELOPMENT.md](file:///Users/zhangyi/insights/DEVELOPMENT.md) - Development guide and best practices
- [STRUCTURE.md](file:///Users/zhangyi/insights/STRUCTURE.md) - This file - project structure overview
