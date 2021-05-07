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

func NoFieldsExample_success() {
	_ = H{}
}

// An example using multiple functions

func Ex1() string {
	// Inner function struct
	type A struct {
		a string
	}

	a := A{
		a: "foo",
	}

	return a.a
}

func Ex2() int {
	// Inner function struct
	type A struct {
		b int
	}

	a := A{
		b: 200,
	}

	return a.b
}
