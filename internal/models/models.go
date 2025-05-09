package models

// --- ENUMы для статусов выражений ---
type ExpressionStatus string

const (
	StatusPending    ExpressionStatus = "pending"
	StatusProcessing ExpressionStatus = "processing"
	StatusCompleted  ExpressionStatus = "completed"
	StatusFailed     ExpressionStatus = "failed"
)

// --- DTO для выражений ---
type Expression struct {
	ID      int              `json:"id"`
	RawExpr string           `json:"raw_expr"`
	Status  ExpressionStatus `json:"status"`
	Result  *string          `json:"result,omitempty"`
	UserID  int              `json:"user_id"`
}

// --- HTTP-запросы/ответы для выражений ---
type ExpressionRequest struct {
	Expression string `json:"expression"`
}

type ExpressionResponse struct {
	ID int `json:"id"`
}

type ExpressionsResponse struct {
	Expressions []Expression `json:"expressions"`
}

type ExpressionDetailResponse struct {
	Expression Expression `json:"expression"`
}

// --- Модель задачи ---
type Task struct {
	ID           int      `json:"id"`
	ExpressionID int      `json:"expression_id"`
	Operation    string   `json:"operation"` // "+", "-", "*", "/"
	Arg1         string   `json:"arg1"`      // может быть число или "res:3"
	Arg2         string   `json:"arg2"`
	Result       *float64 `json:"result,omitempty"`
	Dependencies []int    `json:"dependencies"`
	IsReady      bool     `json:"is_ready"`
}

// --- HTTP/gRPC DTO для задач ---
type TaskResponse struct {
	Task *Task `json:"task"`
}

type TaskResultRequest struct {
	ID     int     `json:"id"`
	Result float64 `json:"result"`
}

// --- Модель пользователя ---
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // не включается в JSON
}

// --- DTO для аутентификации ---
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}