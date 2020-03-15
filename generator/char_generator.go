package generator

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
)

const (
	text_min = 5
	text_max = 50
)

var re_char = regexp.MustCompile(`(?:var)?char\((\d+)\)`)

/* CharGenerator */
type CharGenerator struct {
	*BaseGenerator
	Len int
}

func (g *CharGenerator) init(variation int) error {
	g.variation = variation
	list, err := g.Generate_vals(variation)
	g.list = list
	return err
}

func (g *CharGenerator) Generate_vals(num int) ([]string, error) {
	return generate_vals(g, num)
}

func (g *CharGenerator) Generate_val(seed int64) string {
	return rand_letters_nums(g.Len, seed) // TODO: 文字種を変えられるように
}

/* CharRangeGenerator */
type CharRangeGenerator struct {
	*RangeGenerator
}

func (g *CharRangeGenerator) init(variation int) error {
	g.variation = variation
	list, err := g.Generate_vals(variation)
	g.list = list
	return err
}

func (g *CharRangeGenerator) Generate_vals(num int) ([]string, error) {
	return generate_vals(g, num)
}

func (g *CharRangeGenerator) Generate_val(seed int64) string {
	max, okmax := g.Max.(int)
	min, okmin := g.Min.(int)
	if okmax && okmin {
		length := min
		if diff := max - min; diff > 0 {
			rand.Seed(seed)
			length += rand.Intn(diff)
		}

		return rand_letters_nums(length, seed) // TODO: 文字種を変えられるように
	} else {
		panic("Max & Min must be int value")
	}
}

/* CharEnumGenerator */
type CharEnumGenerator struct {
	*EnumGenerator
}

const (
	letters      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	nums         = "0123456789"
	letters_nums = letters + nums
)

func rand_letters_nums(n int, seed int64) string {
	return rand_string(n, letters_nums, seed)
}

func rand_string(n int, charset string, seed int64) string {
	rand.Seed(seed)
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}

func NewCharGenerator(gtype *Gtype) GenerateVals {
	g := NewBaseGenerator(RETRY_RATIO)
	len := get_charlen(gtype.DBtype)
	switch gtype.Method {
	case EnumMethod:
		eg := NewEnumGenerator(g, gtype.Enum, gtype.Weights)
		return &CharEnumGenerator{EnumGenerator: eg}
	case "":
		return &CharGenerator{BaseGenerator: g, Len: len}
	default:
		panic(fmt.Sprintf("Unsupported method '%s' specified", gtype.Method))
	}
}

func NewVarCharGenerator(gtype *Gtype) GenerateVals {
	g := NewBaseGenerator(RETRY_RATIO)
	len := get_charlen(gtype.DBtype)
	switch gtype.Method {
	case EnumMethod:
		eg := NewEnumGenerator(g, gtype.Enum, gtype.Weights)
		return &CharEnumGenerator{EnumGenerator: eg}
	case "":
		rg := NewRangeGenerator(g, 1, len)
		return &CharRangeGenerator{RangeGenerator: rg}
	default:
		panic(fmt.Sprintf("Unsupported method '%s' specified", gtype.Method))
	}
}

func NewTextGenerator(gtype *Gtype) GenerateVals {
	g := NewBaseGenerator(RETRY_RATIO)
	switch gtype.Method {
	case RangeMethod:
		rg := NewRangeGenerator(g, gtype.Range.Min, gtype.Range.Max)
		return &CharRangeGenerator{RangeGenerator: rg}
	case EnumMethod:
		eg := NewEnumGenerator(g, gtype.Enum, gtype.Weights)
		return &CharEnumGenerator{EnumGenerator: eg}
	case "":
		rg := NewRangeGenerator(g, text_min, text_max)
		return &CharRangeGenerator{RangeGenerator: rg}
	default:
		panic(fmt.Sprintf("Unsupported method '%s' specified", gtype.Method))
	}
}

func get_charlen(chartype string) int {
	if chartype == "char" {
		return 1
	}

	len, err := strconv.Atoi(re_char.ReplaceAllString(chartype, "$1"))
	if err != nil {
		panic(err)
	}

	return len
}
