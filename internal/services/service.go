package services

type LoggerService interface {
	Println(v ...interface{})
}

type StoreService interface {
}

type AuthService interface {
	GenerateTokens() (string, string, error)
	RefreshTokens(string) (string, string, error)
}

type Services struct {
	LoggerService
	StoreService
	AuthService
}
