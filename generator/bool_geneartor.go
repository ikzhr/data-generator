package generator

func NewBoolGenerator(gtype *Gtype) GenerateVals {
	gtype.Method = EnumMethod
	gtype.Enum = []string{"true", "false"}
	return NewTextGenerator(gtype)
}
