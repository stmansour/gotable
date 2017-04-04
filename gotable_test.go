package gotable

import (
	// "bytes"
	"os"
	"testing"
)

// returns a table object
func getTable(rows int) *Table {

	const description = "Lorem ipsum dolor sit amet, elementum fermentum suspendisse"

	var tbl Table
	tbl.Init()
	tbl.SetTitle("Sample Table")
	tbl.SetSection1("Sample Table Section One")
	tbl.SetSection2("Sample Table Section Two")

	// add columns
	tbl.AddColumn("Index", 10, CELLINT, COLJUSTIFYLEFT)
	tbl.AddColumn("Value", 10, CELLINT, COLJUSTIFYRIGHT)
	tbl.AddColumn("Description", 80, CELLSTRING, COLJUSTIFYLEFT)

	for i := 0; i < rows; i++ {
		tbl.AddRow()
		tbl.Puti(-1, 0, int64(i))
		tbl.Puti(-1, 1, int64(i*10))
		tbl.Puts(-1, 2, description)
	}

	return &tbl
}

// // BenchMark of one table printing
// func BenchmarkTableString(b *testing.B) {
// 	// run the table string output function b.N times
// 	for n := 0; n < b.N; n++ {
// 		temp := bytes.Buffer{}
// 		tbl := getTable(10)
// 		temp.WriteString(tbl.String())
// 	}
// }

// BenchMark of one table printing for text export
func BenchmarkTableTEXT(b *testing.B) {
	funcname := "BenchMarkTableTEXT"

	benchmarks := []struct {
		name string
		rows int
	}{
		{"Table-TEXT-Rows-10", 10},
		{"Table-TEXT-Rows-100", 100},
		{"Table-TEXT-Rows-1000", 1000},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {

			for n := 0; n < b.N; n++ {
				tbl := getTable(bm.rows)
				f, err := os.Create("benchmark.txt")
				if err != nil {
					b.Fatalf("Error <%s>: %s", funcname, err.Error())
				}
				err = tbl.TextprintTable(f)
				if err != nil {
					b.Fatalf("Error <%s>: %s", funcname, err.Error())
				}
				f.Close()
				os.Remove(f.Name())
			}

		})
	}
}

// BenchMark of one table printing for csv export
func BenchmarkTableCSV(b *testing.B) {
	funcname := "BenchmarkTableCSV"

	benchmarks := []struct {
		name string
		rows int
	}{
		{"Table-CSV-Rows-10", 10},
		{"Table-CSV-Rows-100", 100},
		{"Table-CSV-Rows-1000", 1000},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {

			for n := 0; n < b.N; n++ {
				tbl := getTable(bm.rows)
				f, err := os.Create("benchmark.csv")
				if err != nil {
					b.Fatalf("Error <%s>: %s", funcname, err.Error())
				}
				err = tbl.CSVprintTable(f)
				if err != nil {
					b.Fatalf("Error <%s>: %s", funcname, err.Error())
				}
				f.Close()
				os.Remove(f.Name())
			}

		})
	}
}

// BenchMark of one table printing for html export
func BenchmarkTableHTML(b *testing.B) {
	funcname := "BenchmarkTableHTML"

	benchmarks := []struct {
		name string
		rows int
	}{
		{"Table-HTML-Rows-10", 10},
		{"Table-HTML-Rows-100", 100},
		{"Table-HTML-Rows-1000", 1000},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {

			for n := 0; n < b.N; n++ {
				tbl := getTable(bm.rows)
				f, err := os.Create("benchmark.html")
				if err != nil {
					b.Fatalf("Error <%s>: %s", funcname, err.Error())
				}
				err = tbl.HTMLprintTable(f)
				if err != nil {
					b.Fatalf("Error <%s>: %s", funcname, err.Error())
				}
				f.Close()
				os.Remove(f.Name())
			}

		})
	}
}

// BenchMark of one table printing for pdf output
func BenchmarkTablePDF(b *testing.B) {
	funcname := "BenchmarkTablePDF"

	benchmarks := []struct {
		name string
		rows int
	}{
		{"Table-PDF-Rows-10", 10},
		{"Table-PDF-Rows-100", 100},
		{"Table-PDF-Rows-1000", 1000},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {

			for n := 0; n < b.N; n++ {
				tbl := getTable(bm.rows)
				f, err := os.Create("benchmark.pdf")
				if err != nil {
					b.Fatalf("Error <%s>: %s", funcname, err.Error())
				}
				err = tbl.PDFprintTable(f)
				if err != nil {
					b.Fatalf("Error <%s>: %s", funcname, err.Error())
				}
				f.Close()
				os.Remove(f.Name())
			}

		})
	}
}
