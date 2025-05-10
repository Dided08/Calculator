package models

// ExpressionRequest представляет входной запрос от пользователя на вычисление выражения.
type ExpressionRequest struct {
	UserID     int    `json:"user_id"`      // ID пользователя
	Expression string `json:"expression"`   // Само выражение
}

// ExpressionResponse содержит ID добавленного выражения.
type ExpressionResponse struct {
	ID int `json:"id"`
}

// ExpressionsResponse содержит список всех выражений.
type ExpressionsResponse struct {
	Expressions []Expression `json:"expressions"`
}

// ExpressionDetailResponse содержит детальную информацию об одном выражении.
type ExpressionDetailResponse struct {
	Expression Expression `json:"expression"`
}