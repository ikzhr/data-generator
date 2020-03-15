package tableschema

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type TableSchema struct {
	Tables []Table `xml:"table"`
}

type Table struct {
	Name    string   `xml:"name,attr"`
	Limit   int      `xml:"_limit,attr"`
	Columns []Column `xml:"column"`
}

type Column struct {
	Name      string `xml:",chardata"`
	Type      string `xml:"type,attr"`
	Pkey      string `xml:"pkey,attr"`
	GtypeName string `xml:"_gtype,attr"`
}

func (c *Column) GetGtypeName() string {
	if c.GtypeName != "" {
		return c.GtypeName
	} else {
		return c.Name
	}
}

func Load(path string) *TableSchema {
	var parsed TableSchema
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err := xml.Unmarshal(buf, &parsed); err != nil {
		panic(err)
	}

	return &parsed
}

func (s *TableSchema) GetAllTableNames() []string {
	l := len(s.Tables)
	names := make([]string, l)
	for i := 0; i < l; i++ {
		names[i] = s.Tables[i].Name
	}

	return names
}

func (s *TableSchema) FindTable(name string) Table {
	for _, table := range s.Tables {
		if table.Name == name {
			return table
		}
	}

	panic(fmt.Sprintf("Could not find table '%s' in schema", name))
}
