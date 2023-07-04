package metrics

import (
	"reflect"
)

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
	ok := registerDef(CounterType, name, description, unit, nil, nil)
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
	keys [1]string,
) CounterDef1[V0] {
	var zero0 V0
	ok := registerDef(
		CounterType,
		name,
		description,
		unit,
		keys[:],
		[]reflect.Type{reflect.TypeOf(zero0)},
	)
	return CounterDef1[V0]{
		name: name,
		keys: keys,
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
	ok := registerDef(GaugeType, name, description, unit, nil, nil)
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
	keys [1]string,
) GaugeDef1[V0] {
	var zero0 V0
	ok := registerDef(
		GaugeType,
		name,
		description,
		unit,
		keys[:],
		[]reflect.Type{reflect.TypeOf(zero0)},
	)
	return GaugeDef1[V0]{
		name: name,
		keys: keys,
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
	ok := registerDef(HistogramType, name, description, unit, nil, nil)
	return HistogramDef{
		name:       name,
		sampleRate: sampleRate,
		ok:         ok,
	}
}

type HistogramDef1[V0 TagValue] struct {
	name       string
	keys       [1]string
	sampleRate float64
	ok         bool
}

func NewHistogramDef1[V0 TagValue](
	name string,
	description string,
	unit Unit,
	keys [1]string,
	sampleRate float64,
) HistogramDef1[V0] {
	var zero0 V0
	ok := registerDef(
		HistogramType,
		name,
		description,
		unit,
		keys[:],
		[]reflect.Type{reflect.TypeOf(zero0)},
	)
	return HistogramDef1[V0]{
		name:       name,
		keys:       keys,
		sampleRate: sampleRate,
		ok:         ok,
	}
}

func (d HistogramDef1[V0]) Values(v0 V0) HistogramDef {
	return HistogramDef{
		name: d.name,
		tags: []string{
			makeTag(d.keys[0], tagValueString(v0)),
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
	ok := registerDef(DistributionType, name, description, unit, nil, nil)
	return DistributionDef{
		name:       name,
		sampleRate: sampleRate,
		ok:         ok,
	}
}

type DistributionDef1[V0 TagValue] struct {
	name       string
	keys       [1]string
	sampleRate float64
	ok         bool
}

func NewDistributionDef1[V0 TagValue](
	name string,
	description string,
	unit Unit,
	keys [1]string,
	sampleRate float64,
) DistributionDef1[V0] {
	var zero0 V0
	ok := registerDef(
		DistributionType,
		name,
		description,
		unit,
		keys[:],
		[]reflect.Type{reflect.TypeOf(zero0)},
	)
	return DistributionDef1[V0]{
		name:       name,
		keys:       keys,
		sampleRate: sampleRate,
		ok:         ok,
	}
}

func (d DistributionDef1[V0]) Values(v0 V0) DistributionDef {
	return DistributionDef{
		name: d.name,
		tags: []string{
			makeTag(d.keys[0], tagValueString(v0)),
		},
		sampleRate: d.sampleRate,
		ok:         d.ok,
	}
}

type SetDef struct {
	name       string
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
	ok := registerDef(SetType, name, description, unit, nil, nil)
	return SetDef{
		name:       name,
		sampleRate: sampleRate,
		ok:         ok,
	}
}

type SetDef1[V0 TagValue] struct {
	name       string
	keys       [1]string
	sampleRate float64
	ok         bool
}

func NewSetDef1[V0 TagValue](
	name string,
	description string,
	unit Unit,
	keys [1]string,
	sampleRate float64,
) SetDef1[V0] {
	var zero0 V0
	ok := registerDef(
		SetType,
		name,
		description,
		unit,
		keys[:],
		[]reflect.Type{reflect.TypeOf(zero0)},
	)
	return SetDef1[V0]{
		name:       name,
		keys:       keys,
		sampleRate: sampleRate,
		ok:         ok,
	}
}

func (d SetDef1[V0]) Values(v0 V0) SetDef {
	return SetDef{
		name: d.name,
		tags: []string{
			makeTag(d.keys[0], tagValueString(v0)),
		},
		sampleRate: d.sampleRate,
		ok:         d.ok,
	}
}
