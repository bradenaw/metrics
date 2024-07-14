package metrics

type tags struct {
	n      int
	keys   [maxTags]string
	values [maxTags]any
}

func (t tags) append(other tags) tags {
	if t.n == 0 {
		return other
	}

	out := tags{n: t.n + other.n}

	copy(out.keys[:], t.keys[:t.n])
	copy(out.keys[t.n:], other.keys[:other.n])
	copy(out.values[:], t.values[:t.n])
	copy(out.values[t.n:], other.values[:other.n])

	return out
}
