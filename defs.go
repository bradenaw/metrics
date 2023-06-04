package metrics

type CounterDef struct {
	name string
}

func NewCounterDef(
	name string,
	unit Unit,
	description string,
) *CounterDef {
	registerDef(counterType, name, unit, description)
	return &CounterDef{
		name: name,
	}
}

// Bind binds a set of tag values to this definition and returns the resulting Counter.
func (d *CounterDef) Bind(m *Metrics) *Counter {
	return m.counter(d.name, nil)
}

type CounterDef1[V0 TagValue] struct {
	name   string
	prefix []string
	key    string
}

func NewCounterDef1[V0 TagValue](
	name string,
	unit Unit,
	key string,
	description string,
) *CounterDef1[V0] {
	registerDef(counterType, name, unit, description)
	return &CounterDef1[V0]{
		name: name,
		key:  key,
	}
}

// Bind binds a set of tag values to this definition and returns the resulting Counter.
func (d *CounterDef1[V0]) Bind(m *Metrics, v0 V0) *Counter {
	return m.counter(d.name, joinStrings(d.prefix, []string{
		makeTag(d.key, tagValueString(v0)),
	}))
}

type GaugeDef struct {
	name string
}

func NewGaugeDef(
	name string,
	unit Unit,
	description string,
) *GaugeDef {
	registerDef(gaugeType, name, unit, description)
	return &GaugeDef{
		name: name,
	}
}

func (d *GaugeDef) Bind(m *Metrics) *Gauge {
	return m.gauge(d.name, nil)
}

type GaugeDef1[V0 TagValue] struct {
	name string
	key  string
}

func NewGaugeDef1[V0 TagValue](
	name string,
	unit Unit,
	key string,
	description string,
) *GaugeDef1[V0] {
	registerDef(gaugeType, name, unit, description)
	return &GaugeDef1[V0]{
		name: name,
		key:  key,
	}
}

func (d *GaugeDef1[V0]) Bind(m *Metrics, v0 V0) *Gauge {
	return m.gauge(d.name, []string{
		makeTag(d.key, tagValueString(v0)),
	})
}

type HistogramDef struct {
	name       string
	key        string
	sampleRate float64
}

func NewHistogramDef(
	name string,
	unit Unit,
	key string,
	description string,
	sampleRate float64,
) *HistogramDef {
	registerDef(histogramType, name, unit, description)
	return &HistogramDef{
		name:       name,
		key:        key,
		sampleRate: sampleRate,
	}
}

func (h *HistogramDef) Bind(m *Metrics) *Histogram {
	return &Histogram{
		m:          m,
		name:       h.name,
		tags:       nil,
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
	unit Unit,
	key string,
	description string,
	sampleRate float64,
) *HistogramDef1[V0] {
	registerDef(histogramType, name, unit, description)
	return &HistogramDef1[V0]{
		name:       name,
		key:        key,
		sampleRate: sampleRate,
	}
}

func (h *HistogramDef1[V0]) Bind(m *Metrics, v0 V0) *Histogram {
	return &Histogram{
		m:    m,
		name: h.name,
		tags: []string{
			makeTag(h.key, tagValueString(v0)),
		},
		sampleRate: h.sampleRate,
	}
}

type DistributionDef struct {
	name       string
	key        string
	sampleRate float64
}

func NewDistributionDef(
	name string,
	unit Unit,
	key string,
	description string,
	sampleRate float64,
) *DistributionDef {
	registerDef(distributionType, name, unit, description)
	return &DistributionDef{
		name:       name,
		key:        key,
		sampleRate: sampleRate,
	}
}

func (h *DistributionDef) Bind(m *Metrics) *Distribution {
	return &Distribution{
		m:          m,
		name:       h.name,
		tags:       nil,
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
	unit Unit,
	key string,
	description string,
	sampleRate float64,
) *DistributionDef1[V0] {
	registerDef(distributionType, name, unit, description)
	return &DistributionDef1[V0]{
		name:       name,
		key:        key,
		sampleRate: sampleRate,
	}
}

func (h *DistributionDef1[V0]) Bind(m *Metrics, v0 V0) *Distribution {
	return &Distribution{
		m:    m,
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
	sampleRate float64
}

func NewSetDef[K any](
	name string,
	unit Unit,
	key string,
	description string,
	sampleRate float64,
) *SetDef[K] {
	registerDef(setType, name, unit, description)
	return &SetDef[K]{
		name:       name,
		key:        key,
		sampleRate: sampleRate,
	}
}

func (h *SetDef[K]) Bind(m *Metrics) *Set[K] {
	return &Set[K]{
		m:          m,
		name:       h.name,
		tags:       nil,
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
	unit Unit,
	key string,
	description string,
	sampleRate float64,
) *SetDef1[K, V0] {
	registerDef(setType, name, unit, description)
	return &SetDef1[K, V0]{
		name:       name,
		key:        key,
		sampleRate: sampleRate,
	}
}

func (h *SetDef1[K, V0]) Bind(m *Metrics, v0 V0) *Set[K] {
	return &Set[K]{
		m:    m,
		name: h.name,
		tags: []string{
			makeTag(h.key, tagValueString(v0)),
		},
		sampleRate: h.sampleRate,
	}
}
