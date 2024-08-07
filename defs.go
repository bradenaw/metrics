package metrics

type CounterDef struct {
	name          string
	tags          tags
	allComparable bool
	ok            bool
}

func NewCounterDef(
	name string,
	description string,
	unit Unit,
) CounterDef {
	ok := registerDef(CounterType, name, description, unit, nil, nil)
	return CounterDef{
		name:          name,
		allComparable: true,
		ok:            ok,
	}
}

type GaugeDef struct {
	name          string
	tags          tags
	allComparable bool
	ok            bool
}

func NewGaugeDef(
	name string,
	description string,
	unit Unit,
) GaugeDef {
	ok := registerDef(GaugeType, name, description, unit, nil, nil)
	return GaugeDef{
		name:          name,
		allComparable: true,
		ok:            ok,
	}
}

type DistributionDef struct {
	name          string
	unit          Unit
	tags          tags
	sampleRate    float64
	allComparable bool
	ok            bool
}

func NewDistributionDef(
	name string,
	description string,
	unit Unit,
	sampleRate float64,
) DistributionDef {
	ok := registerDef(DistributionType, name, description, unit, nil, nil)
	return DistributionDef{
		name:          name,
		unit:          unit,
		sampleRate:    sampleRate,
		allComparable: true,
		ok:            ok,
	}
}

type SetDef struct {
	name          string
	tags          tags
	sampleRate    float64
	allComparable bool
	ok            bool
}

func NewSetDef(
	name string,
	description string,
	unit Unit,
	sampleRate float64,
) SetDef {
	ok := registerDef(SetType, name, description, unit, nil, nil)
	return SetDef{
		name:          name,
		sampleRate:    sampleRate,
		allComparable: true,
		ok:            ok,
	}
}
