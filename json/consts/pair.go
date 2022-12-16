package consts

type Pair[T, R any] struct {
	First  T
	Second R
}

func (p *Pair[T, R]) GetFirst() T {
	return p.First
}

func (p *Pair[T, R]) GetSecond() R {
	return p.Second
}

func (p *Pair[T, R]) UpdateSecond(data R) {
	p.Second = data
}
