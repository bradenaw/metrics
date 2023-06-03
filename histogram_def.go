package metrics

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

type HistogramDef2[V0 TagValue, V1 TagValue] struct {
	name       string
	keys       [2]string
	sampleRate float64
}

func NewHistogramDef2[V0 TagValue, V1 TagValue](
	name string,
	unit Unit,
	keys [2]string,
	description string,
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
