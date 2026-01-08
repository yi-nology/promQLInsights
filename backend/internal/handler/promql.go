package handler

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/insights/backend/internal/dto"
	"github.com/insights/backend/internal/promql/analyzer"
	"github.com/insights/backend/internal/promql/parser"
)

var (
	p *parser.Parser
	a *analyzer.Analyzer
)

func init() {
	p = parser.NewParser()
	a = analyzer.NewAnalyzer()
}

func ParsePromQL(ctx context.Context, c *app.RequestContext) {
	var req dto.ParseRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, dto.ParseResponse{
			Error: "Invalid request body",
		})
		return
	}

	if req.Query == "" {
		c.JSON(400, dto.ParseResponse{
			Error: "Query is required",
		})
		return
	}

	response := dto.ParseResponse{}

	if req.Options == nil {
		req.Options = &dto.ParseOptions{}
	}

	if req.Options.Format || req.Options.Analyze {
		parseResp, err := p.Parse(req.Query)
		if err != nil {
			c.JSON(500, dto.ParseResponse{
				Error: err.Error(),
			})
			return
		}

		response.Formatted = parseResp.Formatted
		response.AST = parseResp.AST
	}

	if req.Options.Analyze {
		dependencies, issues := a.Analyze(req.Query)
		response.Dependencies = dependencies
		response.Issues = issues
	}

	c.JSON(200, response)
}

func QueryRange(ctx context.Context, c *app.RequestContext) {
	var req dto.QueryRangeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, dto.QueryRangeResponse{
			Status: "error",
		})
		return
	}

	result := generateMockData(req.Start, req.End, req.Step)

	c.JSON(200, dto.QueryRangeResponse{
		Status: "success",
		Data: struct {
			ResultType string       `json:"resultType"`
			Result     []dto.Series `json:"result"`
		}{
			ResultType: "matrix",
			Result:     result,
		},
	})
}

func generateMockData(start, end int64, step int) []dto.Series {
	series := []dto.Series{}

	statuses := []string{"200", "201", "400", "500"}
	methods := []string{"GET", "POST", "PUT", "DELETE"}

	rand.Seed(time.Now().UnixNano())

	for _, status := range statuses {
		for _, method := range methods {
			metric := map[string]string{
				"job":     "api",
				"status":  status,
				"method":  method,
				"service": "user-service",
			}

			values := make([][2]interface{}, 0)

			for ts := start; ts <= end; ts += int64(step) {
				timestamp := float64(ts)

				baseValue := 100.0
				switch status {
				case "200":
					baseValue = 500.0
				case "400":
					baseValue = 50.0
				case "500":
					baseValue = 10.0
				}

				switch method {
				case "POST":
					baseValue *= 0.8
				case "PUT":
					baseValue *= 0.6
				case "DELETE":
					baseValue *= 0.4
				}

				value := baseValue + math.Sin(float64(ts)/100.0)*20.0 + (rand.Float64() * 30.0)
				value = math.Max(0, value)

				values = append(values, [2]interface{}{timestamp, fmt.Sprintf("%.2f", value)})
			}

			series = append(series, dto.Series{
				Metric: metric,
				Values: values,
			})
		}
	}

	return series
}
