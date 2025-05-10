package models

// ExpressionStatus определяет статус вычисления выражения.
type ExpressionStatus string

const (
	StatusPending    ExpressionStatus = "pending"
	StatusProcessing ExpressionStatus = "processing"
	StatusCompleted  ExpressionStatus = "completed"
	StatusFailed     ExpressionStatus = "failed"
)

// Expression представляет математическое выражение.
type Expression struct {
	ID       int              `json:"id"`
	UserID   int              `json:"user_id"`
	RawExpr  string           `json:"expression"`
	Result   *string          `json:"result,omitempty"`
	Status   ExpressionStatus `json:"status"`
}