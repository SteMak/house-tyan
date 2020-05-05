package modules

type HandlerXP func(string, uint64)

type XpEvents interface {
	OnXp(HandlerXP)
}

type Events struct {
	XpEvents
}
