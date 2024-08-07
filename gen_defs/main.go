package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func main() {
	n := 6

	fmt.Println("package metrics")
	fmt.Println()
	fmt.Println("import \"reflect\"")
	fmt.Println()
	fmt.Println("// generated by `go run ./gen_defs > defs_generated.go && gofmt -w defs_generated.go`")
	fmt.Println()
	fmt.Printf("const maxTags = %d\n", n-1)

	type vars struct {
		N           int
		Ns          []int
		Metric      string
		MetricLower string
		SampleRate  bool
		Unit        bool
	}

	ns := make([]int, n)
	for j := range ns {
		ns[j] = j
	}

	type metricOpts struct {
		Name       string
		Unit       bool
		SampleRate bool
	}

	for _, metric := range []metricOpts{
		{Name: "Counter", SampleRate: false, Unit: false},
		{Name: "Gauge", SampleRate: false, Unit: false},
		{Name: "Distribution", SampleRate: true, Unit: true},
		{Name: "Set", SampleRate: true, Unit: false},
	} {
		for i := 1; i < n; i++ {
			err := metricTmpl.Execute(os.Stdout, vars{
				N:           i,
				Ns:          ns[:i],
				Metric:      metric.Name,
				MetricLower: strings.ToLower(metric.Name),
				SampleRate:  metric.SampleRate,
				Unit:        metric.Unit,
			})
			if err != nil {
				panic(err)
			}

			for k := 1; k <= i-1; k++ {
				bindPrefixTmpl.Execute(os.Stdout, struct {
					N          int
					Ns         []int
					K          int
					Ks         []int
					NMinusK    int
					NMinusKs   []int
					Metric     string
					SampleRate bool
					Unit       bool
				}{
					N:          i,
					Ns:         ns[:i],
					K:          k,
					Ks:         ns[:k],
					NMinusK:    i - k,
					NMinusKs:   ns[k:i],
					Metric:     metric.Name,
					SampleRate: metric.SampleRate,
					Unit:       metric.Unit,
				})
			}
		}
	}
}

var metricTmpl = template.Must(template.New("name").Parse(`
// {{.Metric}}Def{{.N}} is the definition of a {{.MetricLower}} metric with {{.N}} tag(s).
type {{.Metric}}Def{{.N}}[{{range .Ns}} V{{.}} TagValue, {{end}}] struct {
	name       string
	{{if .Unit}} unit Unit {{end}}
	prefix     tags
	keys       [{{.N}}]string
	{{if .SampleRate}} sampleRate float64 {{end}}
	allComparable bool
	ok            bool
}

// New{{.Metric}}Def{{.N}} defines a {{.MetricLower}} metric with {{.N}} tag(s).
//
// It must be called from a top-level var block in a file called metrics.go, otherwise it will panic
// (if main() has not yet started) or return an inert def that will not produce any data.
func New{{.Metric}}Def{{.N}}[{{range .Ns}} V{{.}} TagValue, {{end}}](
	name string,
	description string,
	unit Unit,
	keys [{{.N}}]string,
	{{if .SampleRate}} sampleRate float64, {{end}}
) {{.Metric}}Def{{.N}}[{{range .Ns}} V{{.}}, {{end}}] {
	{{range .Ns}}var zero{{.}} V{{.}}
	{{ end }}
	ok := registerDef(
		{{.Metric}}Type,
		name,
		description,
		unit,
		keys[:],
		[]reflect.Type{
			{{range .Ns}}reflect.TypeOf(zero{{.}}),
			{{ end }}
		},
	)
	return {{.Metric}}Def{{.N}}[{{range .Ns}} V{{.}}, {{end}}]{
		name:       name,
		{{if .Unit}}unit: unit,{{end}}
		keys:       keys,
		{{if .SampleRate}}sampleRate: sampleRate,{{end}}
		allComparable: {{range .Ns}}reflect.TypeOf(zero{{.}}).Comparable() &&
		{{end}} true,
		ok:         ok,
	}
}

// Values returns a {{.Metric}}Def that has all of the given tag values bound. It can be passed to
// Metrics.{{.Metric}}() to produce a metric to log data to.
func (d {{.Metric}}Def{{.N}}[{{range .Ns}} V{{.}}, {{end}}]) Values({{range .Ns}} v{{.}} V{{.}}, {{end}}) {{.Metric}}Def {
	t := tags{n: {{.N}}}
	copy(t.keys[:], d.keys[:])
	{{range .Ns}}t.values[{{.}}] = v{{.}}
	{{end}}
	return {{.Metric}}Def{
		name: d.name,
		{{if .Unit}}unit: d.unit,{{end}}
		tags: d.prefix.append(t),
		{{if .SampleRate}}sampleRate: d.sampleRate,{{end}}
		allComparable: d.allComparable,
		ok: d.ok,
	}
}
`))

var bindPrefixTmpl = template.Must(template.New("name").Parse(`
// Prefix{{.K}} sets the value of the first {{.K}} tags, returning a {{.Metric}}Def{{.NMinusK}} that
// can be used to set the rest.
func (d {{.Metric}}Def{{.N}}[{{range .Ns}} V{{.}}, {{end}}]) Prefix{{.K}}({{range .Ks}} v{{.}} V{{.}}, {{end}}) {{.Metric}}Def{{.NMinusK}}[{{range .NMinusKs}} V{{.}}, {{end}}] {
	t := tags{n: {{.K}}}
	copy(t.keys[:], d.keys[:{{.K}}])
	{{range .Ks}}t.values[{{.}}] = v{{.}}
	{{end}}

	return {{.Metric}}Def{{.NMinusK}}[{{range .NMinusKs}} V{{.}}, {{end}}]{
		name: d.name,
		{{if .Unit}}unit: d.unit,{{end}}
		prefix: t,
		keys: *((*[{{.NMinusK}}]string)(d.keys[{{.K}}:])),
		{{if .SampleRate}}sampleRate: d.sampleRate,{{end}}
		allComparable: d.allComparable,
		ok:   d.ok,
	}
}
`))
