package modular

type Number interface {
	Signed | Unsigned
}

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

type Unsigned interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}
