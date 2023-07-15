package bit

type Bit int

const (
	O Bit = iota
	I
)

func (b Bit) String() string {
	switch b {
	case O:
		return "O"
	case I:
		return "I"
	default:
		panic("invalid bit")
	}
}

func (b Bit) Not() Bit {
	switch b {
	case O:
		return I
	case I:
		return O
	default:
		panic("invalid Bit")
	}
}

func And(b, b2 Bit) Bit {
	switch {
	case b == O || b2 == O:
		return O
	case b == I && b2 == I:
		return I
	default:
		panic("invalid bits")
	}
}

func Or(b, b2 Bit) Bit {
	switch {
	case b == I || b2 == I:
		return I
	case b == O && b2 == O:
		return O
	default:
		panic("invalid bits")
	}
}

func Nor(b, b2 Bit) Bit {
	return Or(b, b2).Not()
}

func BitSliceToStringSlice(bits []Bit) []string {
	var bitStrings []string
	for _, b := range bits {
		bitStrings = append(bitStrings, b.String())
	}
	return bitStrings
}
