package config

type Secret string

func (s Secret) String() string { return string(s) }

const (
	ZyncSecret Secret = "zync"
)
