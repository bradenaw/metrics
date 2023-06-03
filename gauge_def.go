package metrics

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

type GaugeDef2[V0 TagValue, V1 TagValue] struct {
	name string
	keys [2]string
}

func NewGaugeDef2[V0 TagValue, V1 TagValue](
	name string,
	unit Unit,
	keys [2]string,
	description string,
) *GaugeDef2[V0, V1] {
	registerDef(gaugeType, name, unit, description)
	return &GaugeDef2[V0, V1]{
		name: name,
		keys: keys,
	}
}

func (d *GaugeDef2[V0, V1]) Bind(m *Metrics, v0 V0, v1 V1) *Gauge {
	return m.gauge(d.name, []string{
		makeTag(d.keys[0], tagValueString(v0)),
		makeTag(d.keys[1], tagValueString(v1)),
	})
}

type GaugeDef3[V0 TagValue, V1 TagValue, V2 TagValue] struct {
	name string
	keys [3]string
}

func NewGaugeDef3[V0 TagValue, V1 TagValue, V2 TagValue](
	name string,
	unit Unit,
	keys [3]string,
	description string,
) *GaugeDef3[V0, V1, V2] {
	registerDef(gaugeType, name, unit, description)
	return &GaugeDef3[V0, V1, V2]{
		name: name,
		keys: keys,
	}
}

func (d *GaugeDef3[V0, V1, V2]) Bind(m *Metrics, v0 V0, v1 V1, v2 V2) *Gauge {
	return m.gauge(d.name, []string{
		makeTag(d.keys[0], tagValueString(v0)),
		makeTag(d.keys[1], tagValueString(v1)),
	})
}
