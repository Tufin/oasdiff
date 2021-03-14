package diff_test

// func TestDiff_Unresolved(t *testing.T) {
// 	loader := openapi3.NewSwaggerLoader()
// 	s1, err := loader.LoadSwaggerFromFile("../data/home-iot-api-1.yaml")
// 	require.NoError(t, err)
// 	s2, err := loader.LoadSwaggerFromFile("../data/home-iot-api-2.yaml")
// 	require.NoError(t, err)

// 	s1.Paths["/devices"].Post.RequestBody.Value.Content["application/json"].Schema.Value = nil
// 	diff.Get(diff.NewConfig(), s1, s2)

// 	s1.Paths["/devices"].Post.RequestBody.Value.Content["application/json"].Schema = nil
// 	diff.Get(diff.NewConfig(), s1, s2)

// 	s1.Paths["/devices"].Post.RequestBody.Value.Content["application/json"] = nil
// 	diff.Get(diff.NewConfig(), s1, s2)

// 	s1.Paths["/devices"].Post.RequestBody.Value.Content = nil
// 	diff.Get(diff.NewConfig(), s1, s2)

// 	s1.Paths["/devices"].Post.RequestBody = nil
// 	diff.Get(diff.NewConfig(), s1, s2)

// 	s1.Paths["/devices"].Post = nil
// 	diff.Get(diff.NewConfig(), s1, s2)

// 	s1.Paths["/devices"] = nil
// 	diff.Get(diff.NewConfig(), s1, s2)

// 	s1.Paths = nil
// 	diff.Get(diff.NewConfig(), s1, s2)

// 	s1 = nil
// 	diff.Get(diff.NewConfig(), s1, s2)
// }
