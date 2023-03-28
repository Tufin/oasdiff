package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// ContactDiff describes the changes between a pair of contact objects: https://swagger.io/specification/#contact-object
type ContactDiff struct {
	Added          bool            `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted        bool            `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	ExtensionsDiff *ExtensionsDiff `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	NameDiff       *ValueDiff      `json:"name,omitempty" yaml:"name,omitempty"`
	URLDiff        *ValueDiff      `json:"url,omitempty" yaml:"url,omitempty"`
	EmailDiff      *ValueDiff      `json:"email,omitempty" yaml:"email,omitempty"`
}

// Empty indicates whether a change was found in this element
func (diff *ContactDiff) Empty() bool {
	return diff == nil || *diff == ContactDiff{}
}

func getContactDiff(config *Config, state *state, contact1, contact2 *openapi3.Contact) *ContactDiff {
	diff := getContactDiffInternal(config, state, contact1, contact2)

	if diff.Empty() {
		return nil
	}

	return diff
}

func getContactDiffInternal(config *Config, state *state, contact1, contact2 *openapi3.Contact) *ContactDiff {

	result := ContactDiff{}

	if contact1 == nil && contact2 == nil {
		return &result
	}

	if contact1 == nil && contact2 != nil {
		result.Added = true
		return &result
	}

	if contact1 != nil && contact2 == nil {
		result.Deleted = true
		return &result
	}

	result.ExtensionsDiff = getExtensionsDiff(config, state, contact1.Extensions, contact2.Extensions)
	result.NameDiff = getValueDiff(contact1.Name, contact2.Name)
	result.URLDiff = getValueDiff(contact1.URL, contact2.URL)
	result.EmailDiff = getValueDiff(contact1.Email, contact2.Email)

	return &result
}
