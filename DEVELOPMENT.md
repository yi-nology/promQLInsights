# Development Guide

## Setup Development Environment

### Prerequisites
- Node.js 18+ and npm
- Go 1.21+
- Git
- Docker (optional, for containerized development)

### Initial Setup

1. **Clone the repository**
```bash
git clone <repository-url>
cd insights
```

2. **Install frontend dependencies**
```bash
cd frontend
npm install
cd ..
```

3. **Install backend dependencies**
```bash
cd backend
go mod download
cd ..
```

## Running the Project

### Development Mode

#### Option 1: Using Make (Recommended)

```bash
# Terminal 1: Start backend
make backend-dev

# Terminal 2: Start frontend
make frontend-dev
```

#### Option 2: Manual Commands

**Backend:**
```bash
cd backend
go run cmd/api/main.go
```

**Frontend:**
```bash
cd frontend
npm run dev
```

### Production Build

```bash
# Build frontend
make frontend-build

# Build backend
cd backend
go build -o bin/api cmd/api/main.go
```

## Development Workflow

### Frontend Development

#### Adding a New Component

1. Create component file in `frontend/src/components/`
```bash
touch frontend/src/components/NewComponent.vue
```

2. Implement component with TypeScript
```vue
<script setup lang="ts">
import { ref } from 'vue'

interface Props {
  title: string
}

const props = defineProps<Props>()
const count = ref(0)
</script>

<template>
  <div class="new-component">
    <h2>{{ title }}</h2>
    <p>Count: {{ count }}</p>
  </div>
</template>

<style scoped>
.new-component {
  padding: 20px;
}
</style>
```

3. Import and use in parent component
```vue
<script setup lang="ts">
import NewComponent from '@/components/NewComponent.vue'
</script>

<template>
  <NewComponent title="Hello World" />
</template>
```

#### Adding a New API Call

1. Define types in `frontend/src/types/api.ts`
```typescript
export interface NewRequest {
  param1: string
  param2: number
}

export interface NewResponse {
  result: string
}
```

2. Create API function
```typescript
// frontend/src/api/newApi.ts
import axios from 'axios'
import type { NewRequest, NewResponse } from '@/types/api'

export async function callNewApi(request: NewRequest): Promise<NewResponse> {
  const response = await axios.post<NewResponse>('/api/new-endpoint', request)
  return response.data
}
```

3. Use in component
```vue
<script setup lang="ts">
import { ref } from 'vue'
import { callNewApi } from '@/api/newApi'

const result = ref('')

async function handleApiCall() {
  try {
    const response = await callNewApi({ param1: 'test', param2: 123 })
    result.value = response.result
  } catch (error) {
    console.error('API call failed:', error)
  }
}
</script>
```

### Backend Development

#### Adding a New API Endpoint

1. Define DTO in `backend/internal/dto/api.go`
```go
type NewRequest struct {
    Param1 string `json:"param1"`
    Param2 int    `json:"param2"`
}

type NewResponse struct {
    Result string `json:"result"`
}
```

2. Create handler in `backend/internal/handler/newHandler.go`
```go
package handler

import (
    "context"
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/insights/backend/internal/dto"
)

func NewEndpoint(ctx context.Context, c *app.RequestContext) {
    var req dto.NewRequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, dto.NewResponse{
            Result: "Invalid request",
        })
        return
    }

    c.JSON(200, dto.NewResponse{
        Result: "Success",
    })
}
```

3. Register route in `backend/cmd/api/main.go`
```go
h.POST("/api/new-endpoint", handler.NewEndpoint)
```

#### Adding a New Analysis Rule

1. Add rule function in `backend/internal/promql/analyzer/analyzer.go`
```go
func (a *Analyzer) checkNewRule(expr parser.Expr, deps *dto.Dependencies) []dto.Issue {
    var issues []dto.Issue

    // Your analysis logic here
    if someCondition {
        issues = append(issues, dto.Issue{
            Level:   "warn",
            Code:    "NEW_RULE",
            Message: "New rule detected",
        })
    }

    return issues
}
```

2. Call the rule in `checkIssues` function
```go
func (a *Analyzer) checkIssues(expr parser.Expr, deps *dto.Dependencies) []dto.Issue {
    var issues []dto.Issue

    issues = append(issues, a.checkHighCardinality(expr, deps)...)
    issues = append(issues, a.checkRateIncrease(expr, deps)...)
    issues = append(issues, a.checkRangeSize(expr, deps)...)
    issues = append(issues, a.checkNewRule(expr, deps)...) // Add new rule

    return issues
}
```

## Code Style

### Frontend (Vue + TypeScript)

- Use Composition API with `<script setup>`
- Use TypeScript for all new code
- Follow Vue 3 style guide
- Use Prettier for formatting
- Use ESLint for linting

**Example:**
```vue
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

interface Props {
  title: string
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'update', value: string): void
}>()

const count = ref(0)
const doubled = computed(() => count.value * 2)

onMounted(() => {
  console.log('Component mounted')
})

function increment() {
  count.value++
  emit('update', count.value.toString())
}
</script>
```

### Backend (Go)

- Follow Go best practices
- Use meaningful variable names
- Add comments for exported functions
- Handle errors properly
- Use context for request-scoped values

**Example:**
```go
package parser

import (
    "context"
    "fmt"
)

// Parse parses a PromQL query and returns the AST
func (p *Parser) Parse(ctx context.Context, query string) (*dto.ParseResponse, error) {
    if query == "" {
        return nil, fmt.Errorf("query cannot be empty")
    }

    expr, err := parser.ParseExpr(query)
    if err != nil {
        return nil, fmt.Errorf("parse error: %w", err)
    }

    return &dto.ParseResponse{
        Formatted: expr.String(),
        AST:       p.convertAST(expr),
    }, nil
}
```

## Testing

### Frontend Testing

```bash
cd frontend
npm run test
```

### Backend Testing

```bash
cd backend
go test ./...
```

### Integration Testing

```bash
# Start services
make backend-dev
make frontend-dev

# Run integration tests
npm run test:integration
```

## Debugging

### Frontend Debugging

1. Open Chrome DevTools (F12)
2. Use Vue DevTools extension for component inspection
3. Check Network tab for API calls
4. Use console.log for debugging

### Backend Debugging

1. Use `fmt.Printf` or `log.Printf` for logging
2. Use Delve debugger:
```bash
cd backend
dlv debug cmd/api/main.go
```

3. Check logs in terminal

## Common Issues

### Frontend

**Issue: Monaco Editor not loading**
- Check if monaco-editor is installed: `npm list monaco-editor`
- Clear cache: `rm -rf node_modules/.vite`

**Issue: API calls failing with CORS error**
- Ensure backend is running on port 8888
- Check CORS middleware in backend

**Issue: TypeScript errors**
- Run `npm run lint` to see all errors
- Check tsconfig.json configuration

### Backend

**Issue: Module not found**
- Run `go mod tidy`
- Check go.mod file

**Issue: Port already in use**
- Kill process using port 8888:
```bash
lsof -ti:8888 | xargs kill -9
```

**Issue: Prometheus parser errors**
- Ensure prometheus package is imported correctly
- Check Go version compatibility

## Performance Optimization

### Frontend

- Use `computed` for derived state
- Implement virtual scrolling for large lists
- Lazy load components with `defineAsyncComponent`
- Debounce API calls
- Use `v-show` instead of `v-if` for frequent toggles

### Backend

- Use connection pooling for external services
- Implement caching for metadata
- Use goroutines for parallel processing
- Optimize database queries
- Use streaming for large responses

## Deployment

### Docker Deployment

```bash
# Build and run with Docker Compose
docker-compose up --build

# Or build individual services
docker build -t insights-frontend ./frontend
docker build -t insights-backend ./backend
```

### Manual Deployment

**Frontend:**
```bash
cd frontend
npm run build
# Deploy dist/ directory to web server
```

**Backend:**
```bash
cd backend
go build -o api cmd/api/main.go
./api
```

## Contributing

1. Create a feature branch
2. Make your changes
3. Run tests and linting
4. Commit with clear messages
5. Push and create a pull request

## Resources

- [Vue 3 Documentation](https://vuejs.org/)
- [Vite Documentation](https://vitejs.dev/)
- [Cloudwego Documentation](https://www.cloudwego.io/)
- [Prometheus PromQL Documentation](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [Monaco Editor API](https://microsoft.github.io/monaco-editor/api/index.html)
- [ECharts Documentation](https://echarts.apache.org/en/index.html)
