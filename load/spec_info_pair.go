package load

type SpecInfoPair struct {
	Base     *SpecInfo
	Revision *SpecInfo
}

func (specInfoPair *SpecInfoPair) GetBaseVersion() string {
	if specInfoPair == nil {
		return "n/a"
	}
	return specInfoPair.Base.GetVersion()
}

func (specInfoPair *SpecInfoPair) GetRevisionVersion() string {
	if specInfoPair == nil {
		return "n/a"
	}

	return specInfoPair.Revision.GetVersion()
}

func NewSpecInfoPair(specInfo1, specInfo2 *SpecInfo) *SpecInfoPair {
	return &SpecInfoPair{
		Base:     specInfo1,
		Revision: specInfo2,
	}
}
