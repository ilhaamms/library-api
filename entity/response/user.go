package response

type WebResponseUser struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseUserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type CreateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
