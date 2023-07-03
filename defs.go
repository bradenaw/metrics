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
	ok := registerDef(counterType, name, unit, description)
	return CounterDef{
		name: name,
		ok:   ok,
	}
}

type CounterDef1[V0 TagValue] struct {
	name   string
	prefix []string
	keys   [1]string
	ok     bool
}

func NewCounterDef1[V0 TagValue](
	name string,
	description string,
	unit Unit,
	key string,
) CounterDef1[V0] {
	ok := registerDef(counterType, name, unit, description)
	return CounterDef1[V0]{
		name: name,
		keys: [1]string{key},
		ok:   ok,
	}
}

func (d CounterDef1[V0]) Values(v0 V0) CounterDef {
	return CounterDef{
		name: d.name,
		tags: joinStrings(d.prefix, []string{
			makeTag(d.keys[0], tagValueString(v0)),
		}),
		ok: d.ok,
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
	ok := registerDef(gaugeType, name, unit, description)
	return GaugeDef{
		name: name,
		ok:   ok,
	}
}

type GaugeDef1[V0 TagValue] struct {
	name   string
	prefix []string
	keys   [1]string
	ok     bool
}

func NewGaugeDef1[V0 TagValue](
	name string,
	description string,
	unit Unit,
	key string,
) GaugeDef1[V0] {
	ok := registerDef(gaugeType, name, unit, description)
	return GaugeDef1[V0]{
		name: name,
		keys: [1]string{key},
		ok:   ok,
	}
}

func (d GaugeDef1[V0]) Values(v0 V0) GaugeDef {
	return GaugeDef{
		name: d.name,
		tags: joinStrings(d.prefix, []string{
			makeTag(d.keys[0], tagValueString(v0)),
		}),
		ok: d.ok,
	}
}

type HistogramDef struct {
	name       string
	tags       []string
	sampleRate float64
	ok         bool
}

func NewHistogramDef(
	name string,
	description string,
	unit Unit,
	sampleRate float64,
) HistogramDef {
	ok := registerDef(histogramType, name, unit, description)
	return HistogramDef{
		name:       name,
		sampleRate: sampleRate,
		ok:         ok,
	}
}

type HistogramDef1[V0 TagValue] struct {
	name       string
	key        string
	sampleRate float64
	ok         bool
}

func NewHistogramDef1[V0 TagValue](
	name string,
	description string,
	unit Unit,
	key string,
	sampleRate float64,
) HistogramDef1[V0] {
	ok := registerDef(histogramType, name, unit, description)
	return HistogramDef1[V0]{
		name:       name,
		key:        key,
		sampleRate: sampleRate,
		ok:         ok,
	}
}

func (d HistogramDef1[V0]) Values(v0 V0) HistogramDef {
	return HistogramDef{
		name: d.name,
		tags: []string{
			makeTag(d.key, tagValueString(v0)),
		},
		sampleRate: d.sampleRate,
		ok:         d.ok,
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
	ok := registerDef(distributionType, name, unit, description)
	return DistributionDef{
		name:       name,
		sampleRate: sampleRate,
		ok:         ok,
	}
}

type DistributionDef1[V0 TagValue] struct {
	name       string
	key        string
	sampleRate float64
	ok         bool
}

func NewDistributionDef1[V0 TagValue](
	name string,
	description string,
	unit Unit,
	key string,
	sampleRate float64,
) DistributionDef1[V0] {
	ok := registerDef(distributionType, name, unit, description)
	return DistributionDef1[V0]{
		name:       name,
		key:        key,
		sampleRate: sampleRate,
		ok:         ok,
	}
}

func (d DistributionDef1[V0]) Values(v0 V0) DistributionDef {
	return DistributionDef{
		name: d.name,
		tags: []string{
			makeTag(d.key, tagValueString(v0)),
		},
		sampleRate: d.sampleRate,
		ok:         d.ok,
	}
}

type SetDef struct {
	name       string
	key        string
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
	ok := registerDef(setType, name, unit, description)
	return SetDef{
		name:       name,
		sampleRate: sampleRate,
		ok:         ok,
	}
}

type SetDef1[V0 TagValue] struct {
	name       string
	key        string
	sampleRate float64
	ok         bool
}

func NewSetDef1[V0 TagValue](
	name string,
	description string,
	unit Unit,
	key string,
	sampleRate float64,
) SetDef1[V0] {
	ok := registerDef(setType, name, unit, description)
	return SetDef1[V0]{
		name:       name,
		key:        key,
		sampleRate: sampleRate,
		ok:         ok,
	}
}

func (d SetDef1[V0]) Values(v0 V0) SetDef {
	return SetDef{
		name: d.name,
		tags: []string{
			makeTag(d.key, tagValueString(v0)),
		},
		sampleRate: d.sampleRate,
		ok:         d.ok,
	}
}
