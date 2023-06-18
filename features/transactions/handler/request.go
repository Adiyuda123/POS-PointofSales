package handler

type TransactionRequest1 struct {
	ExternalID  string `json:"external_id" form:"external_id"`
	CallbackURL string `json:"callback_url" form:"callback_url"`
	Type        string `json:"type" form:"type"`
	Amount      int    `json:"amount" form:"amount"`
	ItemID      uint   `json:"item_id" form:"item_id"`
	OrderID     string `json:"order_id" form:"order_id"`
}

type TransactionRequest struct {
	ExternalID  string  `json:"external_id"`
	CallbackURL string  `json:"callback_url"`
	ReferenceID string  `json:"reference_id"`
	Type        string  `json:"type"`
	Currency    string  `json:"currency"`
	Amount      float64 `json:"amount"`
	ExpiresAt   string  `json:"expires_at"`
}

type TransactionDetailsRequest struct {
	Customer string              `json:"customer" form:"customer"`
	Details  []ItemDetailRequest `json:"details" form:"details"`
}

type ItemDetailRequest struct {
	ProductID uint `json:"product_id" form:"product_id"`
	Quantity  int  `json:"product_pcs" form:"product_pcs"`
}
