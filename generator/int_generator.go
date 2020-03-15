package generator

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
)

const (
	min_int16 int64 = -32768
	max_int16 int64 = 32767

	min_int32 int64 = -2147483648
	max_int32 int64 = 2147483647

	min_int64 int64 = -9223372036854775808
	max_int64 int64 = 9223372036854775807
)

/* IntGenerator */
type IntGenerator struct {
	*BaseGenerator
}

func (g *IntGenerator) init(variation int) error {
	g.variation = variation
	list, err := g.Generate_vals(variation)
	g.list = list
	return err
}

func (g *IntGenerator) Generate_vals(num int) ([]string, error) {
	return generate_vals(g, num)
}

func (g *IntGenerator) Generate_val(seed int64) string {
	rand.Seed(seed)
	return strconv.FormatInt(int64(rand.Int()), 10)
}

/* IntRangeGenerator */
type IntRangeGenerator struct {
	*RangeGenerator
}

func (g *IntRangeGenerator) init(variation int) error {
	g.variation = variation
	list, err := g.Generate_vals(variation)
	g.list = list
	return err
}

func (g *IntRangeGenerator) Generate_vals(num int) ([]string, error) {
	return generate_vals(g, num)
}

func (g *IntRangeGenerator) Generate_val(seed int64) string {
	var val int64
	rand.Seed(seed)
	max, okmax := g.Max.(int64)
	min, okmin := g.Min.(int64)
	if !okmax || !okmin {
		log.Fatal("Max & Min must be int64 value")
	}

	diff, ok := safeSubInt64(max, min)
	if ok {
		val = int64(float64(diff+1)*rand.Float64()) + min
	} else {
		// overflowしないように浮動小数点計算
		val = int64((float64(max)-float64(min)+1.0)*rand.Float64() + float64(min))
	}

	return strconv.FormatInt(val, 10)
}

/* IntEnumGenerator */
type IntEnumGenerator struct {
	*EnumGenerator
}

func NewInt16Generator(gtype *Gtype) GenerateVals {
	return newIntGenerator(gtype, min_int16, max_int16)
}

func NewInt32Generator(gtype *Gtype) GenerateVals {
	return newIntGenerator(gtype, min_int32, max_int32)
}

func NewInt64Generator(gtype *Gtype) GenerateVals {
	return newIntGenerator(gtype, min_int64, max_int64)
}

func newIntGenerator(gtype *Gtype, type_min int64, type_max int64) GenerateVals {
	g := NewBaseGenerator(RETRY_RATIO)
	switch gtype.Method {
	case RangeMethod:
		min := convert2int64(gtype.Range.Min)
		max := convert2int64(gtype.Range.Max)
		return &IntRangeGenerator{RangeGenerator: NewRangeGenerator(g, min, max)}
	case EnumMethod:
		return &IntEnumGenerator{EnumGenerator: NewEnumGenerator(g, gtype.Enum, gtype.Weights)}
	default:
		return &IntRangeGenerator{RangeGenerator: NewRangeGenerator(g, type_min, type_max)}
	}
}

func convert2int64(n interface{}) int64 {
	switch n.(type) {
	case int:
		return int64(n.(int))
	case int8:
		return int64(n.(int8))
	case int16:
		return int64(n.(int16))
	case int32:
		return int64(n.(int32))
	case int64:
		return int64(n.(int64))
	default:
		panic(fmt.Sprintf("Can not convert %d to int64", n))
	}
}

func safeSubInt64(a int64, b int64) (int64, bool) {
	if (b < 0 && a > b+math.MaxInt64) || (b > 0 && a < b+math.MinInt64) {
		return 0, false
	}

	return a - b, true
}
