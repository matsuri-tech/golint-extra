package example

type H struct {
	a string
	b int
}

func IncompleteStructExample_fail() H {
	// incomplete struct
	h := H{
		a: "foo",
	}

	return h
}
