package example

type H struct {
	a string
	b int
}

func example1() H {
	h := H{
		a: "foo",
	}

	return h
}
