package diff

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// SecurityRequirementsDiff is a diff between two sets of security requirement objects: https://swagger.io/specification/#security-requirement-object
type SecurityRequirementsDiff struct {
	Added   StringList `json:"added,omitempty"`
	Deleted StringList `json:"deleted,omitempty"`
	// TODO: handle Modified Security Requirements
}

func (diff *SecurityRequirementsDiff) empty() bool {
	if diff == nil {
		return true
	}

	return len(diff.Added) == 0 &&
		len(diff.Deleted) == 0
}

func newSecurityRequirementsDiff() *SecurityRequirementsDiff {
	return &SecurityRequirementsDiff{
		Added:   StringList{},
		Deleted: StringList{},
	}
}

func getSecurityRequirementsDiff(config *Config, securityRequirements1, securityRequirements2 *openapi3.SecurityRequirements) *SecurityRequirementsDiff {
	diff := getSecurityRequirementsDiffInternal(config, securityRequirements1, securityRequirements2)
	if diff.empty() {
		return nil
	}
	return diff
}

func getSecurityRequirementsDiffInternal(config *Config, securityRequirements1, securityRequirements2 *openapi3.SecurityRequirements) *SecurityRequirementsDiff {

	result := newSecurityRequirementsDiff()

	if securityRequirements1 != nil {
		for _, securityRequirement1 := range *securityRequirements1 {
			if findSecurityRequirement(securityRequirement1, securityRequirements2) {
				// TODO: handle modification
			} else {
				result.Deleted = append(result.Deleted, getSecurityRequirementID(securityRequirement1))
			}
		}
	}

	if securityRequirements2 != nil {
		for _, securityRequirement2 := range *securityRequirements2 {
			if !findSecurityRequirement(securityRequirement2, securityRequirements1) {
				result.Added = append(result.Added, getSecurityRequirementID(securityRequirement2))
			}
		}
	}

	return result
}

func findSecurityRequirement(securityRequirement1 openapi3.SecurityRequirement, securityRequirements2 *openapi3.SecurityRequirements) bool {
	if securityRequirements2 == nil {
		return false
	}

	securitySchemes1 := getSecuritySchemes(securityRequirement1)
	for _, securityRequirement2 := range *securityRequirements2 {
		securitySchemes2 := getSecuritySchemes(securityRequirement2)
		if securitySchemes1.equals(securitySchemes2) {
			return true
		}
	}
	return false
}

func getSecuritySchemes(securityRequirement openapi3.SecurityRequirement) StringSet {
	result := StringSet{}
	for name := range securityRequirement {
		result.add(name)
	}
	return result
}

func getSecurityRequirementID(securityRequirement openapi3.SecurityRequirement) string {
	results := make([]string, len(securityRequirement))
	i := 0
	for name := range securityRequirement {
		results[i] = name
		i++
	}
	return strings.Join(results, " AND ")
}

func (diff *SecurityRequirementsDiff) getSummary() *SummaryDetails {
	return &SummaryDetails{
		Added:   len(diff.Added),
		Deleted: len(diff.Deleted),
	}
}
