package main

type Purchase struct {
    ID         int     `json:"id"`
    UserID     int     `json:"user_id"`
    ProductID  int     `json:"product_id"`
    Quantity   int     `json:"quantity"`
    TotalPrice float64 `json:"total_price"`
}
