package services

type LoggerService interface {
	WriteError(data string)
	WriteNotice(data string)
}

type StoreService interface {
}

type AuthService interface {
	GetTokenPair() (*tokenPair, error)
	//RefreshTokens(string) (string, string, error)
}

type Services struct {
	LoggerService
	StoreService
	AuthService
}

func NewServices(l LoggerService, a AuthService) *Services {
	return &Services{
		LoggerService: l,
		AuthService:   a,
	}
}
