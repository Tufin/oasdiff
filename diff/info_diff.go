package diff

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

// InfoDiff describes the changes between a pair of info objects: https://swagger.io/specification/#info-object
type InfoDiff struct {
	ExtensionsDiff     *ExtensionsDiff `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	TitleDiff          *ValueDiff      `json:"title,omitempty" yaml:"title,omitempty"`
	DescriptionDiff    *ValueDiff      `json:"description,omitempty" yaml:"description,omitempty"`
	TermsOfServiceDiff *ValueDiff      `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	ContactDiff        *ContactDiff    `json:"contact,omitempty" yaml:"contact,omitempty"`
	LicenseDiff        *LicenseDiff    `json:"license,omitempty" yaml:"license,omitempty"`
	VersionDiff        *ValueDiff      `json:"version,omitempty" yaml:"version,omitempty"`
}

// Empty indicates whether a change was found in this element
func (diff *InfoDiff) Empty() bool {
	return diff == nil || *diff == InfoDiff{}
}

// Breaking indicates whether this element includes a breaking change
func (diff *InfoDiff) Breaking() bool {
	return false
}

func getInfoDiff(config *Config, info1, info2 *openapi3.Info) (*InfoDiff, error) {
	diff, err := getInfoDiffInternal(config, info1, info2)
	if err != nil {
		return nil, err
	}

	if diff.Empty() {
		return nil, nil
	}

	if config.BreakingOnly && !diff.Breaking() {
		return nil, nil
	}

	return diff, nil
}

func getInfoDiffInternal(config *Config, info1, info2 *openapi3.Info) (*InfoDiff, error) {

	result := InfoDiff{}

	if info1 == nil || info2 == nil {
		return nil, fmt.Errorf("info is nil")
	}

	result.ExtensionsDiff = getExtensionsDiff(config, info1.ExtensionProps, info2.ExtensionProps)
	result.TitleDiff = getValueDiff(config, false, info1.Title, info2.Title)
	result.DescriptionDiff = getValueDiffConditional(config, false, config.ExcludeDescription, info1.Description, info2.Description)
	result.TermsOfServiceDiff = getValueDiff(config, false, info1.TermsOfService, info2.TermsOfService)
	result.ContactDiff = getContactDiff(config, info1.Contact, info2.Contact)
	result.LicenseDiff = getLicenseDiff(config, info1.License, info2.License)
	result.VersionDiff = getValueDiff(config, false, info1.Version, info2.Version)

	return &result, nil
}
