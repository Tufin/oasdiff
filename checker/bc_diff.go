package checker

import "github.com/tufin/oasdiff/diff"

type BCDiff struct {
	diff.Diff
}

func (d *BCDiff) AddModifiedOperation(path string, operation string) *diff.MethodDiff {
	pathDiff := d.AddModifiedPath(path)
	if pathDiff.OperationsDiff == nil {
		pathDiff.OperationsDiff = &diff.OperationsDiff{}
	}
	if pathDiff.OperationsDiff.Modified == nil {
		pathDiff.OperationsDiff.Modified = make(diff.ModifiedOperations)
	}
	if pathDiff.OperationsDiff.Modified[operation] == nil {
		pathDiff.OperationsDiff.Modified[operation] = &diff.MethodDiff{}
	}
	return pathDiff.OperationsDiff.Modified[operation]
}

func (d *BCDiff) AddModifiedPath(path string) *diff.PathDiff {
	if d.PathsDiff == nil {
		d.PathsDiff = &diff.PathsDiff{}
	}
	if d.PathsDiff.Modified == nil {
		d.PathsDiff.Modified = make(diff.ModifiedPaths)
	}
	if d.PathsDiff.Modified[path] == nil {
		d.PathsDiff.Modified[path] = &diff.PathDiff{}
	}
	return d.PathsDiff.Modified[path]
}

func (diffBC *BCDiff) AddModifiedParameter(path string, operation string, paramLocation string, paramName string) *diff.ParameterDiff {
	opDiff := diffBC.AddModifiedOperation(path, operation)
	if opDiff.ParametersDiff == nil {
		opDiff.ParametersDiff = &diff.ParametersDiffByLocation{}
	}
	if opDiff.ParametersDiff.Modified == nil {
		opDiff.ParametersDiff.Modified = make(diff.ParamDiffByLocation)
	}
	if opDiff.ParametersDiff.Modified[paramLocation] == nil {
		opDiff.ParametersDiff.Modified[paramLocation] = make(diff.ParamDiffs)
	}
	if opDiff.ParametersDiff.Modified[paramLocation][paramName] == nil {
		opDiff.ParametersDiff.Modified[paramLocation][paramName] = &diff.ParameterDiff{}
	}
	return opDiff.ParametersDiff.Modified[paramLocation][paramName]
}

func (diffBC *BCDiff) AddRequestPropertiesDiff(path string, operation string, mediaType string) *diff.SchemasDiff {
	opDiff := diffBC.AddModifiedOperation(path, operation)
	if opDiff.RequestBodyDiff == nil {
		opDiff.RequestBodyDiff = &diff.RequestBodyDiff{}
	}
	if opDiff.RequestBodyDiff.ContentDiff == nil {
		opDiff.RequestBodyDiff.ContentDiff = &diff.ContentDiff{}
	}
	if opDiff.RequestBodyDiff.ContentDiff.MediaTypeModified == nil {
		opDiff.RequestBodyDiff.ContentDiff.MediaTypeModified = make(diff.ModifiedMediaTypes)
	}
	if opDiff.RequestBodyDiff.ContentDiff.MediaTypeModified[mediaType] == nil {
		opDiff.RequestBodyDiff.ContentDiff.MediaTypeModified[mediaType] = &diff.MediaTypeDiff{}
	}
	mediaTypeBCDiff := opDiff.RequestBodyDiff.ContentDiff.MediaTypeModified[mediaType]
	if mediaTypeBCDiff.SchemaDiff == nil {
		mediaTypeBCDiff.SchemaDiff = &diff.SchemaDiff{}
	}
	if mediaTypeBCDiff.SchemaDiff.PropertiesDiff == nil {
		mediaTypeBCDiff.SchemaDiff.PropertiesDiff = &diff.SchemasDiff{}
	}
	return mediaTypeBCDiff.SchemaDiff.PropertiesDiff
}
