package dto

type ParseRequest struct {
	Query   string            `json:"query"`
	Dialect string            `json:"dialect,omitempty"`
	Options *ParseOptions     `json:"options,omitempty"`
}

type ParseOptions struct {
	Format  bool `json:"format,omitempty"`
	Analyze bool `json:"analyze,omitempty"`
}

type ParseResponse struct {
	Formatted    string        `json:"formatted,omitempty"`
	AST          *ASTNode      `json:"ast,omitempty"`
	Dependencies *Dependencies `json:"dependencies,omitempty"`
	Issues       []Issue       `json:"issues,omitempty"`
	Error        string        `json:"error,omitempty"`
}

type ASTNode struct {
	Type  string                 `json:"type"`
	Span  *[2]int                `json:"span,omitempty"`
	Text  string                 `json:"text,omitempty"`
	Meta  map[string]interface{} `json:"meta,omitempty"`
	Expr  *ASTNode               `json:"expr,omitempty"`
	Op    string                 `json:"op,omitempty"`
	By    []string               `json:"by,omitempty"`
	Left  *ASTNode               `json:"left,omitempty"`
	Right *ASTNode               `json:"right,omitempty"`
	Args  []*ASTNode             `json:"args,omitempty"`
}

type Dependencies struct {
	Metrics      []string   `json:"metrics"`
	Labels       []string   `json:"labels"`
	Matchers     []Matcher  `json:"matchers"`
	Functions    []string   `json:"functions"`
	Aggregations []string   `json:"aggregations"`
}

type Matcher struct {
	Label string `json:"label"`
	Op    string `json:"op"`
	Value string `json:"value"`
}

type Issue struct {
	Level   string  `json:"level"`
	Code    string  `json:"code"`
	Message string  `json:"message"`
	Span    *[2]int `json:"span,omitempty"`
}

type QueryRangeRequest struct {
	Query      string `json:"query"`
	Start      int64  `json:"start"`
	End        int64  `json:"end"`
	Step       int    `json:"step"`
	Datasource string `json:"datasource,omitempty"`
}

type QueryRangeResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string   `json:"resultType"`
		Result     []Series `json:"result"`
	} `json:"data"`
}

type Series struct {
	Metric map[string]string `json:"metric"`
	Values [][2]interface{}  `json:"values"`
}
