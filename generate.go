package main

import (
	"data_generator_go/generator"
	"data_generator_go/gtype"
	"data_generator_go/tablegenerator"
	"data_generator_go/tableschema"
	"flag"
	"strings"
)

type Generator = generator.Generator
type Gtype = gtype.Gtype
type TableSchema = tableschema.TableSchema
type Column = tableschema.Column
type Table = tableschema.Table
type GenerateConf = tablegenerator.GenerateConf

const DEFAULT_GENERATE_NUM = 10
const DEFAULT_MAX_WORKERS = 1
const DEFAULT_BUF_SIZE = 1024 * 1024

func main() {
	var (
		gtype_path           = flag.String("c", "", "Path to column types def file")
		out_dir              = flag.String("o", "./", "Path to dir to output")
		file_ext             = flag.String("e", "tsv", "Extension of file to output")
		out_sep              = flag.String("s", "\t", "Separator of output data")
		table_names_raw      = flag.String("t", "", "List of table names")
		default_generate_num = flag.Int("n", DEFAULT_GENERATE_NUM, "The number of values generated")
		output_buf_size      = flag.Int("b", DEFAULT_BUF_SIZE, "Buffer size of writer")
		nworkers             = flag.Int("w", DEFAULT_MAX_WORKERS, "Woker size to generate table data")
	)

	flag.Parse()
	schema_path := flag.Arg(0)
	if schema_path == "" {
		panic("Path to schema file must be specifed")
	}

	schema := tableschema.Load(schema_path)
	genconf := &GenerateConf{
		Sep: *out_sep, Ext: *file_ext, Dir: *out_dir, Buf: *output_buf_size}
	table_names := strings.Split(*table_names_raw, ",")
	if table_names[0] == "" {
		table_names = schema.GetAllTableNames()
	}

	gtype_conf := gtype.LoadConfig(*gtype_path)
	generator_map := makeGeneratorMap(schema, gtype_conf.GenerateGtypesMap(), *default_generate_num)
	generate_tables_data(schema, table_names, generator_map, genconf, *nworkers)
}

func makeGeneratorMap(schema *TableSchema, gtypes map[string]Gtype, defaultVariation int) map[string]Generator {
	genMap := map[string]Generator{}
	for _, table := range schema.Tables {
		for _, col := range table.Columns {
			gtypeName := col.GetGtypeName()
			if _, ok := genMap[gtypeName]; ok {
				continue
			} else if val, ok := gtypes[gtypeName]; ok {
				genMap[gtypeName] = generator.NewGenerator(&val, defaultVariation)
			} else {
				genMap[gtypeName] = generator.NewGenerator(
					&Gtype{Name: gtypeName, DBtype: col.Type}, defaultVariation)
			}
		}
	}

	return genMap
}

func generate_tables_data(schema *TableSchema, table_names []string,
	generator_map map[string]Generator, genconf *GenerateConf, nworkers int) {
	for _, table_name := range table_names {
		table := schema.FindTable(table_name)
		tablegenerator.Generate(table, generator_map, nworkers, genconf)
	}
}
