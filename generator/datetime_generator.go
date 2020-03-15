package generator

import (
	"fmt"
	"math/rand"
	"time"
)

var max_time = time.Date(9999, time.December, 31, 12, 59, 59, 9999999, time.Local)
var min_time = time.Time{}
var duration_time = max_time.Sub(min_time)

/* DateTimeGenerator */
type DateTimeGenerator struct {
	*BaseGenerator
	Format string
}

func (g *DateTimeGenerator) init(variation int) error {
	g.variation = variation
	list, err := g.Generate_vals(variation)
	g.list = list
	return err
}

func (g *DateTimeGenerator) Generate_vals(num int) ([]string, error) {
	return generate_vals(g, num)
}

func (g *DateTimeGenerator) Generate_val(seed int64) string {
	s := rand_time(min_time, duration_time, seed).Format(g.Format)
	return s
}

/* DateTimeRangeGenerator */
type DateTimeRangeGenerator struct {
	*RangeGenerator
	Format string
}

func (g *DateTimeRangeGenerator) init(variation int) error {
	g.variation = variation
	list, err := g.Generate_vals(variation)
	g.list = list
	return err
}

func (g *DateTimeRangeGenerator) Generate_vals(num int) ([]string, error) {
	return generate_vals(g, num)
}

func (g *DateTimeRangeGenerator) Generate_val(seed int64) string {
	max, okmax := g.Max.(string)
	min, okmin := g.Min.(string)
	if !okmax || !okmin {
		panic("Max & Min must be date format string")
	}

	tmax, err := time.Parse(g.Format, max)
	if err != nil {
		panic(fmt.Sprintf("Max does not match with time format string '%s'", g.Format))
	}

	tmin, err := time.Parse(g.Format, min)
	if err != nil {
		panic(fmt.Sprintf("Min does not match with time format string '%s'", g.Format))
	}

	return rand_time(tmin, tmax.Sub(tmin), seed).Format(g.Format)
}

func rand_time(min time.Time, duration time.Duration, seed int64) time.Time {
	rand.Seed(seed)
	return min.Add(time.Duration(float64(duration) * rand.Float64()))
}

/* DateTimeEnumGenerator */
type DateTimeEnumGenerator struct {
	*EnumGenerator
}

func NewDateGenerator(gtype *Gtype) GenerateVals {
	format := gtype.TimeFmt
	if format == "" {
		// default date format
		format = `2006-01-02`
	}

	return newDateTimeGenerator(gtype, format)
}

func NewTimestampGenerator(gtype *Gtype) GenerateVals {
	format := gtype.TimeFmt
	if format == "" {
		// default timestamp format
		format = `2006-01-02 15:04:05`
	}
	return newDateTimeGenerator(gtype, format)
}

func NewTimeGenerator(gtype *Gtype) GenerateVals {
	format := gtype.TimeFmt
	if format == "" {
		// default timestamp format
		format = `15:04`
	}
	return newDateTimeGenerator(gtype, format)
}

func newDateTimeGenerator(gtype *Gtype, format string) GenerateVals {
	g := NewBaseGenerator(RETRY_RATIO)
	switch gtype.Method {
	case RangeMethod:
		rg := NewRangeGenerator(g, gtype.Range.Min, gtype.Range.Max)
		return &DateTimeRangeGenerator{RangeGenerator: rg, Format: format}
	case EnumMethod:
		eg := NewEnumGenerator(g, gtype.Enum, gtype.Weights)
		return &DateTimeEnumGenerator{EnumGenerator: eg}
	default:
		return &DateTimeGenerator{BaseGenerator: g, Format: format}
	}
}
