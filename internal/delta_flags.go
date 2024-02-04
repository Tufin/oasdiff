package internal

type DeltaFlags struct {
	DiffFlags
	asymmetric bool
}

func (flags *DeltaFlags) getAsymmetric() bool {
	return flags.asymmetric
}

func (flags *DeltaFlags) refComposed() *bool {
	return &flags.composed
}

func (flags *DeltaFlags) refExcludeElements() *[]string {
	return &flags.excludeElements
}

func (flags *DeltaFlags) refMatchPath() *string {
	return &flags.matchPath
}

func (flags *DeltaFlags) refFilterExtension() *string {
	return &flags.filterExtension
}

func (flags *DeltaFlags) refCircularReferenceCounter() *int {
	return &flags.circularReferenceCounter
}

func (flags *DeltaFlags) refPrefixBase() *string {
	return &flags.prefixBase
}

func (flags *DeltaFlags) refPrefixRevision() *string {
	return &flags.prefixRevision
}

func (flags *DeltaFlags) refStripPrefixBase() *string {
	return &flags.stripPrefixBase
}

func (flags *DeltaFlags) refStripPrefixRevision() *string {
	return &flags.stripPrefixRevision
}

func (flags *DeltaFlags) refIncludePathParams() *bool {
	return &flags.includePathParams
}

func (flags *DeltaFlags) refFlattenAllOf() *bool {
	return &flags.flattenAllOf
}

func (flags *DeltaFlags) refFlattenParams() *bool {
	return &flags.flattenParams
}
