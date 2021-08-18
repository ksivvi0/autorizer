package services

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
