package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/Dided08/Calculator/internal/calculator"
)

type Request struct {
    Expression string `json:"expression"`
}

type Response struct {
    Result float64 `json:"result"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    var req Request
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    result, err := calculator.Calculate(req.Expression)
    if err != nil {
        errorResp := ErrorResponse{Error: fmt.Sprintf("%v", err)}
        json.NewEncoder(w).Encode(errorResp)
        w.WriteHeader(http.StatusUnprocessableEntity)
        return
    }

    resp := Response{Result: result}
    json.NewEncoder(w).Encode(resp)
    w.WriteHeader(http.StatusOK)
}