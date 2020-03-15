package tablegenerator

import (
	"bufio"
	"data_generator_go/generator"
	"data_generator_go/helper"
	"data_generator_go/tableschema"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Generator = generator.Generator

type GenerateConf struct {
	Sep string
	Ext string
	Dir string
	Buf int
}

func Generate(table tableschema.Table,
	generator_map map[string]Generator, nworkers int, genconf *GenerateConf) {

	// ユニークなpkeyを生成する際に利用するmap
	// 各gtypeの値を参照する際のindexをどのような周期で変更するかを保持する
	pkeys_multi_rows := map[string]int{}
	rows := 1
	for _, col := range table.Columns {
		if col.Pkey == "1" {
			pkeys_multi_rows[col.Name] = rows
			rows *= generator_map[col.GetGtypeName()].GetVariation()
		}
	}

	fmt.Printf("Table: %s  Rows: %d\n", table.Name, rows)
	file_path := fmt.Sprintf("%s/%s.%s", genconf.Dir, table.Name, genconf.Ext)
	file, err := os.Create(file_path)
	if err != nil {
		fmt.Printf("Could not open file '%s'.\n", file_path)
	}

	w := bufio.NewWriterSize(file, genconf.Buf)
	wg := &sync.WaitGroup{}
	ch_in := make(chan int, rows)
	ch_out := make(chan []string, rows)
	for i := 0; i < nworkers; i++ {
		wg.Add(1)
		go rowGenerateWorker(wg, table.Columns, pkeys_multi_rows, generator_map, ch_in, ch_out)
	}

	go func() {
		if table.Limit > 0 {
			irows := helper.RandomChoice(0, rows, table.Limit)
			for _, irow := range irows {
				ch_in <- irow
			}
		} else {
			for irow := 0; irow < rows; irow++ {
				ch_in <- irow
			}
		}

		close(ch_in)
		wg.Wait()
		close(ch_out)
	}()

	for row := range ch_out {
		rowstr := strings.Join(row, genconf.Sep) + "\n"
		if _, err := w.WriteString(rowstr); err != nil {
			panic(err)
		}
	}

	w.Flush()
	return
}

func rowGenerateWorker(wg *sync.WaitGroup, cols []tableschema.Column,
	pkeys_multi_rows map[string]int, generator_map map[string]Generator, ch_in chan int, ch_out chan []string) {
	defer wg.Done()
	for {
		irow, ok := <-ch_in
		if !ok {
			return
		}

		ncols := len(cols)
		row := make([]string, ncols)
		for icol := 0; icol < ncols; icol++ {
			col := cols[icol]
			generator := generator_map[col.GetGtypeName()]
			if col.Pkey == "1" {
				row[icol] = generator.GeneratePkeyVal(irow / pkeys_multi_rows[col.Name])
			} else {
				row[icol] = generator.GenerateNonPkeyVal(irow)
			}
		}

		ch_out <- row
	}
}
