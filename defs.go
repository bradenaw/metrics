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

// Bind binds the metric definition to a *Metrics used to output, returning the metric.
func (d *CounterDef) Bind(m *Metrics) *Counter {
	return m.counter(d.name, d.tags)
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

// Bind binds the metric definition to a *Metrics used to output, returning the metric.
func (d *GaugeDef) Bind(m *Metrics) *Gauge {
	return m.gauge(d.name, d.tags)
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

// Bind binds the metric definition to a *Metrics used to output, returning the metric.
func (h *HistogramDef) Bind(m *Metrics) *Histogram {
	return &Histogram{
		m:          m,
		name:       h.name,
		tags:       h.tags,
		sampleRate: h.sampleRate,
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

// Bind binds the metric definition to a *Metrics used to output, returning the metric.
func (h *DistributionDef) Bind(m *Metrics) *Distribution {
	return &Distribution{
		m:          m,
		name:       h.name,
		tags:       h.tags,
		sampleRate: h.sampleRate,
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

type SetDef[K any] struct {
	name       string
	key        string
	tags       []string
	sampleRate float64
}

func NewSetDef[K any](
	name string,
	description string,
	unit Unit,
	sampleRate float64,
) *SetDef[K] {
	registerDef(setType, name, unit, description)
	return &SetDef[K]{
		name:       name,
		sampleRate: sampleRate,
	}
}

// Bind binds the metric definition to a *Metrics used to output, returning the metric.
func (h *SetDef[K]) Bind(m *Metrics) *Set[K] {
	return &Set[K]{
		m:          m,
		name:       h.name,
		tags:       h.tags,
		sampleRate: h.sampleRate,
	}
}

type SetDef1[K any, V0 TagValue] struct {
	name       string
	key        string
	sampleRate float64
}

func NewSetDef1[K any, V0 TagValue](
	name string,
	description string,
	unit Unit,
	key string,
	sampleRate float64,
) *SetDef1[K, V0] {
	registerDef(setType, name, unit, description)
	return &SetDef1[K, V0]{
		name:       name,
		key:        key,
		sampleRate: sampleRate,
	}
}

func (h *SetDef1[K, V0]) Values(v0 V0) *SetDef[K] {
	return &SetDef[K]{
		name: h.name,
		tags: []string{
			makeTag(h.key, tagValueString(v0)),
		},
		sampleRate: h.sampleRate,
	}
}
