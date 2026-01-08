# Project Architecture

## System Overview

```
┌─────────────────────────────────────────────────────────────┐
│                         Browser                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              Vue 3 + Vite Frontend                  │   │
│  │  ┌──────────────┐  ┌──────────────┐  ┌────────────┐ │   │
│  │  │ Monaco Editor│  │  ECharts     │  │  Vue Router│ │   │
│  │  │   (PromQL)   │  │  (Charts)    │  │            │ │   │
│  │  └──────────────┘  └──────────────┘  └────────────┘ │   │
│  │  ┌──────────────────────────────────────────────┐  │   │
│  │  │         PromQL Workbench Component            │  │   │
│  │  │  - Editor Panel                                │  │   │
│  │  │  - Analysis Panel (AST, Dependencies, Issues)  │  │   │
│  │  │  - Results Panel (Charts, Tables)              │  │   │
│  │  └──────────────────────────────────────────────┘  │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ HTTP/REST API
                              │
┌─────────────────────────────────────────────────────────────┐
│              Cloudwego (Hertz) Backend                      │
│  ┌──────────────────────────────────────────────────────┐   │
│  │                  HTTP Server                        │   │
│  │  - CORS Middleware                                  │   │
│  │  - Request/Response Handling                        │   │
│  └──────────────────────────────────────────────────────┘   │
│                              │                               │
│  ┌──────────────────────────────────────────────────────┐   │
│  │                 Handler Layer                        │   │
│  │  - ParsePromQL Handler                               │   │
│  │  - QueryRange Handler                                │   │
│  └──────────────────────────────────────────────────────┘   │
│                              │                               │
│  ┌──────────────────────────────────────────────────────┐   │
│  │                Business Logic Layer                  │   │
│  │  ┌──────────────┐  ┌──────────────┐                 │   │
│  │  │   Parser     │  │   Analyzer   │                 │   │
│  │  │  Module      │  │   Module     │                 │   │
│  │  │              │  │              │                 │   │
│  │  │ - Parse      │  │ - Dependency │                 │   │
│  │  │ - Format     │  │   Extraction │                 │   │
│  │  │ - AST        │  │ - Issues     │                 │   │
│  │  │   Conversion │  │   Detection  │                 │   │
│  │  └──────────────┘  └──────────────┘                 │   │
│  └──────────────────────────────────────────────────────┘   │
│                              │                               │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              Prometheus Parser Library               │   │
│  │  - Official PromQL Parser                            │   │
│  │  - AST Generation                                     │   │
│  │  - Expression Validation                             │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ (Optional)
                              │
┌─────────────────────────────────────────────────────────────┐
│              External Services                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  Prometheus   │  │    Thanos    │  │VictoriaMetrics│      │
│  │              │  │              │  │              │      │
│  │  Query API   │  │  Query API   │  │  Query API   │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

## Data Flow

### 1. PromQL Parse Flow

```
User Input (PromQL)
       │
       ▼
┌──────────────┐
│ Monaco Editor│
└──────────────┘
       │
       │ POST /api/promql/parse
       │ { query, options: { format, analyze } }
       ▼
┌──────────────┐
│ Hertz Server │
└──────────────┘
       │
       ▼
┌──────────────┐
│  Handler     │
│  - Validate  │
│  - Dispatch  │
└──────────────┘
       │
       ├─────────────────┐
       │                 │
       ▼                 ▼
┌──────────────┐  ┌──────────────┐
│   Parser     │  │  Analyzer    │
│              │  │              │
│ - Parse      │  │ - Extract    │
│ - Format     │  │   Metrics    │
│ - AST → JSON │  │ - Check      │
└──────────────┘  │   Issues     │
       │         └──────────────┘
       │                 │
       └────────┬────────┘
                │
                ▼
         ┌──────────────┐
         │  Response    │
         │ { formatted, │
         │   ast,       │
         │   deps,      │
         │   issues }   │
         └──────────────┘
                │
                ▼
         ┌──────────────┐
         │   Frontend   │
         │  - Update UI │
         │  - Show AST  │
         │  - Show Issues│
         └──────────────┘
```

### 2. Query Range Flow

```
User Click "Run"
       │
       ▼
┌──────────────┐
│ Frontend     │
│ - Get Query  │
│ - Get Time   │
│   Range      │
└──────────────┘
       │
       │ POST /api/query/range
       │ { query, start, end, step, datasource }
       ▼
┌──────────────┐
│ Hertz Server │
└──────────────┘
       │
       ▼
┌──────────────┐
│  Handler     │
│  - Validate  │
│  - Proxy to  │
│    Prometheus│
└──────────────┘
       │
       ▼
┌──────────────┐
│  Prometheus  │
│  Query API   │
└──────────────┘
       │
       ▼
┌──────────────┐
│  Response    │
│ { status,    │
│   data: {    │
│     result   │
│   } }        │
└──────────────┘
       │
       ▼
┌──────────────┐
│   Frontend   │
│  - ECharts   │
│  - Table     │
└──────────────┘
```

## Module Responsibilities

### Frontend Modules

#### 1. Editor Module
- **Component**: Monaco Editor integration
- **Responsibilities**:
  - PromQL syntax highlighting
  - Code completion
  - Error highlighting
  - Format on demand

#### 2. Analysis Module
- **Components**: ASTTree, Dependencies, Issues
- **Responsibilities**:
  - AST visualization
  - Dependency display
  - Issue presentation
  - Interactive highlighting

#### 3. Visualization Module
- **Component**: QueryResults
- **Responsibilities**:
  - Chart rendering (ECharts)
  - Table view
  - Data filtering

### Backend Modules

#### 1. Parser Module
- **Location**: `internal/promql/parser/`
- **Responsibilities**:
  - Parse PromQL string
  - Format PromQL
  - Convert AST to JSON DTO
  - Error handling

#### 2. Analyzer Module
- **Location**: `internal/promql/analyzer/`
- **Responsibilities**:
  - Extract dependencies
  - Detect issues
  - Static analysis rules
  - Performance hints

#### 3. Handler Module
- **Location**: `internal/handler/`
- **Responsibilities**:
  - HTTP request handling
  - Request validation
  - Response formatting
  - Error responses

#### 4. DTO Module
- **Location**: `internal/dto/`
- **Responsibilities**:
  - Request/Response structures
  - Type definitions
  - JSON serialization

## Technology Stack

### Frontend
- **Framework**: Vue 3 (Composition API)
- **Build Tool**: Vite
- **Language**: TypeScript
- **Editor**: Monaco Editor
- **Charts**: ECharts
- **State**: Pinia
- **Router**: Vue Router

### Backend
- **Framework**: Cloudwego (Hertz)
- **Language**: Go 1.21+
- **Parser**: Prometheus PromQL Parser
- **JSON**: Sonic (high-performance JSON)

## API Specification

### POST /api/promql/parse
**Request**:
```json
{
  "query": "sum(rate(http_requests_total[5m]))",
  "dialect": "prometheus",
  "options": {
    "format": true,
    "analyze": true
  }
}
```

**Response**:
```json
{
  "formatted": "sum(rate(http_requests_total[5m]))",
  "ast": {
    "type": "AggregateExpr",
    "op": "sum",
    "expr": {
      "type": "Call",
      "op": "rate",
      "args": [...]
    }
  },
  "dependencies": {
    "metrics": ["http_requests_total"],
    "labels": [],
    "functions": ["rate"],
    "aggregations": ["sum"]
  },
  "issues": []
}
```

### POST /api/query/range
**Request**:
```json
{
  "query": "sum(rate(http_requests_total[5m]))",
  "start": 1700000000,
  "end": 1700003600,
  "step": 15,
  "datasource": "default"
}
```

**Response**:
```json
{
  "status": "success",
  "data": {
    "resultType": "matrix",
    "result": [
      {
        "metric": {},
        "values": [[1700000000, "1.23"]]
      }
    ]
  }
}
```

## Deployment Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Load Balancer                          │
└─────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┴───────────────┐
              │                               │
              ▼                               ▼
┌─────────────────────────┐     ┌─────────────────────────┐
│   Frontend Server 1     │     │   Frontend Server 2     │
│   (Nginx + Vue)         │     │   (Nginx + Vue)         │
└─────────────────────────┘     └─────────────────────────┘
              │                               │
              └───────────────┬───────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────┐
│                    API Gateway                            │
└─────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┴───────────────┐
              │                               │
              ▼                               ▼
┌─────────────────────────┐     ┌─────────────────────────┐
│   Backend Server 1      │     │   Backend Server 2      │
│   (Cloudwego)           │     │   (Cloudwego)           │
└─────────────────────────┘     └─────────────────────────┘
              │                               │
              └───────────────┬───────────────┘
                              │
              ┌───────────────┼───────────────┐
              │               │               │
              ▼               ▼               ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│   Prometheus    │ │     Thanos      │ │VictoriaMetrics  │
└─────────────────┘ └─────────────────┘ └─────────────────┘
```

## Security Considerations

1. **Authentication**: JWT-based authentication for API access
2. **Authorization**: Role-based access control for datasources
3. **Rate Limiting**: Per-user and per-datasource rate limits
4. **Input Validation**: Strict validation of PromQL queries
5. **CORS**: Configurable CORS policies
6. **Audit Logging**: All queries logged for compliance

## Performance Optimization

### Frontend
- Code splitting with dynamic imports
- Lazy loading of components
- Virtual scrolling for large datasets
- Debounced API calls
- Caching of parse results

### Backend
- Connection pooling for Prometheus clients
- Caching of metadata
- Parallel processing of independent analyses
- Efficient AST traversal
- Streaming of large query results

## Monitoring & Observability

- **Metrics**: Request latency, error rates, query complexity
- **Logging**: Structured logging with correlation IDs
- **Tracing**: Distributed tracing for request flows
- **Health Checks**: Liveness and readiness probes
