package services

type LoggerService interface {
	WriteError(data string)
	WriteNotice(data string)
}

type StoreService interface {
}

type AuthService interface {
	GetTokenPair() (*tokenPair, error)
	RefreshTokens(rToken string) (*tokenPair, error)
	ParseToken(token string) (string, error)
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
