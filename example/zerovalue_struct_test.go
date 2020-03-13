package example

type H struct {
	a string
	b int
}

func IncompleteStructExample_ignore_success() H {
	// incomplete struct
	h := H{
		// @ignore-golint-extra
		a: "foo",
	}

	return h
}

func IncompleteStructExample_fail() H {
	// incomplete struct
	h := H{
		a: "foo",
	}

	return h
}
