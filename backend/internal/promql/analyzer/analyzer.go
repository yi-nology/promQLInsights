package analyzer

import (
	"fmt"
	"strings"

	"github.com/insights/backend/internal/dto"
	"github.com/prometheus/prometheus/promql/parser"
)

type Analyzer struct{}

func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

func (a *Analyzer) Analyze(query string) (*dto.Dependencies, []dto.Issue) {
	expr, err := parser.ParseExpr(query)
	if err != nil {
		return nil, []dto.Issue{
			{
				Level:   "error",
				Code:    "PARSE_ERROR",
				Message: err.Error(),
			},
		}
	}

	dependencies := a.extractDependencies(expr)
	issues := a.checkIssues(expr, dependencies)

	return dependencies, issues
}

func (a *Analyzer) extractDependencies(expr parser.Expr) *dto.Dependencies {
	deps := &dto.Dependencies{
		Metrics:      []string{},
		Labels:       []string{},
		Matchers:     []dto.Matcher{},
		Functions:    []string{},
		Aggregations: []string{},
	}

	a.walkExpr(expr, deps)

	return deps
}

func (a *Analyzer) walkExpr(expr parser.Expr, deps *dto.Dependencies) {
	if expr == nil {
		return
	}

	switch e := expr.(type) {
	case *parser.AggregateExpr:
		deps.Aggregations = append(deps.Aggregations, e.Op.String())
		a.walkExpr(e.Expr, deps)

	case *parser.Call:
		deps.Functions = append(deps.Functions, e.Func.Name)
		for _, arg := range e.Args {
			a.walkExpr(arg, deps)
		}

	case *parser.BinaryExpr:
		a.walkExpr(e.LHS, deps)
		a.walkExpr(e.RHS, deps)

	case *parser.UnaryExpr:
		a.walkExpr(e.Expr, deps)

	case *parser.ParenExpr:
		a.walkExpr(e.Expr, deps)

	case *parser.VectorSelector:
		if e.Name != "" {
			deps.Metrics = append(deps.Metrics, e.Name)
		}
		for _, matcher := range e.LabelMatchers {
			deps.Labels = append(deps.Labels, matcher.Name)
			deps.Matchers = append(deps.Matchers, dto.Matcher{
				Label: matcher.Name,
				Op:    matcher.Type.String(),
				Value: matcher.Value,
			})
		}

	case *parser.MatrixSelector:
		a.walkExpr(e.VectorSelector, deps)

	case *parser.SubqueryExpr:
		a.walkExpr(e.Expr, deps)
	}
}

func (a *Analyzer) checkIssues(expr parser.Expr, deps *dto.Dependencies) []dto.Issue {
	var issues []dto.Issue

	issues = append(issues, a.checkHighCardinality(expr, deps)...)
	issues = append(issues, a.checkRateIncrease(expr, deps)...)
	issues = append(issues, a.checkRangeSize(expr, deps)...)

	return issues
}

func (a *Analyzer) checkHighCardinality(expr parser.Expr, deps *dto.Dependencies) []dto.Issue {
	var issues []dto.Issue

	if len(deps.Metrics) > 0 && len(deps.Matchers) == 0 {
		issues = append(issues, dto.Issue{
			Level:   "warn",
			Code:    "HIGH_CARDINALITY",
			Message: "Metric selector has no label filters, which may lead to high cardinality",
		})
	}

	if len(deps.Aggregations) > 0 && len(deps.Matchers) > 3 {
		issues = append(issues, dto.Issue{
			Level:   "warn",
			Code:    "HIGH_CARDINALITY",
			Message: "Aggregation with many label matchers may lead to high cardinality",
		})
	}

	return issues
}

func (a *Analyzer) checkRateIncrease(expr parser.Expr, deps *dto.Dependencies) []dto.Issue {
	var issues []dto.Issue

	for _, fn := range deps.Functions {
		if fn == "rate" || fn == "increase" {
			for _, metric := range deps.Metrics {
				if !strings.HasSuffix(metric, "_total") && !strings.HasSuffix(metric, "_count") && !strings.HasSuffix(metric, "_sum") && !strings.HasSuffix(metric, "_bucket") {
					issues = append(issues, dto.Issue{
						Level:   "warn",
						Code:    "RATE_ON_NON_COUNTER",
						Message: fmt.Sprintf("rate() or increase() is used on metric '%s' which may not be a counter", metric),
					})
				}
			}
		}
	}

	return issues
}

func (a *Analyzer) checkRangeSize(expr parser.Expr, deps *dto.Dependencies) []dto.Issue {
	var issues []dto.Issue

	a.walkExprForRange(expr, &issues)

	return issues
}

func (a *Analyzer) walkExprForRange(expr parser.Expr, issues *[]dto.Issue) {
	if expr == nil {
		return
	}

	switch e := expr.(type) {
	case *parser.MatrixSelector:
		if e.Range > 3600 {
			*issues = append(*issues, dto.Issue{
				Level:   "warn",
				Code:    "LARGE_RANGE",
				Message: fmt.Sprintf("Range selector with duration %s is very large, consider using a smaller range", e.Range.String()),
			})
		}

	case *parser.SubqueryExpr:
		if e.Range > 3600 {
			*issues = append(*issues, dto.Issue{
				Level:   "warn",
				Code:    "LARGE_RANGE",
				Message: fmt.Sprintf("Subquery with range %s is very large, consider using a smaller range", e.Range.String()),
			})
		}

	default:
		a.walkChildrenForRange(expr, issues)
	}
}

func (a *Analyzer) walkChildrenForRange(expr parser.Expr, issues *[]dto.Issue) {
	switch e := expr.(type) {
	case *parser.AggregateExpr:
		a.walkExprForRange(e.Expr, issues)
	case *parser.Call:
		for _, arg := range e.Args {
			a.walkExprForRange(arg, issues)
		}
	case *parser.BinaryExpr:
		a.walkExprForRange(e.LHS, issues)
		a.walkExprForRange(e.RHS, issues)
	case *parser.UnaryExpr:
		a.walkExprForRange(e.Expr, issues)
	case *parser.ParenExpr:
		a.walkExprForRange(e.Expr, issues)
	case *parser.SubqueryExpr:
		a.walkExprForRange(e.Expr, issues)
	}
}
