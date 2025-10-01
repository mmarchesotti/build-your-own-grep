package predefinedclass

type PredefinedClass int

const (
	ClassDigit PredefinedClass = iota
	ClassAlphanumeric
	ClassWhitespace
)
