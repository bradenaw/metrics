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
	name string
	key  string
}

func NewCounterDef1[V0 TagValue, V1 TagValue](
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
	return m.counter(d.name, []string{
		makeTag(d.key, tagValueString(v0)),
	})
}

type CounterDef2[V0 TagValue, V1 TagValue] struct {
	name string
	keys [2]string
}

func NewCounterDef2[V0 TagValue, V1 TagValue](
	name string,
	unit Unit,
	keys [2]string,
	description string,
) *CounterDef2[V0, V1] {
	registerDef(counterType, name, unit, description)
	return &CounterDef2[V0, V1]{
		name: name,
		keys: keys,
	}
}

// Bind binds a set of tag values to this definition and returns the resulting Counter.
func (d *CounterDef2[V0, V1]) Bind(m *Metrics, v0 V0, v1 V1) *Counter {
	return m.counter(d.name, []string{
		makeTag(d.keys[0], tagValueString(v0)),
		makeTag(d.keys[1], tagValueString(v1)),
	})
}

type CounterDef3[V0 TagValue, V1 TagValue, V2 TagValue] struct {
	name string
	keys [3]string
}

func NewCounterDef3[V0 TagValue, V1 TagValue, V2 TagValue](
	name string,
	unit Unit,
	keys [3]string,
	description string,
) *CounterDef3[V0, V1, V2] {
	registerDef(counterType, name, unit, description)
	return &CounterDef3[V0, V1, V2]{
		name: name,
		keys: keys,
	}
}

// Bind binds a set of tag values to this definition and returns the resulting Counter.
func (d *CounterDef3[V0, V1, V2]) Bind(m *Metrics, v0 V0, v1 V1, v2 V2) *Counter {
	return m.counter(d.name, []string{
		makeTag(d.keys[0], tagValueString(v0)),
		makeTag(d.keys[1], tagValueString(v1)),
		makeTag(d.keys[2], tagValueString(v2)),
	})
}
