package services

type LoggerService interface {
	WriteError(data string)
	WriteNotice(data string)
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

func NewServices(l LoggerService) *Services {
	return &Services{
		LoggerService: l,
	}
}
