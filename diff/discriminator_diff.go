package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// DiscriminatorDiff describes the changes between a pair of discriminator objects: https://swagger.io/specification/#discriminator-object
type DiscriminatorDiff struct {
	Added            bool            `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted          bool            `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	ExtensionsDiff   *ExtensionsDiff `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	PropertyNameDiff *ValueDiff      `json:"propertyName,omitempty" yaml:"propertyName,omitempty"`
	MappingDiff      *StringMapDiff  `json:"mapping,omitempty" yaml:"mapping,omitempty"`
}

// Empty indicates whether a change was found in this element
func (diff *DiscriminatorDiff) Empty() bool {
	return diff == nil || *diff == DiscriminatorDiff{}
}

func (diff *DiscriminatorDiff) removeNonBreaking() {

	if diff.Empty() {
		return
	}

	diff.ExtensionsDiff = nil
}

func newDiscriminatorDiff() *DiscriminatorDiff {
	return &DiscriminatorDiff{}

}

func getDiscriminatorDiff(config *Config, discriminator1, discriminator2 *openapi3.Discriminator) *DiscriminatorDiff {
	diff := getDiscriminatorDiffInternal(config, discriminator1, discriminator2)

	if config.BreakingOnly {
		diff.removeNonBreaking()
	}

	if diff.Empty() {
		return nil
	}

	return diff
}

func getDiscriminatorDiffInternal(config *Config, discriminator1, discriminator2 *openapi3.Discriminator) *DiscriminatorDiff {

	result := newDiscriminatorDiff()

	if discriminator1 == nil && discriminator2 == nil {
		return result
	}

	if discriminator1 == nil && discriminator2 != nil {
		result.Added = true
		return result
	}

	if discriminator1 != nil && discriminator2 == nil {
		result.Deleted = true
		return result
	}

	result.ExtensionsDiff = getExtensionsDiff(config, discriminator1.ExtensionProps, discriminator2.ExtensionProps)
	result.PropertyNameDiff = getValueDiff(config, discriminator1.PropertyName, discriminator2.PropertyName)
	result.MappingDiff = getStringMapDiff(config, discriminator1.Mapping, discriminator2.Mapping)

	return result
}
