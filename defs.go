package metrics

type CounterDef struct {
	name string
	tags []string
}

func NewCounterDef(
	name string,
	description string,
	unit Unit,
) *CounterDef {
	registerDef(counterType, name, unit, description)
	return &CounterDef{
		name: name,
	}
}

type CounterDef1[V0 TagValue] struct {
	name   string
	prefix []string
	keys   [1]string
}

func NewCounterDef1[V0 TagValue](
	name string,
	description string,
	unit Unit,
	key string,
) *CounterDef1[V0] {
	registerDef(counterType, name, unit, description)
	return &CounterDef1[V0]{
		name: name,
		keys: [1]string{key},
	}
}

func (d *CounterDef1[V0]) Values(v0 V0) *CounterDef {
	return &CounterDef{
		name: d.name,
		tags: joinStrings(d.prefix, []string{
			makeTag(d.keys[0], tagValueString(v0)),
		}),
	}
}

type GaugeDef struct {
	name string
	tags []string
}

func NewGaugeDef(
	name string,
	description string,
	unit Unit,
) *GaugeDef {
	registerDef(gaugeType, name, unit, description)
	return &GaugeDef{
		name: name,
	}
}

type GaugeDef1[V0 TagValue] struct {
	name   string
	prefix []string
	keys   [1]string
}

func NewGaugeDef1[V0 TagValue](
	name string,
	description string,
	unit Unit,
	key string,
) *GaugeDef1[V0] {
	registerDef(gaugeType, name, unit, description)
	return &GaugeDef1[V0]{
		name: name,
		keys: [1]string{key},
	}
}

func (d *GaugeDef1[V0]) Values(v0 V0) *GaugeDef {
	return &GaugeDef{
		name: d.name,
		tags: joinStrings(d.prefix, []string{
			makeTag(d.keys[0], tagValueString(v0)),
		}),
	}
}

type HistogramDef struct {
	name       string
	tags       []string
	sampleRate float64
}

func NewHistogramDef(
	name string,
	description string,
	unit Unit,
	sampleRate float64,
) *HistogramDef {
	registerDef(histogramType, name, unit, description)
	return &HistogramDef{
		name:       name,
		sampleRate: sampleRate,
	}
}

type HistogramDef1[V0 TagValue] struct {
	name       string
	key        string
	sampleRate float64
}

func NewHistogramDef1[V0 TagValue](
	name string,
	description string,
	unit Unit,
	key string,
	sampleRate float64,
) *HistogramDef1[V0] {
	registerDef(histogramType, name, unit, description)
	return &HistogramDef1[V0]{
		name:       name,
		key:        key,
		sampleRate: sampleRate,
	}
}

func (h *HistogramDef1[V0]) Values(v0 V0) *HistogramDef {
	return &HistogramDef{
		name: h.name,
		tags: []string{
			makeTag(h.key, tagValueString(v0)),
		},
		sampleRate: h.sampleRate,
	}
}

type DistributionDef struct {
	name       string
	tags       []string
	sampleRate float64
}

func NewDistributionDef(
	name string,
	description string,
	unit Unit,
	sampleRate float64,
) *DistributionDef {
	registerDef(distributionType, name, unit, description)
	return &DistributionDef{
		name:       name,
		sampleRate: sampleRate,
	}
}

type DistributionDef1[V0 TagValue] struct {
	name       string
	key        string
	sampleRate float64
}

func NewDistributionDef1[V0 TagValue](
	name string,
	description string,
	unit Unit,
	key string,
	sampleRate float64,
) *DistributionDef1[V0] {
	registerDef(distributionType, name, unit, description)
	return &DistributionDef1[V0]{
		name:       name,
		key:        key,
		sampleRate: sampleRate,
	}
}

func (h *DistributionDef1[V0]) Values(v0 V0) *DistributionDef {
	return &DistributionDef{
		name: h.name,
		tags: []string{
			makeTag(h.key, tagValueString(v0)),
		},
		sampleRate: h.sampleRate,
	}
}

type SetDef struct {
	name       string
	key        string
	tags       []string
	sampleRate float64
}

func NewSetDef(
	name string,
	description string,
	unit Unit,
	sampleRate float64,
) *SetDef {
	registerDef(setType, name, unit, description)
	return &SetDef{
		name:       name,
		sampleRate: sampleRate,
	}
}

type SetDef1[V0 TagValue] struct {
	name       string
	key        string
	sampleRate float64
}

func NewSetDef1[V0 TagValue](
	name string,
	description string,
	unit Unit,
	key string,
	sampleRate float64,
) *SetDef1[V0] {
	registerDef(setType, name, unit, description)
	return &SetDef1[V0]{
		name:       name,
		key:        key,
		sampleRate: sampleRate,
	}
}

func (h *SetDef1[V0]) Values(v0 V0) *SetDef {
	return &SetDef{
		name: h.name,
		tags: []string{
			makeTag(h.key, tagValueString(v0)),
		},
		sampleRate: h.sampleRate,
	}
}
