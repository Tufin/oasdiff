package internal

type DeltaFlags struct {
	DiffFlags
	asymmetric bool
}

func (flags *DeltaFlags) getAsymmetric() bool {
	return flags.asymmetric
}
