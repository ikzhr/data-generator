package generator

import (
	"data_generator_go/gtype"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
)

const RETRY_RATIO = 0.5

type Gtype = gtype.Gtype

var RangeMethod = gtype.RangeMethod
var EnumMethod = gtype.EnumMethod

type Generator interface {
	init(int) error
	GeneratePkeyVal(int) string
	GenerateNonPkeyVal(int) string
	GetVariation() int
}

type GenerateVals interface {
	Generator
	Generate_vals(int) ([]string, error)
}

type GenerateVal interface {
	Generate_val(int64) string
	getRetryRatio() float32
}

type BaseGenerator struct {
	variation  int
	list       []string
	retryRatio float32
}

func NewBaseGenerator(retryRatio float32) *BaseGenerator {
	return &BaseGenerator{retryRatio: retryRatio}
}

func (g *BaseGenerator) getRetryRatio() float32 {
	return g.retryRatio
}

func (g *BaseGenerator) GetVariation() int {
	return g.variation
}

func (g *BaseGenerator) GeneratePkeyVal(idx int) string {
	return g.list[idx%g.variation]
}

func (g *BaseGenerator) GenerateNonPkeyVal(idx int) string {
	return g.list[idx%g.variation]
}

type RangeGenerator struct {
	*BaseGenerator
	Min interface{}
	Max interface{}
}

func NewRangeGenerator(g *BaseGenerator, min interface{}, max interface{}) *RangeGenerator {
	return &RangeGenerator{BaseGenerator: g, Min: min, Max: max}
}

type EnumGenerator struct {
	*BaseGenerator
	Enum    []string
	Weights []float32
}

func NewEnumGenerator(g *BaseGenerator, enum []string, weights []float32) *EnumGenerator {
	// TODO: wを作る処理を関数に外だし
	l := len(enum)
	w := make([]float32, l)
	wl := len(weights)
	if wl == 0 {
		fl := float32(l)
		for idx, _ := range enum {
			w[idx] = 1 / fl
		}
	} else {
		if wl != l {
			log.Fatal(fmt.Sprintf("Weights length(%d) must match with Enum length(%d)", wl, l))
		}

		var sum float32
		for _, v := range weights {
			sum += v
		}

		for idx, _ := range enum {
			w[idx] = weights[idx] / sum
		}
	}

	return &EnumGenerator{BaseGenerator: g, Enum: enum, Weights: w}
}

func (g *EnumGenerator) init(variation int) error {
	g.variation = len(g.Enum) // enum型では指定された値のvariationをそのまま使う
	list, err := g.Generate_vals(variation)
	g.list = list
	return err
}

func (g *EnumGenerator) GenerateNonPkeyVal(idx int) string {
	return g.list[g.weightIdx(idx)]
}

func (g *EnumGenerator) weightIdx(idx int) int {
	rand.Seed(int64(idx))
	v := rand.Float32()
	weightedIdx := 0
	for i, weight := range g.Weights {
		if weight < v {
			v -= weight
		} else {
			weightedIdx = i
			break
		}
	}

	return weightedIdx
}

func (g *EnumGenerator) Generate_vals(num int) ([]string, error) {
	return g.Enum, nil
}

func generate_vals(g GenerateVal, num int) ([]string, error) {
	l := make([]string, num)
	memo := map[string]struct{}{}
	max_loop := int(float32(num) * (1.0 + g.getRetryRatio()))
	count := 0
	for i := 0; i < max_loop; i++ {
		v := g.Generate_val(int64(num + i)) // TODO: seedが妥当か検討
		if _, ok := memo[v]; !ok {
			memo[v] = struct{}{}
			l[count] = v
			count += 1
			if count == num {
				break
			}
		}
	}

	if count < num {
		return l, errors.New("Enough values are not generated")
	}

	return l, nil
}

func NewGenerator(gtype *Gtype, defaultVariation int) Generator {
	var generator Generator
	dbtype := gtype.DBtype

	switch {
	case dbtype == "int2":
		generator = NewInt16Generator(gtype)
	case dbtype == "int":
		generator = NewInt32Generator(gtype)
	case dbtype == "int4":
		generator = NewInt32Generator(gtype)
	case dbtype == "int8":
		generator = NewInt64Generator(gtype)
	case strings.Contains(dbtype, "varchar"): // TODO: より安全なマッチングにする
		generator = NewVarCharGenerator(gtype)
	case strings.Contains(dbtype, "char"): // TODO: より安全なマッチングにする
		generator = NewCharGenerator(gtype)
	case dbtype == "text":
		generator = NewTextGenerator(gtype)
	case dbtype == "date":
		generator = NewDateGenerator(gtype)
	case dbtype == "timestamp":
		generator = NewTimestampGenerator(gtype)
	case dbtype == "time":
		generator = NewTimeGenerator(gtype)
	case dbtype == "boolean":
		generator = NewBoolGenerator(gtype)
	default:
		panic(fmt.Sprintf("Unsupported type '%s' specified", dbtype))
	}

	n := gtype.Num
	if n == 0 {
		n = defaultVariation
	}

	err := generator.init(n)
	if err != nil {
		log.Fatal(gtype.Name, ": ", err)
	}

	return generator
}
