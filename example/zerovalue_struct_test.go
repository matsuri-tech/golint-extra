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

func MapExample_success() interface{} {
	type Map = map[string]string

	m := Map{
		"fooo": "bar",
	}

	return m
}

// An example using multiple functions

func Ex1_success() string {
	// Inner function struct
	type A struct {
		a string
	}

	a := A{
		a: "foo",
	}

	return a.a
}

func Ex2_success() int {
	// Inner function struct
	type A struct {
		b int
	}

	a := A{
		b: 200,
	}

	return a.b
}
