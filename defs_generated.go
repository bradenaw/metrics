package metrics

// generated by `go run ./gen_defs > defs_generated.go && gofmt -w defs_generated.go`

type CounterDef2[V0 TagValue, V1 TagValue] struct {
	name   string
	prefix []string
	keys   [2]string
}

func NewCounterDef2[V0 TagValue, V1 TagValue](
	name string,
	description string,
	unit Unit,
	keys [2]string,
) *CounterDef2[V0, V1] {
	registerDef(counterType, name, unit, description)
	return &CounterDef2[V0, V1]{
		name: name,
		keys: keys,
	}
}

func (h *CounterDef2[V0, V1]) Bind(m *Metrics, v0 V0, v1 V1) *Counter {
	return m.counter(h.name, joinStrings(h.prefix, []string{

		makeTag(h.keys[0], tagValueString(v0)),

		makeTag(h.keys[1], tagValueString(v1)),
	}))
}

func (h *CounterDef2[V0, V1]) BindPrefix1(v0 V0) *CounterDef1[V1] {
	return &CounterDef1[V1]{
		name: h.name,
		prefix: []string{

			makeTag(h.keys[0], tagValueString(v0)),
		},
		keys: *((*[1]string)(h.keys[1:])),
	}
}

type CounterDef3[V0 TagValue, V1 TagValue, V2 TagValue] struct {
	name   string
	prefix []string
	keys   [3]string
}

func NewCounterDef3[V0 TagValue, V1 TagValue, V2 TagValue](
	name string,
	description string,
	unit Unit,
	keys [3]string,
) *CounterDef3[V0, V1, V2] {
	registerDef(counterType, name, unit, description)
	return &CounterDef3[V0, V1, V2]{
		name: name,
		keys: keys,
	}
}

func (h *CounterDef3[V0, V1, V2]) Bind(m *Metrics, v0 V0, v1 V1, v2 V2) *Counter {
	return m.counter(h.name, joinStrings(h.prefix, []string{

		makeTag(h.keys[0], tagValueString(v0)),

		makeTag(h.keys[1], tagValueString(v1)),

		makeTag(h.keys[2], tagValueString(v2)),
	}))
}

func (h *CounterDef3[V0, V1, V2]) BindPrefix1(v0 V0) *CounterDef2[V1, V2] {
	return &CounterDef2[V1, V2]{
		name: h.name,
		prefix: []string{

			makeTag(h.keys[0], tagValueString(v0)),
		},
		keys: *((*[2]string)(h.keys[1:])),
	}
}

func (h *CounterDef3[V0, V1, V2]) BindPrefix2(v0 V0, v1 V1) *CounterDef1[V2] {
	return &CounterDef1[V2]{
		name: h.name,
		prefix: []string{

			makeTag(h.keys[0], tagValueString(v0)),

			makeTag(h.keys[1], tagValueString(v1)),
		},
		keys: *((*[1]string)(h.keys[2:])),
	}
}

type CounterDef4[V0 TagValue, V1 TagValue, V2 TagValue, V3 TagValue] struct {
	name   string
	prefix []string
	keys   [4]string
}

func NewCounterDef4[V0 TagValue, V1 TagValue, V2 TagValue, V3 TagValue](
	name string,
	description string,
	unit Unit,
	keys [4]string,
) *CounterDef4[V0, V1, V2, V3] {
	registerDef(counterType, name, unit, description)
	return &CounterDef4[V0, V1, V2, V3]{
		name: name,
		keys: keys,
	}
}

func (h *CounterDef4[V0, V1, V2, V3]) Bind(m *Metrics, v0 V0, v1 V1, v2 V2, v3 V3) *Counter {
	return m.counter(h.name, joinStrings(h.prefix, []string{

		makeTag(h.keys[0], tagValueString(v0)),

		makeTag(h.keys[1], tagValueString(v1)),

		makeTag(h.keys[2], tagValueString(v2)),

		makeTag(h.keys[3], tagValueString(v3)),
	}))
}

func (h *CounterDef4[V0, V1, V2, V3]) BindPrefix1(v0 V0) *CounterDef3[V1, V2, V3] {
	return &CounterDef3[V1, V2, V3]{
		name: h.name,
		prefix: []string{

			makeTag(h.keys[0], tagValueString(v0)),
		},
		keys: *((*[3]string)(h.keys[1:])),
	}
}

func (h *CounterDef4[V0, V1, V2, V3]) BindPrefix2(v0 V0, v1 V1) *CounterDef2[V2, V3] {
	return &CounterDef2[V2, V3]{
		name: h.name,
		prefix: []string{

			makeTag(h.keys[0], tagValueString(v0)),

			makeTag(h.keys[1], tagValueString(v1)),
		},
		keys: *((*[2]string)(h.keys[2:])),
	}
}

func (h *CounterDef4[V0, V1, V2, V3]) BindPrefix3(v0 V0, v1 V1, v2 V2) *CounterDef1[V3] {
	return &CounterDef1[V3]{
		name: h.name,
		prefix: []string{

			makeTag(h.keys[0], tagValueString(v0)),

			makeTag(h.keys[1], tagValueString(v1)),

			makeTag(h.keys[2], tagValueString(v2)),
		},
		keys: *((*[1]string)(h.keys[3:])),
	}
}

type CounterDef5[V0 TagValue, V1 TagValue, V2 TagValue, V3 TagValue, V4 TagValue] struct {
	name   string
	prefix []string
	keys   [5]string
}

func NewCounterDef5[V0 TagValue, V1 TagValue, V2 TagValue, V3 TagValue, V4 TagValue](
	name string,
	description string,
	unit Unit,
	keys [5]string,
) *CounterDef5[V0, V1, V2, V3, V4] {
	registerDef(counterType, name, unit, description)
	return &CounterDef5[V0, V1, V2, V3, V4]{
		name: name,
		keys: keys,
	}
}

func (h *CounterDef5[V0, V1, V2, V3, V4]) Bind(m *Metrics, v0 V0, v1 V1, v2 V2, v3 V3, v4 V4) *Counter {
	return m.counter(h.name, joinStrings(h.prefix, []string{

		makeTag(h.keys[0], tagValueString(v0)),

		makeTag(h.keys[1], tagValueString(v1)),

		makeTag(h.keys[2], tagValueString(v2)),

		makeTag(h.keys[3], tagValueString(v3)),

		makeTag(h.keys[4], tagValueString(v4)),
	}))
}

func (h *CounterDef5[V0, V1, V2, V3, V4]) BindPrefix1(v0 V0) *CounterDef4[V1, V2, V3, V4] {
	return &CounterDef4[V1, V2, V3, V4]{
		name: h.name,
		prefix: []string{

			makeTag(h.keys[0], tagValueString(v0)),
		},
		keys: *((*[4]string)(h.keys[1:])),
	}
}

func (h *CounterDef5[V0, V1, V2, V3, V4]) BindPrefix2(v0 V0, v1 V1) *CounterDef3[V2, V3, V4] {
	return &CounterDef3[V2, V3, V4]{
		name: h.name,
		prefix: []string{

			makeTag(h.keys[0], tagValueString(v0)),

			makeTag(h.keys[1], tagValueString(v1)),
		},
		keys: *((*[3]string)(h.keys[2:])),
	}
}

func (h *CounterDef5[V0, V1, V2, V3, V4]) BindPrefix3(v0 V0, v1 V1, v2 V2) *CounterDef2[V3, V4] {
	return &CounterDef2[V3, V4]{
		name: h.name,
		prefix: []string{

			makeTag(h.keys[0], tagValueString(v0)),

			makeTag(h.keys[1], tagValueString(v1)),

			makeTag(h.keys[2], tagValueString(v2)),
		},
		keys: *((*[2]string)(h.keys[3:])),
	}
}

func (h *CounterDef5[V0, V1, V2, V3, V4]) BindPrefix4(v0 V0, v1 V1, v2 V2, v3 V3) *CounterDef1[V4] {
	return &CounterDef1[V4]{
		name: h.name,
		prefix: []string{

			makeTag(h.keys[0], tagValueString(v0)),

			makeTag(h.keys[1], tagValueString(v1)),

			makeTag(h.keys[2], tagValueString(v2)),

			makeTag(h.keys[3], tagValueString(v3)),
		},
		keys: *((*[1]string)(h.keys[4:])),
	}
}

type GaugeDef2[V0 TagValue, V1 TagValue] struct {
	name   string
	prefix []string
	keys   [2]string
}

func NewGaugeDef2[V0 TagValue, V1 TagValue](
	name string,
	description string,
	unit Unit,
	keys [2]string,
) *GaugeDef2[V0, V1] {
	registerDef(gaugeType, name, unit, description)
	return &GaugeDef2[V0, V1]{
		name: name,
		keys: keys,
	}
}

func (h *GaugeDef2[V0, V1]) Bind(m *Metrics, v0 V0, v1 V1) *Gauge {
	return m.gauge(h.name, joinStrings(h.prefix, []string{

		makeTag(h.keys[0], tagValueString(v0)),

		makeTag(h.keys[1], tagValueString(v1)),
	}))
}

func (h *GaugeDef2[V0, V1]) BindPrefix1(v0 V0) *GaugeDef1[V1] {
	return &GaugeDef1[V1]{
		name: h.name,
		prefix: []string{

			makeTag(h.keys[0], tagValueString(v0)),
		},
		keys: *((*[1]string)(h.keys[1:])),
	}
}

type GaugeDef3[V0 TagValue, V1 TagValue, V2 TagValue] struct {
	name   string
	prefix []string
	keys   [3]string
}

func NewGaugeDef3[V0 TagValue, V1 TagValue, V2 TagValue](
	name string,
	description string,
	unit Unit,
	keys [3]string,
) *GaugeDef3[V0, V1, V2] {
	registerDef(gaugeType, name, unit, description)
	return &GaugeDef3[V0, V1, V2]{
		name: name,
		keys: keys,
	}
}

func (h *GaugeDef3[V0, V1, V2]) Bind(m *Metrics, v0 V0, v1 V1, v2 V2) *Gauge {
	return m.gauge(h.name, joinStrings(h.prefix, []string{

		makeTag(h.keys[0], tagValueString(v0)),

		makeTag(h.keys[1], tagValueString(v1)),

		makeTag(h.keys[2], tagValueString(v2)),
	}))
}

func (h *GaugeDef3[V0, V1, V2]) BindPrefix1(v0 V0) *GaugeDef2[V1, V2] {
	return &GaugeDef2[V1, V2]{
		name: h.name,
		prefix: []string{

			makeTag(h.keys[0], tagValueString(v0)),
		},
		keys: *((*[2]string)(h.keys[1:])),
	}
}

func (h *GaugeDef3[V0, V1, V2]) BindPrefix2(v0 V0, v1 V1) *GaugeDef1[V2] {
	return &GaugeDef1[V2]{
		name: h.name,
		prefix: []string{

			makeTag(h.keys[0], tagValueString(v0)),

			makeTag(h.keys[1], tagValueString(v1)),
		},
		keys: *((*[1]string)(h.keys[2:])),
	}
}

type GaugeDef4[V0 TagValue, V1 TagValue, V2 TagValue, V3 TagValue] struct {
	name   string
	prefix []string
	keys   [4]string
}

func NewGaugeDef4[V0 TagValue, V1 TagValue, V2 TagValue, V3 TagValue](
	name string,
	description string,
	unit Unit,
	keys [4]string,
) *GaugeDef4[V0, V1, V2, V3] {
	registerDef(gaugeType, name, unit, description)
	return &GaugeDef4[V0, V1, V2, V3]{
		name: name,
		keys: keys,
	}
}

func (h *GaugeDef4[V0, V1, V2, V3]) Bind(m *Metrics, v0 V0, v1 V1, v2 V2, v3 V3) *Gauge {
	return m.gauge(h.name, joinStrings(h.prefix, []string{

		makeTag(h.keys[0], tagValueString(v0)),

		makeTag(h.keys[1], tagValueString(v1)),

		makeTag(h.keys[2], tagValueString(v2)),

		makeTag(h.keys[3], tagValueString(v3)),
	}))
}

func (h *GaugeDef4[V0, V1, V2, V3]) BindPrefix1(v0 V0) *GaugeDef3[V1, V2, V3] {
	return &GaugeDef3[V1, V2, V3]{
		name: h.name,
		prefix: []string{

			makeTag(h.keys[0], tagValueString(v0)),
		},
		keys: *((*[3]string)(h.keys[1:])),
	}
}

func (h *GaugeDef4[V0, V1, V2, V3]) BindPrefix2(v0 V0, v1 V1) *GaugeDef2[V2, V3] {
	return &GaugeDef2[V2, V3]{
		name: h.name,
		prefix: []string{

			makeTag(h.keys[0], tagValueString(v0)),

			makeTag(h.keys[1], tagValueString(v1)),
		},
		keys: *((*[2]string)(h.keys[2:])),
	}
}

func (h *GaugeDef4[V0, V1, V2, V3]) BindPrefix3(v0 V0, v1 V1, v2 V2) *GaugeDef1[V3] {
	return &GaugeDef1[V3]{
		name: h.name,
		prefix: []string{

			makeTag(h.keys[0], tagValueString(v0)),

			makeTag(h.keys[1], tagValueString(v1)),

			makeTag(h.keys[2], tagValueString(v2)),
		},
		keys: *((*[1]string)(h.keys[3:])),
	}
}

type GaugeDef5[V0 TagValue, V1 TagValue, V2 TagValue, V3 TagValue, V4 TagValue] struct {
	name   string
	prefix []string
	keys   [5]string
}

func NewGaugeDef5[V0 TagValue, V1 TagValue, V2 TagValue, V3 TagValue, V4 TagValue](
	name string,
	description string,
	unit Unit,
	keys [5]string,
) *GaugeDef5[V0, V1, V2, V3, V4] {
	registerDef(gaugeType, name, unit, description)
	return &GaugeDef5[V0, V1, V2, V3, V4]{
		name: name,
		keys: keys,
	}
}

func (h *GaugeDef5[V0, V1, V2, V3, V4]) Bind(m *Metrics, v0 V0, v1 V1, v2 V2, v3 V3, v4 V4) *Gauge {
	return m.gauge(h.name, joinStrings(h.prefix, []string{

		makeTag(h.keys[0], tagValueString(v0)),

		makeTag(h.keys[1], tagValueString(v1)),

		makeTag(h.keys[2], tagValueString(v2)),

		makeTag(h.keys[3], tagValueString(v3)),

		makeTag(h.keys[4], tagValueString(v4)),
	}))
}

func (h *GaugeDef5[V0, V1, V2, V3, V4]) BindPrefix1(v0 V0) *GaugeDef4[V1, V2, V3, V4] {
	return &GaugeDef4[V1, V2, V3, V4]{
		name: h.name,
		prefix: []string{

			makeTag(h.keys[0], tagValueString(v0)),
		},
		keys: *((*[4]string)(h.keys[1:])),
	}
}

func (h *GaugeDef5[V0, V1, V2, V3, V4]) BindPrefix2(v0 V0, v1 V1) *GaugeDef3[V2, V3, V4] {
	return &GaugeDef3[V2, V3, V4]{
		name: h.name,
		prefix: []string{

			makeTag(h.keys[0], tagValueString(v0)),

			makeTag(h.keys[1], tagValueString(v1)),
		},
		keys: *((*[3]string)(h.keys[2:])),
	}
}

func (h *GaugeDef5[V0, V1, V2, V3, V4]) BindPrefix3(v0 V0, v1 V1, v2 V2) *GaugeDef2[V3, V4] {
	return &GaugeDef2[V3, V4]{
		name: h.name,
		prefix: []string{

			makeTag(h.keys[0], tagValueString(v0)),

			makeTag(h.keys[1], tagValueString(v1)),

			makeTag(h.keys[2], tagValueString(v2)),
		},
		keys: *((*[2]string)(h.keys[3:])),
	}
}

func (h *GaugeDef5[V0, V1, V2, V3, V4]) BindPrefix4(v0 V0, v1 V1, v2 V2, v3 V3) *GaugeDef1[V4] {
	return &GaugeDef1[V4]{
		name: h.name,
		prefix: []string{

			makeTag(h.keys[0], tagValueString(v0)),

			makeTag(h.keys[1], tagValueString(v1)),

			makeTag(h.keys[2], tagValueString(v2)),

			makeTag(h.keys[3], tagValueString(v3)),
		},
		keys: *((*[1]string)(h.keys[4:])),
	}
}

type HistogramDef2[V0 TagValue, V1 TagValue] struct {
	name       string
	keys       [2]string
	sampleRate float64
}

func NewHistogramDef2[V0 TagValue, V1 TagValue](
	name string,
	description string,
	unit Unit,
	keys [2]string,
	sampleRate float64,
) *HistogramDef2[V0, V1] {
	registerDef(histogramType, name, unit, description)
	return &HistogramDef2[V0, V1]{
		name:       name,
		keys:       keys,
		sampleRate: sampleRate,
	}
}

func (h *HistogramDef2[V0, V1]) Bind(m *Metrics, v0 V0, v1 V1) *Histogram {
	return &Histogram{
		m:    m,
		name: h.name,
		tags: []string{

			makeTag(h.keys[0], tagValueString(v0)),

			makeTag(h.keys[1], tagValueString(v1)),
		},
		sampleRate: h.sampleRate,
	}
}

type HistogramDef3[V0 TagValue, V1 TagValue, V2 TagValue] struct {
	name       string
	keys       [3]string
	sampleRate float64
}

func NewHistogramDef3[V0 TagValue, V1 TagValue, V2 TagValue](
	name string,
	description string,
	unit Unit,
	keys [3]string,
	sampleRate float64,
) *HistogramDef3[V0, V1, V2] {
	registerDef(histogramType, name, unit, description)
	return &HistogramDef3[V0, V1, V2]{
		name:       name,
		keys:       keys,
		sampleRate: sampleRate,
	}
}

func (h *HistogramDef3[V0, V1, V2]) Bind(m *Metrics, v0 V0, v1 V1, v2 V2) *Histogram {
	return &Histogram{
		m:    m,
		name: h.name,
		tags: []string{

			makeTag(h.keys[0], tagValueString(v0)),

			makeTag(h.keys[1], tagValueString(v1)),

			makeTag(h.keys[2], tagValueString(v2)),
		},
		sampleRate: h.sampleRate,
	}
}

type HistogramDef4[V0 TagValue, V1 TagValue, V2 TagValue, V3 TagValue] struct {
	name       string
	keys       [4]string
	sampleRate float64
}

func NewHistogramDef4[V0 TagValue, V1 TagValue, V2 TagValue, V3 TagValue](
	name string,
	description string,
	unit Unit,
	keys [4]string,
	sampleRate float64,
) *HistogramDef4[V0, V1, V2, V3] {
	registerDef(histogramType, name, unit, description)
	return &HistogramDef4[V0, V1, V2, V3]{
		name:       name,
		keys:       keys,
		sampleRate: sampleRate,
	}
}

func (h *HistogramDef4[V0, V1, V2, V3]) Bind(m *Metrics, v0 V0, v1 V1, v2 V2, v3 V3) *Histogram {
	return &Histogram{
		m:    m,
		name: h.name,
		tags: []string{

			makeTag(h.keys[0], tagValueString(v0)),

			makeTag(h.keys[1], tagValueString(v1)),

			makeTag(h.keys[2], tagValueString(v2)),

			makeTag(h.keys[3], tagValueString(v3)),
		},
		sampleRate: h.sampleRate,
	}
}

type HistogramDef5[V0 TagValue, V1 TagValue, V2 TagValue, V3 TagValue, V4 TagValue] struct {
	name       string
	keys       [5]string
	sampleRate float64
}

func NewHistogramDef5[V0 TagValue, V1 TagValue, V2 TagValue, V3 TagValue, V4 TagValue](
	name string,
	description string,
	unit Unit,
	keys [5]string,
	sampleRate float64,
) *HistogramDef5[V0, V1, V2, V3, V4] {
	registerDef(histogramType, name, unit, description)
	return &HistogramDef5[V0, V1, V2, V3, V4]{
		name:       name,
		keys:       keys,
		sampleRate: sampleRate,
	}
}

func (h *HistogramDef5[V0, V1, V2, V3, V4]) Bind(m *Metrics, v0 V0, v1 V1, v2 V2, v3 V3, v4 V4) *Histogram {
	return &Histogram{
		m:    m,
		name: h.name,
		tags: []string{

			makeTag(h.keys[0], tagValueString(v0)),

			makeTag(h.keys[1], tagValueString(v1)),

			makeTag(h.keys[2], tagValueString(v2)),

			makeTag(h.keys[3], tagValueString(v3)),

			makeTag(h.keys[4], tagValueString(v4)),
		},
		sampleRate: h.sampleRate,
	}
}

type DistributionDef2[V0 TagValue, V1 TagValue] struct {
	name       string
	keys       [2]string
	sampleRate float64
}

func NewDistributionDef2[V0 TagValue, V1 TagValue](
	name string,
	description string,
	unit Unit,
	keys [2]string,
	sampleRate float64,
) *DistributionDef2[V0, V1] {
	registerDef(distributionType, name, unit, description)
	return &DistributionDef2[V0, V1]{
		name:       name,
		keys:       keys,
		sampleRate: sampleRate,
	}
}

func (h *DistributionDef2[V0, V1]) Bind(m *Metrics, v0 V0, v1 V1) *Distribution {
	return &Distribution{
		m:    m,
		name: h.name,
		tags: []string{

			makeTag(h.keys[0], tagValueString(v0)),

			makeTag(h.keys[1], tagValueString(v1)),
		},
		sampleRate: h.sampleRate,
	}
}

type DistributionDef3[V0 TagValue, V1 TagValue, V2 TagValue] struct {
	name       string
	keys       [3]string
	sampleRate float64
}

func NewDistributionDef3[V0 TagValue, V1 TagValue, V2 TagValue](
	name string,
	description string,
	unit Unit,
	keys [3]string,
	sampleRate float64,
) *DistributionDef3[V0, V1, V2] {
	registerDef(distributionType, name, unit, description)
	return &DistributionDef3[V0, V1, V2]{
		name:       name,
		keys:       keys,
		sampleRate: sampleRate,
	}
}

func (h *DistributionDef3[V0, V1, V2]) Bind(m *Metrics, v0 V0, v1 V1, v2 V2) *Distribution {
	return &Distribution{
		m:    m,
		name: h.name,
		tags: []string{

			makeTag(h.keys[0], tagValueString(v0)),

			makeTag(h.keys[1], tagValueString(v1)),

			makeTag(h.keys[2], tagValueString(v2)),
		},
		sampleRate: h.sampleRate,
	}
}

type DistributionDef4[V0 TagValue, V1 TagValue, V2 TagValue, V3 TagValue] struct {
	name       string
	keys       [4]string
	sampleRate float64
}

func NewDistributionDef4[V0 TagValue, V1 TagValue, V2 TagValue, V3 TagValue](
	name string,
	description string,
	unit Unit,
	keys [4]string,
	sampleRate float64,
) *DistributionDef4[V0, V1, V2, V3] {
	registerDef(distributionType, name, unit, description)
	return &DistributionDef4[V0, V1, V2, V3]{
		name:       name,
		keys:       keys,
		sampleRate: sampleRate,
	}
}

func (h *DistributionDef4[V0, V1, V2, V3]) Bind(m *Metrics, v0 V0, v1 V1, v2 V2, v3 V3) *Distribution {
	return &Distribution{
		m:    m,
		name: h.name,
		tags: []string{

			makeTag(h.keys[0], tagValueString(v0)),

			makeTag(h.keys[1], tagValueString(v1)),

			makeTag(h.keys[2], tagValueString(v2)),

			makeTag(h.keys[3], tagValueString(v3)),
		},
		sampleRate: h.sampleRate,
	}
}

type DistributionDef5[V0 TagValue, V1 TagValue, V2 TagValue, V3 TagValue, V4 TagValue] struct {
	name       string
	keys       [5]string
	sampleRate float64
}

func NewDistributionDef5[V0 TagValue, V1 TagValue, V2 TagValue, V3 TagValue, V4 TagValue](
	name string,
	description string,
	unit Unit,
	keys [5]string,
	sampleRate float64,
) *DistributionDef5[V0, V1, V2, V3, V4] {
	registerDef(distributionType, name, unit, description)
	return &DistributionDef5[V0, V1, V2, V3, V4]{
		name:       name,
		keys:       keys,
		sampleRate: sampleRate,
	}
}

func (h *DistributionDef5[V0, V1, V2, V3, V4]) Bind(m *Metrics, v0 V0, v1 V1, v2 V2, v3 V3, v4 V4) *Distribution {
	return &Distribution{
		m:    m,
		name: h.name,
		tags: []string{

			makeTag(h.keys[0], tagValueString(v0)),

			makeTag(h.keys[1], tagValueString(v1)),

			makeTag(h.keys[2], tagValueString(v2)),

			makeTag(h.keys[3], tagValueString(v3)),

			makeTag(h.keys[4], tagValueString(v4)),
		},
		sampleRate: h.sampleRate,
	}
}

type SetDef2[K any, V0 TagValue, V1 TagValue] struct {
	name       string
	keys       [2]string
	sampleRate float64
}

func NewSetDef2[K any, V0 TagValue, V1 TagValue](
	name string,
	description string,
	unit Unit,
	keys [2]string,
	sampleRate float64,
) *SetDef2[K, V0, V1] {
	registerDef(setType, name, unit, description)
	return &SetDef2[K, V0, V1]{
		name:       name,
		keys:       keys,
		sampleRate: sampleRate,
	}
}

func (h *SetDef2[K, V0, V1]) Bind(m *Metrics, v0 V0, v1 V1) *Set[K] {
	return &Set[K]{
		m:    m,
		name: h.name,
		tags: []string{

			makeTag(h.keys[0], tagValueString(v0)),

			makeTag(h.keys[1], tagValueString(v1)),
		},
		sampleRate: h.sampleRate,
	}
}

type SetDef3[K any, V0 TagValue, V1 TagValue, V2 TagValue] struct {
	name       string
	keys       [3]string
	sampleRate float64
}

func NewSetDef3[K any, V0 TagValue, V1 TagValue, V2 TagValue](
	name string,
	description string,
	unit Unit,
	keys [3]string,
	sampleRate float64,
) *SetDef3[K, V0, V1, V2] {
	registerDef(setType, name, unit, description)
	return &SetDef3[K, V0, V1, V2]{
		name:       name,
		keys:       keys,
		sampleRate: sampleRate,
	}
}

func (h *SetDef3[K, V0, V1, V2]) Bind(m *Metrics, v0 V0, v1 V1, v2 V2) *Set[K] {
	return &Set[K]{
		m:    m,
		name: h.name,
		tags: []string{

			makeTag(h.keys[0], tagValueString(v0)),

			makeTag(h.keys[1], tagValueString(v1)),

			makeTag(h.keys[2], tagValueString(v2)),
		},
		sampleRate: h.sampleRate,
	}
}

type SetDef4[K any, V0 TagValue, V1 TagValue, V2 TagValue, V3 TagValue] struct {
	name       string
	keys       [4]string
	sampleRate float64
}

func NewSetDef4[K any, V0 TagValue, V1 TagValue, V2 TagValue, V3 TagValue](
	name string,
	description string,
	unit Unit,
	keys [4]string,
	sampleRate float64,
) *SetDef4[K, V0, V1, V2, V3] {
	registerDef(setType, name, unit, description)
	return &SetDef4[K, V0, V1, V2, V3]{
		name:       name,
		keys:       keys,
		sampleRate: sampleRate,
	}
}

func (h *SetDef4[K, V0, V1, V2, V3]) Bind(m *Metrics, v0 V0, v1 V1, v2 V2, v3 V3) *Set[K] {
	return &Set[K]{
		m:    m,
		name: h.name,
		tags: []string{

			makeTag(h.keys[0], tagValueString(v0)),

			makeTag(h.keys[1], tagValueString(v1)),

			makeTag(h.keys[2], tagValueString(v2)),

			makeTag(h.keys[3], tagValueString(v3)),
		},
		sampleRate: h.sampleRate,
	}
}

type SetDef5[K any, V0 TagValue, V1 TagValue, V2 TagValue, V3 TagValue, V4 TagValue] struct {
	name       string
	keys       [5]string
	sampleRate float64
}

func NewSetDef5[K any, V0 TagValue, V1 TagValue, V2 TagValue, V3 TagValue, V4 TagValue](
	name string,
	description string,
	unit Unit,
	keys [5]string,
	sampleRate float64,
) *SetDef5[K, V0, V1, V2, V3, V4] {
	registerDef(setType, name, unit, description)
	return &SetDef5[K, V0, V1, V2, V3, V4]{
		name:       name,
		keys:       keys,
		sampleRate: sampleRate,
	}
}

func (h *SetDef5[K, V0, V1, V2, V3, V4]) Bind(m *Metrics, v0 V0, v1 V1, v2 V2, v3 V3, v4 V4) *Set[K] {
	return &Set[K]{
		m:    m,
		name: h.name,
		tags: []string{

			makeTag(h.keys[0], tagValueString(v0)),

			makeTag(h.keys[1], tagValueString(v1)),

			makeTag(h.keys[2], tagValueString(v2)),

			makeTag(h.keys[3], tagValueString(v3)),

			makeTag(h.keys[4], tagValueString(v4)),
		},
		sampleRate: h.sampleRate,
	}
}
