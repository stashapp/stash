package exp

type (
	LockStrength int
	WaitOption   int
	Lock         interface {
		Strength() LockStrength
		WaitOption() WaitOption
		Of() []IdentifierExpression
	}
	lock struct {
		strength   LockStrength
		waitOption WaitOption
		of         []IdentifierExpression
	}
)

const (
	ForNolock LockStrength = iota
	ForUpdate
	ForNoKeyUpdate
	ForShare
	ForKeyShare

	Wait WaitOption = iota
	NoWait
	SkipLocked
)

func NewLock(strength LockStrength, option WaitOption, of ...IdentifierExpression) Lock {
	return lock{
		strength:   strength,
		waitOption: option,
		of:         of,
	}
}

func (l lock) Strength() LockStrength {
	return l.strength
}

func (l lock) WaitOption() WaitOption {
	return l.waitOption
}

func (l lock) Of() []IdentifierExpression {
	return l.of
}
