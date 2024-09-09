package entity

// Load -> объединение типов { V, N, F }
type Load interface{}

type V struct {
	V string
}

type N struct {
	Value int64
	Len   int
}

type F struct {
	Frame Frame
}