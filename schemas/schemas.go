package schemas

type AuthHeader struct {
	Authorization string `json:"Authorization"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Symbol struct {
	Symbol string `json:"symbol"`
}

type BuyTradePost struct {
	Symbol string  `json:"symbol"`
	Amount float64 `json:"amount"`
}

type PnlResponse struct {
	UnrealizedPL float64 `json:"unrealized_pl"`
	CumulativePL float64 `json:"cumulative_pl"`
	NetPNL       float64 `json:"net_pnl"`
}

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
