package models

type Status string

// Константы для возможных статусов выражения
const (
	StatusPending    Status = "PENDING"    
	StatusProcessing Status = "PROCESSING" 
	StatusCompleted  Status = "COMPLETED"  
	StatusError      Status = "ERROR"      
)