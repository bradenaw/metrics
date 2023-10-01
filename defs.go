package metrics

type CounterDef struct {
	name string
	tags []string
	ok   bool
}

func NewCounterDef(
	name string,
	description string,
	unit Unit,
) CounterDef {
	ok := registerDef(CounterType, name, description, unit, nil, nil)
	return CounterDef{
		name: name,
		ok:   ok,
	}
}

type GaugeDef struct {
	name string
	tags []string
	ok   bool
}

func NewGaugeDef(
	name string,
	description string,
	unit Unit,
) GaugeDef {
	ok := registerDef(GaugeType, name, description, unit, nil, nil)
	return GaugeDef{
		name: name,
		ok:   ok,
	}
}

type DistributionDef struct {
	name       string
	tags       []string
	sampleRate float64
	ok         bool
}

func NewDistributionDef(
	name string,
	description string,
	unit Unit,
	sampleRate float64,
) DistributionDef {
	ok := registerDef(DistributionType, name, description, unit, nil, nil)
	return DistributionDef{
		name:       name,
		sampleRate: sampleRate,
		ok:         ok,
	}
}

type SetDef struct {
	name       string
	tags       []string
	sampleRate float64
	ok         bool
}

func NewSetDef(
	name string,
	description string,
	unit Unit,
	sampleRate float64,
) SetDef {
	ok := registerDef(SetType, name, description, unit, nil, nil)
	return SetDef{
		name:       name,
		sampleRate: sampleRate,
		ok:         ok,
	}
}
