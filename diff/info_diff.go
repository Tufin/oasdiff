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

func getInfoDiff(config *Config, state *state, info1, info2 *openapi3.Info) (*InfoDiff, error) {
	diff, err := getInfoDiffInternal(config, state, info1, info2)
	if err != nil {
		return nil, err
	}

	if diff.Empty() {
		return nil, nil
	}

	return diff, nil
}

func getInfoDiffInternal(config *Config, state *state, info1, info2 *openapi3.Info) (*InfoDiff, error) {

	result := InfoDiff{}

	if info1 == nil || info2 == nil {
		return nil, fmt.Errorf("info is nil")
	}

	result.ExtensionsDiff = getExtensionsDiff(config, state, info1.ExtensionProps, info2.ExtensionProps)
	result.TitleDiff = getValueDiff(info1.Title, info2.Title)
	result.DescriptionDiff = getValueDiffConditional(config.ExcludeDescription, info1.Description, info2.Description)
	result.TermsOfServiceDiff = getValueDiff(info1.TermsOfService, info2.TermsOfService)
	result.ContactDiff = getContactDiff(config, state, info1.Contact, info2.Contact)
	result.LicenseDiff = getLicenseDiff(config, state, info1.License, info2.License)
	result.VersionDiff = getValueDiff(info1.Version, info2.Version)

	return &result, nil
}
