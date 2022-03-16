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
