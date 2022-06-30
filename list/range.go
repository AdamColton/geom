package list

type Range struct {
	Start, End float64
}

func (r Range) D() float64 {
	return r.End - r.Start
}

func (r Range) Steps(n int) RangeStep {
	return RangeStep{
		Range: r,
		Step:  r.D() / float64(n),
	}
}

type RangeStep struct {
	Range
	Step float64
}

func (rs RangeStep) Steps(n int) RangeStep {
	rs.End = rs.Start + rs.Step*float64(n)
	return rs
}

func (rs RangeStep) Len() int {
	return int(rs.D() / rs.Step)
}

func (rs RangeStep) Idx(idx int) float64 {
	return rs.Start + float64(idx)*rs.Step
}
