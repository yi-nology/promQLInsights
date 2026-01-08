# PromQL Insights

A powerful PromQL query analyzer and visualization tool with real-time monitoring data integration.

## Features

- **PromQL Parser & Formatter**: Parse and format PromQL queries with syntax highlighting
- **AST Visualization**: Interactive abstract syntax tree with expand/collapse functionality
- **Dependency Analysis**: Extract metrics, labels, functions, and aggregations from queries
- **Static Analysis**: Detect potential issues and provide optimization suggestions
- **Query Execution**: Execute PromQL queries against Prometheus with real-time results
- **Data Visualization**: Interactive charts and table views for query results
- **Mock Data Generation**: Built-in mock data generator for testing

## Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                     Frontend (Vue 3)                          │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │  Monaco Editor  │  AST Tree  │  Charts  │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                            ↓ HTTP API
┌─────────────────────────────────────────────────────────────────────┐
│                  Backend (Go + Hertz)                      │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │  Parser  │  Analyzer  │  Query Proxy  │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                            ↓ HTTP API
┌─────────────────────────────────────────────────────────────────────┐
│              Prometheus + Pushgateway (Docker)               │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │  Time Series DB  │  Query Engine  │  Pushgateway │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

## Tech Stack

### Frontend
- **Framework**: Vue 3 (Composition API)
- **Build Tool**: Vite
- **Language**: TypeScript
- **Editor**: Monaco Editor
- **Charts**: ECharts
- **State Management**: Pinia
- **Router**: Vue Router

### Backend
- **Framework**: Cloudwego (Hertz)
- **Language**: Go 1.21+
- **PromQL Parser**: Prometheus official parser
- **JSON Library**: Sonic

### Monitoring
- **Time Series DB**: Prometheus v2.48.0
- **Metrics Gateway**: Pushgateway v1.6.2
- **Deployment**: Docker

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Node.js 16 or higher
- npm 8 or higher
- Docker (optional, for Prometheus deployment)

### One-Click Build & Install

Run the build script to build and install all components:

```bash
chmod +x build.sh
./build.sh
```

This will:
1. Check all dependencies (Go, Node.js, npm, Docker)
2. Build the backend binary
3. Build the frontend distribution
4. Start Prometheus and Pushgateway containers
5. Display service URLs and usage instructions

### Manual Installation

#### Backend

```bash
cd backend
go mod download
go build -o bin/api cmd/api/main.go
./bin/api
```

Backend will be available at `http://localhost:8888`

#### Frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend will be available at `http://localhost:3000`

#### Prometheus (Optional)

```bash
# Start Prometheus
docker run -d --name prometheus -p 9090:9090 \
  -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
  prom/prometheus:v2.48.0

# Start Pushgateway
docker run -d --name pushgateway -p 9091:9091 \
  prom/pushgateway:v1.6.2
```

## Usage

### 1. Start Services

```bash
# Start backend
cd backend && ./bin/api

# Start frontend (dev mode)
cd frontend && npm run dev

# Start frontend (production)
cd frontend && npm run preview

# Push metrics (optional)
python3 push_metrics.py
```

### 2. Access Services

| Service | URL | Description |
|----------|-----|-------------|
| Frontend | http://localhost:3000 | PromQL Workbench |
| Backend API | http://localhost:8888 | REST API |
| Prometheus | http://localhost:9090 | Monitoring Data |
| Pushgateway | http://localhost:9091 | Metrics Gateway |

### 3. Use the Application

#### PromQL Editor

1. Enter your PromQL query in the Monaco Editor
2. Click "Format" to format the query
3. Click "Parse" to analyze the query

#### AST Visualization

- View the abstract syntax tree structure
- Click on nodes to expand/collapse
- Click on nodes to highlight corresponding code

#### Dependency Analysis

- View extracted metrics, labels, functions
- See aggregations and matchers
- Understand query dependencies

#### Static Analysis

- Review detected issues (error/warn/info)
- See optimization suggestions
- Improve query performance

#### Query Execution

1. Set time range (start, end, step)
2. Click "Run" to execute the query
3. View results in chart or table format
4. Filter series by labels

## API Documentation

### POST /api/promql/parse

Parse and analyze a PromQL query.

**Request:**
```json
{
  "query": "sum(rate(http_requests_total[5m])) by (status)",
  "dialect": "prometheus",
  "options": {
    "format": true,
    "analyze": true
  }
}
```

**Response:**
```json
{
  "formatted": "sum by (status) (rate(http_requests_total[5m]))",
  "ast": {
    "type": "AggregateExpr",
    "op": "sum",
    "by": ["status"],
    "expr": { ... }
  },
  "dependencies": {
    "metrics": ["http_requests_total"],
    "labels": ["status"],
    "functions": ["rate"],
    "aggregations": ["sum"]
  },
  "issues": [
    {
      "level": "warn",
      "code": "LARGE_RANGE",
      "message": "Range selector with duration 5m0s is very large"
    }
  ]
}
```

### POST /api/query/range

Execute a PromQL range query.

**Request:**
```json
{
  "query": "http_requests_total",
  "start": 1700000000,
  "end": 1700003600,
  "step": 15,
  "datasource": "default"
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "resultType": "matrix",
    "result": [
      {
        "metric": {
          "__name__": "http_requests_total",
          "job": "api",
          "status": "200"
        },
        "values": [
          [1700000000, "500"],
          [1700000015, "520"],
          ...
        ]
      }
    ]
  }
}
```

## Project Structure

```
insights/
├── backend/                          # Go Backend
│   ├── cmd/
│   │   └── api/
│   │       └── main.go            # Entry point
│   ├── internal/
│   │   ├── handler/               # HTTP handlers
│   │   │   └── promql.go       # API endpoints
│   │   ├── middleware/            # Middleware (CORS)
│   │   ├── dto/                  # Data transfer objects
│   │   └── promql/               # PromQL processing
│   │       ├── parser/            # Parser module
│   │       └── analyzer/          # Analyzer module
│   ├── go.mod                    # Go module definition
│   └── go.sum                    # Go dependencies
├── frontend/                         # Vue Frontend
│   ├── src/
│   │   ├── components/            # Vue components
│   │   │   ├── ASTTree.vue      # AST visualization
│   │   │   ├── ASTNode.vue      # Recursive AST node
│   │   │   ├── Dependencies.vue # Dependencies panel
│   │   │   ├── Issues.vue       # Issues panel
│   │   │   └── QueryResults.vue # Query results
│   │   ├── views/                 # Page views
│   │   │   └── PromQLWorkbench.vue
│   │   ├── router/                # Vue Router
│   │   ├── types/                 # TypeScript types
│   │   └── main.ts               # Entry point
│   ├── package.json
│   ├── vite.config.ts
│   └── tsconfig.json
├── prometheus.yml                    # Prometheus config
├── push_metrics.py                  # Metrics push script
├── build.sh                        # Build & install script
└── README.md                       # This file
```

## Development

### Backend Development

```bash
cd backend
go run cmd/api/main.go
```

### Frontend Development

```bash
cd frontend
npm run dev
```

### Running Tests

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
npm run test
```

### Code Quality

```bash
# Backend linting
cd backend
go fmt ./...
go vet ./...

# Frontend linting
cd frontend
npm run lint
npm run type-check
```

## Deployment

### Production Build

```bash
# Build backend
cd backend
go build -ldflags="-s -w" -o bin/api cmd/api/main.go

# Build frontend
cd frontend
npm run build
```

### Docker Deployment

```bash
# Build and start all services
./build.sh

# Or use Docker Compose
docker-compose -f docker-compose.prometheus.yml up -d
```

### Environment Variables

Create a `.env` file in the project root:

```env
# Backend
BACKEND_PORT=8888
BACKEND_HOST=0.0.0.0

# Frontend
FRONTEND_PORT=3000
VITE_API_BASE_URL=http://localhost:8888

# Prometheus
PROMETHEUS_PORT=9090
PUSHGATEWAY_PORT=9091
```

## Monitoring

### Pushing Metrics

Use the provided Python script to push mock metrics:

```bash
python3 push_metrics.py
```

This will push HTTP request metrics every 10 seconds with various labels:
- Services: user-service, order-service, payment-service
- Methods: GET, POST, PUT, DELETE
- Status codes: 200, 201, 400, 500

### Prometheus Queries

Access the Prometheus UI to query metrics:

```promql
# Query all HTTP requests
http_requests_total

# Query by service
http_requests_total{service="user-service"}

# Query by status
http_requests_total{status="200"}

# Calculate rate
rate(http_requests_total[5m])

# Calculate increase
increase(http_requests_total[1h])

# Group by status
sum by (status) (http_requests_total)
```

## Troubleshooting

### Backend won't start

```bash
# Check if port is in use
lsof -ti:8888

# Kill process using port
kill -9 <PID>

# Check logs
./bin/api
```

### Frontend won't start

```bash
# Clear cache and reinstall
cd frontend
rm -rf node_modules package-lock.json
npm install

# Check port
lsof -ti:3000
```

### Prometheus can't connect to Pushgateway

```bash
# Check container status
docker ps

# Check container logs
docker logs prometheus
docker logs pushgateway

# Check network
docker network inspect insights_monitoring

# Restart containers
docker restart prometheus pushgateway
```

### Build errors

```bash
# Clean Go cache
cd backend
go clean -cache

# Update dependencies
go mod tidy

# Rebuild
go build -o bin/api cmd/api/main.go
```

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see LICENSE file for details

## Acknowledgments

- [Prometheus](https://prometheus.io/) - Time series database
- [Cloudwego](https://www.cloudwego.io/) - Go HTTP framework
- [Vue.js](https://vuejs.org/) - Progressive JavaScript framework
- [Monaco Editor](https://microsoft.github.io/monaco-editor/) - Code editor
- [ECharts](https://echarts.apache.org/) - Visualization library

## Contact

- Project Name: PromQL Insights
- Issues: [GitHub Issues](https://github.com/your-org/insights/issues)
- Discussions: [GitHub Discussions](https://github.com/your-org/insights/discussions)

---

**Made with ❤️ by the PromQL Insights Team**
