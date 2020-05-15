package modules

// Events
var (
	Event Events
)

type HandlerXP func(string, uint64)

type XpEvents interface {
	OnXp(HandlerXP)
}

type Events struct {
	XpEvents
}
