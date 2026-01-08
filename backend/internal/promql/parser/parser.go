package parser

import (
	"fmt"

	"github.com/insights/backend/internal/dto"
	"github.com/prometheus/prometheus/promql/parser"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(query string) (*dto.ParseResponse, error) {
	expr, err := parser.ParseExpr(query)
	if err != nil {
		return &dto.ParseResponse{
			Error: err.Error(),
		}, nil
	}

	response := &dto.ParseResponse{}

	response.Formatted = expr.String()
	response.AST = p.convertAST(expr)

	return response, nil
}

func (p *Parser) convertAST(expr parser.Expr) *dto.ASTNode {
	if expr == nil {
		return nil
	}

	node := &dto.ASTNode{
		Type: getExprType(expr),
	}

	switch e := expr.(type) {
	case *parser.AggregateExpr:
		node.Op = e.Op.String()
		if len(e.Grouping) > 0 {
			node.By = e.Grouping
		}
		node.Expr = p.convertAST(e.Expr)

	case *parser.Call:
		node.Op = e.Func.Name
		args := make([]*dto.ASTNode, len(e.Args))
		for i, arg := range e.Args {
			args[i] = p.convertAST(arg)
		}
		node.Args = args

	case *parser.BinaryExpr:
		node.Op = e.Op.String()
		node.Left = p.convertAST(e.LHS)
		node.Right = p.convertAST(e.RHS)

	case *parser.UnaryExpr:
		node.Op = e.Op.String()
		node.Expr = p.convertAST(e.Expr)

	case *parser.ParenExpr:
		node.Expr = p.convertAST(e.Expr)

	case *parser.VectorSelector:
		node.Text = e.Name

	case *parser.MatrixSelector:
		node.Expr = p.convertAST(e.VectorSelector)

	case *parser.SubqueryExpr:
		node.Expr = p.convertAST(e.Expr)

	case *parser.NumberLiteral:
		node.Text = fmt.Sprintf("%g", e.Val)

	case *parser.StringLiteral:
		node.Text = e.Val
	}

	return node
}

func getExprType(expr parser.Expr) string {
	switch expr.(type) {
	case *parser.AggregateExpr:
		return "AggregateExpr"
	case *parser.Call:
		return "Call"
	case *parser.BinaryExpr:
		return "BinaryExpr"
	case *parser.UnaryExpr:
		return "UnaryExpr"
	case *parser.ParenExpr:
		return "ParenExpr"
	case *parser.VectorSelector:
		return "VectorSelector"
	case *parser.MatrixSelector:
		return "MatrixSelector"
	case *parser.SubqueryExpr:
		return "SubqueryExpr"
	case *parser.NumberLiteral:
		return "NumberLiteral"
	case *parser.StringLiteral:
		return "StringLiteral"
	default:
		return "Unknown"
	}
}

func FormatQuery(query string) string {
	expr, err := parser.ParseExpr(query)
	if err != nil {
		return query
	}
	return expr.String()
}
