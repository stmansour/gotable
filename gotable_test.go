package gotable

import (
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
	tbl.AddColumn("Index", 35, CELLSTRING, COLJUSTIFYLEFT)
	tbl.AddColumn("Value", 10, CELLINT, COLJUSTIFYRIGHT)
	tbl.AddColumn("Description", 80, CELLINT, COLJUSTIFYRIGHT)

	for i := 0; i < rows; i++ {
		tbl.AddRow()
		tbl.Puti(-1, 0, int64(i))
		tbl.Puti(-1, 1, int64(i*10))
		tbl.Puts(-1, 2, description)
	}

	return &tbl
}

// BenchMark of one table printing
func BenchmarkTableString(b *testing.B) {
	// run the table string output function b.N times
	for n := 0; n < b.N; n++ {
		tbl := getTable(10)
		_ = tbl.String()
	}
}

// BenchMark of one table printing for text export
func BenchmarkTableTEXT(b *testing.B) {
	funcname := "BenchmarkTableTEXT"

	// run the table string output function b.N times
	for n := 0; n < b.N; n++ {
		tbl := getTable(10)
		csvF, err := os.Create("benchmark.txt")
		if err != nil {
			b.Fatalf("Error <%s>: %s", funcname, err.Error())
		}
		err = tbl.CSVprintTable(csvF)
		if err != nil {
			b.Fatalf("Error <%s>: %s", funcname, err.Error())
		}
		csvF.Close()
		os.Remove(csvF.Name())
	}
}

// BenchMark of one table printing for csv export
func BenchmarkTableCSV(b *testing.B) {
	funcname := "BenchmarkTableCSV"

	// run the table string output function b.N times
	for n := 0; n < b.N; n++ {
		tbl := getTable(10)
		csvF, err := os.Create("benchmark.csv")
		if err != nil {
			b.Fatalf("Error <%s>: %s", funcname, err.Error())
		}
		err = tbl.CSVprintTable(csvF)
		if err != nil {
			b.Fatalf("Error <%s>: %s", funcname, err.Error())
		}
		csvF.Close()
		os.Remove(csvF.Name())
	}
}

// BenchMark of one table printing for html export
func BenchmarkTableHTML(b *testing.B) {
	funcname := "BenchmarkTableHTML"

	// run the table string output function b.N times
	for n := 0; n < b.N; n++ {
		tbl := getTable(10)
		csvF, err := os.Create("benchmark.html")
		if err != nil {
			b.Fatalf("Error <%s>: %s", funcname, err.Error())
		}
		err = tbl.CSVprintTable(csvF)
		if err != nil {
			b.Fatalf("Error <%s>: %s", funcname, err.Error())
		}
		csvF.Close()
		os.Remove(csvF.Name())
	}
}

// BenchMark of one table printing for pdf output
func BenchmarkTablePDF(b *testing.B) {
	funcname := "BenchmarkTablePDF"

	// run the table string output function b.N times
	for n := 0; n < b.N; n++ {
		tbl := getTable(10)
		csvF, err := os.Create("benchmark.pdf")
		if err != nil {
			b.Fatalf("Error <%s>: %s", funcname, err.Error())
		}
		err = tbl.CSVprintTable(csvF)
		if err != nil {
			b.Fatalf("Error <%s>: %s", funcname, err.Error())
		}
		csvF.Close()
		os.Remove(csvF.Name())
	}
}
