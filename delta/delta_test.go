package delta_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/delta"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/utils"
)

func TestEmpty(t *testing.T) {
	d := &diff.Diff{}
	require.Equal(t, 0.0, delta.Get(false, d))
}

func TestEndpointAdded(t *testing.T) {
	d := &diff.Diff{
		EndpointsDiff: &diff.EndpointsDiff{
			Added: diff.Endpoints{
				diff.Endpoint{
					Method: "GET",
					Path:   "/test",
				},
			},
			Unchanged: diff.Endpoints{
				diff.Endpoint{
					Method: "POST",
					Path:   "/test",
				},
			},
		},
	}

	require.Equal(t, 0.5, delta.Get(false, d))
}

func TestEndpointDeletedAsym(t *testing.T) {
	d := &diff.Diff{
		EndpointsDiff: &diff.EndpointsDiff{
			Deleted: diff.Endpoints{
				diff.Endpoint{
					Method: "GET",
					Path:   "/test",
				},
			},
			Unchanged: diff.Endpoints{
				diff.Endpoint{
					Method: "POST",
					Path:   "/test",
				},
			},
		},
	}

	require.Equal(t, 0.5, delta.Get(true, d))
}

func TestEndpointAddedAndDeleted(t *testing.T) {
	d := &diff.Diff{
		EndpointsDiff: &diff.EndpointsDiff{
			Added: diff.Endpoints{
				diff.Endpoint{
					Method: "GET",
					Path:   "/test",
				},
			},
			Deleted: diff.Endpoints{
				diff.Endpoint{
					Method: "POST",
					Path:   "/test",
				},
			},
		},
	}

	require.Equal(t, 1.0, delta.Get(false, d))
}

func TestParameters(t *testing.T) {
	d := &diff.Diff{
		EndpointsDiff: &diff.EndpointsDiff{
			Modified: diff.ModifiedEndpoints{
				diff.Endpoint{
					Method: "GET",
					Path:   "/test",
				}: &diff.MethodDiff{
					ParametersDiff: &diff.ParametersDiffByLocation{
						Deleted: diff.ParamNamesByLocation{
							"query": utils.StringList{"a"},
						},
					},
				},
			},
		},
	}

	require.Equal(t, 0.5, delta.Get(true, d))
}

func TestResponses(t *testing.T) {
	d := &diff.Diff{
		EndpointsDiff: &diff.EndpointsDiff{
			Modified: diff.ModifiedEndpoints{
				diff.Endpoint{
					Method: "GET",
					Path:   "/test",
				}: &diff.MethodDiff{
					ResponsesDiff: &diff.ResponsesDiff{
						Added:   utils.StringList{"201"},
						Deleted: utils.StringList{"200"},
					},
				},
			},
		},
	}

	require.Equal(t, 0.5, delta.Get(false, d))
}

func TestSchema(t *testing.T) {
	d := &diff.Diff{
		EndpointsDiff: &diff.EndpointsDiff{
			Modified: diff.ModifiedEndpoints{
				diff.Endpoint{
					Method: "GET",
					Path:   "/test",
				}: &diff.MethodDiff{
					ParametersDiff: &diff.ParametersDiffByLocation{
						Modified: diff.ParamDiffByLocation{
							"query": diff.ParamDiffs{
								"a": &diff.ParameterDiff{
									SchemaDiff: &diff.SchemaDiff{
										TypeDiff: &diff.ValueDiff{
											From: "integer",
											To:   "string",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	require.Equal(t, 0.25, delta.Get(false, d))
}

func TestSymmetric(t *testing.T) {
	specs := utils.StringList{"../data/simple.yaml", "../data/simple1.yaml", "../data/simple2.yaml", "../data/simple3.yaml", "../data/simple4.yaml", "../data/simple5.yaml"}
	specPairs := specs.CartesianProduct(specs)

	loader := openapi3.NewLoader()
	for _, pair := range specPairs {
		s1, err := loader.LoadFromFile(pair.X)
		require.NoError(t, err)

		s2, err := loader.LoadFromFile(pair.Y)
		require.NoError(t, err)

		d1, err := diff.Get(diff.NewConfig(), s1, s2)
		require.NoError(t, err)

		d2, err := diff.Get(diff.NewConfig(), s2, s1)
		require.NoError(t, err)

		require.Equal(t, delta.Get(false, d1), delta.Get(false, d2), pair)
	}
}

func TestAsymmetric(t *testing.T) {
	specs := utils.StringList{"../data/simple.yaml", "../data/simple1.yaml", "../data/simple2.yaml", "../data/simple3.yaml", "../data/simple4.yaml", "../data/simple5.yaml"}
	specPairs := specs.CartesianProduct(specs)

	loader := openapi3.NewLoader()
	for _, pair := range specPairs {
		s1, err := loader.LoadFromFile(pair.X)
		require.NoError(t, err)

		s2, err := loader.LoadFromFile(pair.Y)
		require.NoError(t, err)

		d1, err := diff.Get(diff.NewConfig(), s1, s2)
		require.NoError(t, err)
		asymmetric1 := delta.Get(true, d1)

		d2, err := diff.Get(diff.NewConfig(), s2, s1)
		require.NoError(t, err)
		asymmetric2 := delta.Get(true, d2)

		symmetric := delta.Get(false, d2)

		require.Equal(t, asymmetric1+asymmetric2, symmetric, pair)
	}
}
