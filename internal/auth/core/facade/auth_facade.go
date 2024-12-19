package facade

type IAuthFacade interface {
}

type authFacade struct {
}

func NewAuthFacade() IAuthFacade {
	return &authFacade{}
}
