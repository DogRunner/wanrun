package facade

type IDogrunmgFacade interface {
}

type dogrunmgFacade struct {
}

func NewDogrunmgFacade() IDogrunmgFacade {
	return &dogrunmgFacade{}
}
