package services

type LoggerService interface {
	WriteError(data string)
	WriteNotice(data string)
}

type StoreService interface {
	WriteTokensInfo(token string) error
}

type AuthService interface {
	CreateTokenPair() (*tokenPair, error)
	RefreshTokens(rToken string) (*tokenPair, error)
	ValidateToken(token string, refresh bool) (string, error)
}

type Services struct {
	LoggerService
	StoreService
	AuthService
}

func NewServices(l LoggerService, a AuthService, s StoreService) *Services {
	return &Services{
		LoggerService: l,
		AuthService:   a,
		StoreService:  s,
	}
}
