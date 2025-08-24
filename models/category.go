package models

type Category struct {
    ID        string `json:"id,omitempty"`
    UserID    string `json:"user_id"`
    Name      string `json:"name"`
    CreatedAt string `json:"created_at,omitempty"`
}
