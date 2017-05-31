package models

//User Struct
type USER struct {
	ID    int64  `range:"RANGE" json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
