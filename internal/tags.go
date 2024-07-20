package internal

import "github.com/tufin/oasdiff/checker"

func getAllTags() []string {
	return []string{"request", "response", "add", "remove", "change", "generalize", "specialize", "increase", "decrease", "set", "body", "parameters", "properties", "headers", "security", "components"}
}

// matchTags returns true if the rule matches all the tags
func matchTags(tags []string, rule checker.BackwardCompatibilityRule) bool {
	if len(tags) == 0 {
		return true
	}

	for _, tag := range tags {
		if !matchTag(tag, rule) {
			return false
		}
	}

	return true
}

func matchTag(tag string, rule checker.BackwardCompatibilityRule) bool {
	if matchLocationTag(tag, rule.Location) {
		return true
	}

	if matchActionTag(tag, rule.Action) {
		return true
	}

	if matchDirectionTag(tag, rule.Direction) {
		return true
	}

	return false
}

func matchDirectionTag(tag string, direction checker.Direction) bool {
	switch tag {
	case "request":
		return direction == checker.DirectionRequest
	case "response":
		return direction == checker.DirectionResponse
	}

	return false
}

func matchActionTag(tag string, action checker.Action) bool {
	switch tag {
	case "add":
		return action == checker.ActionAdd
	case "remove":
		return action == checker.ActionRemove
	case "change":
		return action == checker.ActionChange
	case "generalize":
		return action == checker.ActionGeneralize
	case "specialize":
		return action == checker.ActionSpecialize
	case "increase":
		return action == checker.ActionIncrease
	case "decrease":
		return action == checker.ActionDecrease
	case "set":
		return action == checker.ActionSet
	}

	return false
}

func matchLocationTag(tag string, location checker.Location) bool {
	switch tag {
	case "body":
		return location == checker.LocationBody
	case "parameters":
		return location == checker.LocationParameters
	case "properties":
		return location == checker.LocationProperties
	case "headers":
		return location == checker.LocationHeaders
	case "security":
		return location == checker.LocationSecurity
	case "components":
		return location == checker.LocationComponents
	}

	return false
}
