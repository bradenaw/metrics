package metrics

type gaugeGroup struct {
	curr map[*Gauge]float64
	prev map[*Gauge]float64
}

func newGaugeGroup() gaugeGroup {
	return gaugeGroup{
		curr: make(map[*Gauge]float64),
		prev: make(map[*Gauge]float64),
	}
}

func (gg *gaugeGroup) set(g *Gauge, value float64) {
	gg.curr[g] = value
}

func (gg *gaugeGroup) EmitAndUnset() {
	for gauge := range gg.prev {
		_, ok := gg.curr[gauge]
		if !ok {
			gauge.Unset()
		}
	}
	for gauge, value := range gg.curr {
		gauge.Set(value)
	}
	gg.prev = gg.curr
	gg.curr = make(map[*Gauge]float64, len(gg.prev))
}

type GaugeGroup1[V0 TagValue] struct {
	m     *Metrics
	d     GaugeDef1[V0]
	inner gaugeGroup
}

func NewGaugeGroup1[V0 TagValue](m *Metrics, d GaugeDef1[V0]) *GaugeGroup1[V0] {
	return &GaugeGroup1[V0]{
		m:     m,
		d:     d,
		inner: newGaugeGroup(),
	}
}

func (g *GaugeGroup1[V0]) Set(v0 V0, value float64) {
	g.inner.set(g.m.Gauge(g.d.Values(v0)), value)
}

// EmitAndUnset emits the gauges added using Set since the last EmitAndUnset, and unsets any tagsets
// that were not Set since the last EmitAndUnset.
func (g *GaugeGroup1[V0]) EmitAndUnset() { g.inner.EmitAndUnset() }

type GaugeGroup2[V0 TagValue, V1 TagValue] struct {
	d     *GaugeDef2[V0, V1]
	m     *Metrics
	inner gaugeGroup
}

func NewGaugeGroup2[V0 TagValue, V1 TagValue](m *Metrics, d *GaugeDef2[V0, V1]) *GaugeGroup2[V0, V1] {
	return &GaugeGroup2[V0, V1]{
		m:     m,
		d:     d,
		inner: newGaugeGroup(),
	}
}

func (g *GaugeGroup2[V0, V1]) Set(v0 V0, v1 V1, value float64) {
	g.inner.set(g.m.Gauge(g.d.Values(v0, v1)), value)
}

// EmitAndUnset emits the gauges added using Set since the last EmitAndUnset, and unsets any tagsets
// that were not Set since the last EmitAndUnset.
func (g *GaugeGroup2[V0, V1]) EmitAndUnset() { g.inner.EmitAndUnset() }
