package metrics

type GaugeGroup1[V0 TagValue] struct {
	m    *Metrics
	d    *GaugeDef1[V0]
	curr map[*Gauge]float64
	prev map[*Gauge]float64
}

func NewGaugeGroup1[V0 TagValue](m *Metrics, d *GaugeDef1[V0]) *GaugeGroup1[V0] {
	return &GaugeGroup1[V0]{
		m:    m,
		d:    d,
		curr: make(map[*Gauge]float64),
		prev: make(map[*Gauge]float64),
	}
}

func (g *GaugeGroup1[V0]) Set(v0 V0, value float64) {
	gauge := g.d.Values(v0).Bind(g.m)
	g.curr[gauge] = value
}

// EmitAndUnset emits the gauges added using Set since the last EmitAndUnset, and unsets any tagsets
// that were not Set since the last EmitAndUnset.
func (g *GaugeGroup1[V0]) EmitAndUnset() {
	for gauge := range g.prev {
		_, ok := g.curr[gauge]
		if !ok {
			gauge.Unset()
		}
	}
	for gauge, value := range g.curr {
		gauge.Set(value)
	}
	g.prev = g.curr
	g.curr = make(map[*Gauge]float64, len(g.prev))
}

type GaugeGroup2[V0 TagValue, V1 TagValue] struct {
	d    *GaugeDef2[V0, V1]
	m    *Metrics
	curr map[*Gauge]float64
	prev map[*Gauge]float64
}

func NewGaugeGroup2[V0 TagValue, V1 TagValue](m *Metrics, d *GaugeDef2[V0, V1]) *GaugeGroup2[V0, V1] {
	return &GaugeGroup2[V0, V1]{
		m:    m,
		d:    d,
		curr: make(map[*Gauge]float64),
		prev: make(map[*Gauge]float64),
	}
}

func (g *GaugeGroup2[V0, V1]) Set(v0 V0, v1 V1, value float64) {
	gauge := g.d.Values(v0, v1).Bind(g.m)
	g.curr[gauge] = value
}

// EmitAndUnset emits the gauges added using Set since the last EmitAndUnset, and unsets any tagsets
// that were not Set since the last EmitAndUnset.
func (g *GaugeGroup2[V0, V1]) EmitAndUnset() {
	for gauge := range g.prev {
		_, ok := g.curr[gauge]
		if !ok {
			gauge.Unset()
		}
	}
	for gauge, value := range g.curr {
		gauge.Set(value)
	}
	g.prev = g.curr
	g.curr = make(map[*Gauge]float64, len(g.prev))
}
