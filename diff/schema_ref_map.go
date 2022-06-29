package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type refWithCount struct {
	schemaRef *openapi3.SchemaRef
	count     int
}

type schemaRefMap map[string]*refWithCount

func (schemaRefMap schemaRefMap) add(schemaRef *openapi3.SchemaRef) {
	if val, found := schemaRefMap[schemaRef.Ref]; found {
		val.count++
	} else {
		schemaRefMap[schemaRef.Ref] = &refWithCount{
			schemaRef: schemaRef,
			count:     1,
		}
	}
}

func (schemaRefMap schemaRefMap) delete(ref string) {
	if val, found := schemaRefMap[ref]; found {
		val.count--
		if val.count == 0 {
			delete(schemaRefMap, ref)
		}
	}
}

func (schemaRefMapOrig schemaRefMap) copy() schemaRefMap {
	result := schemaRefMap{}
	for k, v := range schemaRefMapOrig {
		result[k] = &refWithCount{
			schemaRef: v.schemaRef,
			count:     v.count,
		}
	}
	return result
}

func toSchemaRefMap(schemaRefs openapi3.SchemaRefs, filter schemaRefsFilter) schemaRefMap {
	result := schemaRefMap{}
	for _, schemaRef := range schemaRefs {
		if filter(schemaRef) {
			result.add(schemaRef)
		}
	}
	return result
}
