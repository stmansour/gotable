package gotable

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"
	"time"
)

var pdfProps = []*PDFProperty{
	//
	{Option: "--no-collate"},
	// top margin
	{Option: "-T", Value: "15"},
	// header center content
	{Option: "--header-center", Value: "Smoke Test Report Table"},
	// header font size
	{Option: "--header-font-size", Value: "7"},
	// header font
	{Option: "--header-font-name", Value: "opensans"},
	// header spacing
	{Option: "--header-spacing", Value: "3"},
	// bottom margin
	{Option: "-B", Value: "15"},
	// footer spacing
	{Option: "--footer-spacing", Value: "5"},
	// footer font
	{Option: "--footer-font-name", Value: "opensans"},
	// footer font size
	{Option: "--footer-font-size", Value: "7"},
	// footer left content
	{Option: "--footer-left", Value: time.Now().Format(DATETIMEFMT)},
	// footer right content
	{Option: "--footer-right", Value: "Page [page] of [toPage]"},
	// page size
	{Option: "--page-size", Value: "Letter"},
	// orientation
	{Option: "--orientation", Value: "Landscape"},
}

func TestSmoke(t *testing.T) {

	var tbl Table

	// stringln case
	if title := stringln(tbl.Title); title != "" {
		t.Errorf("smoke_test: Expected blank title, but found: %s\n", title)
	}
	titleSmpl := "GOTable\n"
	if titleOut := stringln(titleSmpl); titleSmpl != titleOut {
		t.Errorf("smoke_test: Expected %s title, but found: %s\n", titleSmpl, titleOut)
	}

	title := "GOTABLE"
	section1 := "A Smoke Test"
	section2 := "February 21, 2017"
	section3 := "section3"
	tbl.Init() //sets column spacing and date format to default

	tbl.SetNoRowsCSS([]*CSSProperty{{Name: "font-family", Value: "monospace"}})
	tbl.SetNoHeadersCSS([]*CSSProperty{{Name: "font-family", Value: "monospace"}})

	// force some edge condition errors...
	errExp := "No Header Columns"
	// headers check
	err := tbl.HasHeaders()
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}
	// output buffer for each type for table
	var temp bytes.Buffer
	// text version
	textOut, err := tbl.SprintTable()
	if !strings.Contains(textOut, errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}
	// FprintTable of text version
	err = tbl.FprintTable(&temp)
	if !strings.Contains(temp.String(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}
	// csv version
	err = tbl.CSVprintTable(&temp)
	if !strings.Contains(temp.String(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}
	// html version
	err = tbl.HTMLprintTable(&temp)
	if !strings.Contains(temp.String(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}

	tbl.SetTitle(title)
	tbl.SetSection1(section1)
	tbl.SetSection2(section2)
	tbl.SetSection3(section3)
	tbl.AddColumn("Name", 35, CELLSTRING, COLJUSTIFYLEFT)               // 0 Name
	tbl.AddColumn("Person Age", 3, CELLINT, COLJUSTIFYRIGHT)            // 1 Age
	tbl.AddColumn("Height (cm)", 0, CELLINT, COLJUSTIFYRIGHT)           // 2 Height in centimeters
	tbl.AddColumn("Date of Birth", 10, CELLDATE, COLJUSTIFYLEFT)        // 3 DOB
	tbl.AddColumn("Country of Birth", 14, CELLSTRING, COLJUSTIFYLEFT)   // 4 COB
	tbl.AddColumn("Winnings", 12, CELLFLOAT, COLJUSTIFYRIGHT)           // 5 total winnings
	tbl.AddColumn("Notes", 20, CELLSTRING, COLJUSTIFYLEFT)              // 6 Notes
	tbl.AddColumn("Random Date/Time", 25, CELLDATETIME, COLJUSTIFYLEFT) // 7 totally random datetime

	errExp = "No Records"
	// headers check
	err = tbl.HasData()
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}
	// print table check with text version
	tableOutText := fmt.Sprintf("%s", tbl)
	if !strings.Contains(tableOutText, errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}
	// csv output
	err = tbl.CSVprintTable(&temp)
	if !strings.Contains(temp.String(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}
	// html output
	err = tbl.HTMLprintTable(&temp)
	if !strings.Contains(temp.String(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}

	const (
		Name = iota
		Age
		Height
		DOB
		COB
		Winnings
		Notest
	)

	bd1 := time.Date(1969, time.March, 2, 0, 0, 0, 0, time.UTC)
	bd2 := time.Date(1960, time.October, 4, 0, 0, 0, 0, time.UTC)
	bd3 := time.Date(1974, time.April, 10, 0, 0, 0, 0, time.UTC)
	bd4 := time.Date(1950, time.April, 21, 0, 0, 0, 0, time.UTC)
	bd5 := time.Date(1977, time.August, 6, 0, 0, 0, 0, time.UTC)

	type tdata struct {
		Name     string
		Age      int64
		Height   int64
		DOB      time.Time
		COB      string
		Winnings float64
		Notes    string
		UnxNano  int64
	}
	var d = []tdata{
		{Name: "Mary M. Oneil", Age: 47, Height: 165, DOB: bd1, COB: "United States", Winnings: float64(17633.21), Notes: "A few notes here withaverylongnoteword", UnxNano: 7564006824999651664},
		{Name: "Lynette C. Allen", Age: 56, Height: 156, DOB: bd2, COB: "United States", Winnings: float64(45373.00), Notes: "A lot more notes. A whole, big, line with lots and lots and lots and lots of notes. And some more notes.", UnxNano: 7733402883116878723},
		{Name: "Stanislaus Aliyeva", Age: 42, Height: 172, DOB: bd3, COB: "Slovinia", Winnings: 106632.36, Notes: "A few notes here", UnxNano: 1584693382958379231},
		{Name: "Casandra Åberg", Age: 66, Height: 158, DOB: bd4, COB: "Sweden", Winnings: 93883.25, Notes: "2000 Seat Toledo", UnxNano: 7796987096200859545},
		{Name: "Amanda Melo Ferreira", Age: 55, Height: 174, DOB: bd5, COB: "Brazil", Winnings: 46673.42, Notes: "2006 Ford Falcon", UnxNano: 3267110399458248377},
	}

	totalsRSet := tbl.CreateRowset()
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < len(d); i++ {
		tbl.AddRow()
		tbl.AppendToRowset(totalsRSet, tbl.RowCount()-1)
		tbl.Puts(-1, Name, d[i].Name)
		tbl.Puti(-1, Age, d[i].Age)
		tbl.Puti(-1, Height, d[i].Height)
		tbl.Putd(-1, DOB, d[i].DOB)
		tbl.Puts(-1, COB, d[i].COB)
		tbl.Putf(-1, Winnings, d[i].Winnings)
		tbl.Puts(-1, Notest, d[i].Notes)
		tbl.Putdt(-1, 7, time.Unix(0, d[i].UnxNano).UTC()) // random date in the future
	}
	// Start with simple checks...
	if tbl.ColCount() != 8 {
		t.Errorf("smoke_test: Expected %d,  found %d\n", 8, tbl.ColCount())
	}
	if len(tbl.GetRowset(99)) != 0 {
		t.Errorf("smoke_test: Expected emty rowset,  found %#v\n", tbl.GetRowset(99))
	}
	if 0 != tbl.Type(999, 999) {
		t.Errorf("smoke_test: Expected %d,  found %d\n", 0, tbl.Type(999, 999))
	}
	if tbl.GetTitle() != title {
		t.Errorf("smoke_test: Expected %s,  found %s\n", tbl.GetTitle(), title)
	}
	if false != tbl.Puti(999, 999, 1) {
		t.Errorf("smoke_test: Expected return value of false, but got true\n")
	}
	if false != tbl.Putf(999, 999, 1) {
		t.Errorf("smoke_test: Expected return value of false, but got true\n")
	}
	if false != tbl.Puts(999, 999, "ignore") {
		t.Errorf("smoke_test: Expected return value of false, but got true\n")
	}
	if false != tbl.Putd(999, 999, time.Now()) {
		t.Errorf("smoke_test: Expected return value of false, but got true\n")
	}
	if tbl.GetSection1() != section1 {
		t.Errorf("smoke_test: Expected %s,  found %s\n", tbl.GetSection1(), section1)
	}
	if tbl.GetSection2() != section2 {
		t.Errorf("smoke_test: Expected %s,  found %s\n", tbl.GetSection2(), section2)
	}
	if tbl.GetSection3() != section3 {
		t.Errorf("smoke_test: Expected %s,  found %s\n", tbl.GetSection3(), section3)
	}
	cell := tbl.Get(0, 0)
	if cell.Sval != d[0].Name {
		t.Errorf("smoke_test: Expected %s,  found %s\n", cell.Sval, d[0].Name)
	}
	if tbl.Geti(1, Age) != d[1].Age {
		t.Errorf("smoke_test: Expected %d,  found %d\n", tbl.Geti(1, Age), d[1].Age)
	}
	if tbl.Getf(1, Winnings) != d[1].Winnings {
		t.Errorf("smoke_test: Expected %f,  found %f\n", tbl.Getf(1, Winnings), d[1].Winnings)
	}
	if tbl.Gets(1, Name) != d[1].Name {
		t.Errorf("smoke_test: Expected %s,  found %s\n", tbl.Gets(1, Name), d[1].Name)
	}
	if tbl.Getd(1, DOB) != d[1].DOB {
		t.Logf("smoke_test: Expected %s,  found %s\n", tbl.Getd(1, DOB).Format("1/2/2006"), d[1].DOB.Format("1/Errorf06"))
	}
	if tbl.Type(1, Name) != CELLSTRING {
		t.Errorf("smoke_test: Expected %d,  found %d\n", tbl.Type(1, Name), CELLSTRING)
	}

	// Bang it a bit...
	tbl.Sort(0, tbl.RowCount()-1, DOB)
	tbl.AddLineAfter(tbl.RowCount() - 1) // a line after the last row in the table
	tbl.AddLineBefore(tbl.RowCount())    // a line after the last row in the table
	tbl.InsertSumRowsetCols(totalsRSet, tbl.RowCount(), []int{Winnings})

	// css property on table for html
	cssSmpl := &CSSProperty{Name: "color", Value: "orange"}
	cssStrExp := `"color:orange;"`
	cStrOut := fmt.Sprintf("%s", cssSmpl)
	if cssStrExp != cStrOut {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", cssStrExp, cStrOut)
	}

	cssList := []*CSSProperty{}
	cssList = append(cssList, cssSmpl)

	// --------------------------------
	// validation over rows
	// --------------------------------
	errExp = "Row number > no of rows in table"
	err = tbl.HasValidRow(999)
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}
	err = tbl.SetRowCSS(999, cssList)
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}
	err = tbl.SetCellCSS(999, tbl.ColCount()-1, cssList)
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}

	errExp = "Row number is less than zero"
	err = tbl.HasValidRow(-999)
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}
	err = tbl.SetRowCSS(-999, cssList)
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}
	err = tbl.SetCellCSS(-999, tbl.ColCount()-1, cssList)
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}

	// valid row case
	if err = tbl.HasValidRow(tbl.RowCount() - 1); err != nil {
		t.Errorf("smoke_test: Expected `nil` Error, but found: %s\n", err.Error())
	}

	// --------------------------------
	// validation over columns
	// --------------------------------
	errExp = "Column number > no of columns in table"
	err = tbl.HasValidColumn(999)
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}
	err = tbl.SetHeaderCellCSS(999, cssList)
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}
	err = tbl.SetColCSS(999, cssList)
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}
	err = tbl.SetCellCSS(tbl.RowCount()-1, 999, cssList)
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}
	err = tbl.SetColHTMLWidth(999, 200, "px")
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}

	errExp = "Column number is less than zero"
	err = tbl.HasValidColumn(-999)
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}
	err = tbl.SetHeaderCellCSS(-999, cssList)
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}
	err = tbl.SetColCSS(-999, cssList)
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}
	err = tbl.SetCellCSS(tbl.RowCount()-1, -999, cssList)
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}
	err = tbl.SetColHTMLWidth(-999, 200, "px")
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}

	// valid column case
	if err = tbl.HasValidColumn(tbl.ColCount() - 1); err != nil {
		t.Errorf("smoke_test: Expected `nil` Error, but found: %s\n", err.Error())
	}

	// valid row css - make first row text color orange
	if err = tbl.SetRowCSS(0, cssList); err != nil {
		t.Errorf("smoke_test: Expected `nil` Error, but found: %s\n", err.Error())
	}
	// valid column css - make first column text color oranage
	cssList = []*CSSProperty{}
	cssList = append(cssList, &CSSProperty{Name: "color", Value: "blue"})
	if err = tbl.SetColCSS(0, cssList); err != nil {
		t.Errorf("smoke_test: Expected `nil` Error, but found: %s\n", err.Error())
	}
	// valid cell css case - make all cells background color yellow
	cssList = []*CSSProperty{}
	cssList = append(cssList, &CSSProperty{Name: "background-color", Value: "yellow"})
	tbl.SetAllCellCSS(cssList)

	// valid set html width case - make second column width wider
	if err = tbl.SetColHTMLWidth(1, 10, "ch"); err != nil {
		t.Errorf("smoke_test: Expected `nil` Error, but found: %s\n", err.Error())
	}

	//set title css
	cssList = []*CSSProperty{}
	cssList = append(cssList, &CSSProperty{Name: "color", Value: "blue"})
	cssList = append(cssList, &CSSProperty{Name: "font-style", Value: "italic"})
	// cssList = append(cssList, &CSSProperty{Name: "font-size", Value: "20px"})
	tbl.SetTitleCSS(cssList)

	// set header css
	cssList = []*CSSProperty{}
	cssList = append(cssList, &CSSProperty{Name: "color", Value: "orange"})
	cssList = append(cssList, &CSSProperty{Name: "font-style", Value: "italic"})
	tbl.SetHeaderCSS(cssList)
	cssList = append(cssList, &CSSProperty{Name: "background-color", Value: "blue"})
	tbl.SetHeaderCSS(cssList)

	// set section1 css
	cssList = []*CSSProperty{}
	cssList = append(cssList, &CSSProperty{Name: "color", Value: "white"})
	cssList = append(cssList, &CSSProperty{Name: "background-color", Value: "black"})
	tbl.SetSection1CSS(cssList)

	// set section2 css
	cssList = []*CSSProperty{}
	cssList = append(cssList, &CSSProperty{Name: "color", Value: "red"})
	cssList = append(cssList, &CSSProperty{Name: "background-color", Value: "yellow"})
	tbl.SetSection2CSS(cssList)

	// set section3 css
	cssList = []*CSSProperty{}
	cssList = append(cssList, &CSSProperty{Name: "color", Value: "green"})
	tbl.SetSection3CSS(cssList)

	SetLogger(ioutil.Discard, "debug")

	// errorneous pdf output, cover the errorneous output
	wrongProps := []*PDFProperty{{Option: "-fake"}}
	err = tbl.PDFprintTable(ioutil.Discard, wrongProps)
	if err == nil {
		t.Error("smoke_test: expected error while using with fake options for pdf output")
	}

	// Now hit it hard...
	DoTextOutput(t, &tbl)
	DoCSVOutput(t, &tbl)
	DoHTMLOutput(t, &tbl)
	DoPDFOutput(t, &tbl)
	DoCustomTemplateHTMLOutput(t, &tbl)

	// multiple table test
	var m []Table
	for i := 0; i < 10; i++ {
		var t Table
		t = tbl
		t.TightenColumns()
		m = append(m, t)
	}

	DoMultiTableTextOutput(t, &m)
	DoMultiTableCSVOutput(t, &m)
	DoMultiTableHTMLOutput(t, &m)
	DoMultiTablePDFOutput(t, &m)
}

func DoTextOutput(t *testing.T, tbl *Table) {
	(*tbl).TightenColumns()

	fname := "smoke_test.txt"
	f, err := os.Create(fname)
	if nil != err {
		t.Errorf("smoke_test: Error creating file %s: %s\n", fname, err.Error())
		// fmt.Printf("smoke_test: Error creating file: %s\n", err.Error())
	}

	if err := tbl.TextprintTable(f); err != nil {
		t.Errorf("smoke_test: Error creating TEXT output: %s\n", err.Error())
	}
	// close file after operation
	f.Close()

	// now compare what we have to the known-good output
	b, _ := ioutil.ReadFile("./testdata/smoke_test.txt")
	sb, _ := ioutil.ReadFile("./smoke_test.txt")

	if len(b) != len(sb) {
		// fmt.Printf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
		t.Errorf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
	}
	if len(sb) > 0 && len(b) > 0 {
		for i := 0; i < len(b); i++ {
			if i < len(sb) && sb[i] != b[i] {
				t.Errorf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
				// fmt.Printf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
				break
			}
		}
	}

	// add test case for SprintTable, it should be same as TextPrintTable output
	s, err := tbl.SprintTable()
	st := []byte(s)
	if err != nil {
		t.Errorf("smoke_test: Error creating TEXT output: %s\n", err.Error())
	}
	if len(st) != len(sb) {
		t.Errorf("smoke_test: Expected len = %d,  found len = %d\n", len(st), len(sb))
	}
	if len(sb) > 0 && len(st) > 0 {
		for i := 0; i < len(st); i++ {
			if i < len(sb) && sb[i] != st[i] {
				t.Errorf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, st[i], st[i], sb[i], sb[i])
				break
			}
		}
	}

	// add test case for String, it should be same as TextPrintTable output
	s = tbl.String()
	st = []byte(s)
	if err != nil {
		t.Errorf("smoke_test: Error creating TEXT output: %s\n", err.Error())
	}
	if len(st) != len(sb) {
		t.Errorf("smoke_test: Expected len = %d,  found len = %d\n", len(st), len(sb))
	}
	if len(sb) > 0 && len(st) > 0 {
		for i := 0; i < len(st); i++ {
			if i < len(sb) && sb[i] != st[i] {
				t.Errorf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, st[i], st[i], sb[i], sb[i])
				break
			}
		}
	}

	// add test case for FprintTable, it should be same as TextPrintTable output
	var temp bytes.Buffer
	err = tbl.FprintTable(&temp)
	st = temp.Bytes()
	if err != nil {
		t.Errorf("smoke_test: Error creating TEXT output: %s\n", err.Error())
	}
	if len(st) != len(sb) {
		t.Errorf("smoke_test: Expected len = %d,  found len = %d\n", len(st), len(sb))
	}
	if len(sb) > 0 && len(st) > 0 {
		for i := 0; i < len(st); i++ {
			if i < len(sb) && sb[i] != st[i] {
				t.Errorf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, st[i], st[i], sb[i], sb[i])
				break
			}
		}
	}
}

func DoCSVOutput(t *testing.T, tbl *Table) {
	fname := "smoke_test.csv"
	f, err := os.Create(fname)
	if nil != err {
		t.Errorf("smoke_test: Error creating file %s: %s\n", fname, err.Error())
		// fmt.Printf("smoke_test: Error creating file: %s\n", err.Error())
	}

	if err := tbl.CSVprintTable(f); err != nil {
		t.Errorf("smoke_test: Error creating CSV output: %s\n", err.Error())
	}
	// close file after operation
	f.Close()

	// now compare what we have to the known-good output
	b, _ := ioutil.ReadFile("./testdata/smoke_test.csv")
	sb, _ := ioutil.ReadFile("./smoke_test.csv")

	if len(b) != len(sb) {
		// fmt.Printf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
		t.Errorf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
	}
	if len(sb) > 0 && len(b) > 0 {
		for i := 0; i < len(b); i++ {
			if i < len(sb) && sb[i] != b[i] {
				t.Logf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
				// fmt.Printf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
				break
			}
		}
	}
}

func DoHTMLOutput(t *testing.T, tbl *Table) {
	fname := "smoke_test.html"
	f, err := os.Create(fname)
	if nil != err {
		t.Errorf("smoke_test: Error creating file %s: %s\n", fname, err.Error())
		// fmt.Printf("smoke_test: Error creating file: %s\n", err.Error())
	}

	if err := tbl.HTMLprintTable(f); err != nil {
		t.Errorf("smoke_test: Error creating HTML output: %s\n", err.Error())
	}
	// close file after operation
	f.Close()

	// now compare what we have to the known-good output
	b, _ := ioutil.ReadFile("./testdata/smoke_test.html")
	sb, _ := ioutil.ReadFile("./smoke_test.html")

	if len(b) != len(sb) {
		// fmt.Printf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
		t.Errorf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
	}
	if len(sb) > 0 && len(b) > 0 {
		for i := 0; i < len(b); i++ {
			if i < len(sb) && sb[i] != b[i] {
				t.Logf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
				// fmt.Printf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
				break
			}
		}
	}
}

func DoCustomTemplateHTMLOutput(t *testing.T, tbl *Table) {
	_, filename, _, _ := runtime.Caller(0)
	customDir := path.Dir(filename)
	tbl.SetHTMLTemplate(path.Join(path.Join(customDir, "testdata"), "custom.tmpl"))
	tbl.SetHTMLTemplateCSS(path.Join(path.Join(customDir, "testdata"), "custom.css"))
	fname := "smoke_test_custom_template.html"
	f, err := os.Create(fname)
	if nil != err {
		t.Errorf("smoke_test: Error creating file %s: %s\n", fname, err.Error())
		// fmt.Printf("smoke_test: Error creating file: %s\n", err.Error())
	}

	if err := tbl.HTMLprintTable(f); err != nil {
		t.Errorf("smoke_test: Error creating HTML output: %s\n", err.Error())
	}
	// close file after operation
	f.Close()

	// now compare what we have to the known-good output
	b, _ := ioutil.ReadFile("./testdata/smoke_test_custom_template.html")
	sb, _ := ioutil.ReadFile("./smoke_test_custom_template.html")

	if len(b) != len(sb) {
		// fmt.Printf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
		t.Errorf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
	}
	if len(sb) > 0 && len(b) > 0 {
		for i := 0; i < len(b); i++ {
			if i < len(sb) && sb[i] != b[i] {
				t.Logf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
				// fmt.Printf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
				break
			}
		}
	}
}

func DoPDFOutput(t *testing.T, tbl *Table) {
	fname := "smoke_test.pdf"
	f, err := os.Create(fname)
	if nil != err {
		t.Errorf("smoke_test: Error creating file %s: %s\n", fname, err.Error())
		// fmt.Printf("smoke_test: Error creating file: %s\n", err.Error())
	}

	if err := tbl.PDFprintTable(f, pdfProps); err != nil {
		t.Errorf("smoke_test: Error creating PDF output: %s\n", err.Error())
	}
	// close file after operation
	f.Close()

	// now compare what we have to the known-good output
	// b, _ := ioutil.ReadFile("./testdata/smoke_test.pdf")
	sb, _ := ioutil.ReadFile("./smoke_test.pdf")

	if len(sb) == 0 {
		t.Errorf("smoke_test: Expected some content in PDF output file,  found len = 0")
	}

	// if len(b) != len(sb) {
	// 	// fmt.Printf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
	// 	t.Errorf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
	// }
	// if len(sb) > 0 && len(b) > 0 {
	// 	for i := 0; i < len(b); i++ {
	// 		if i < len(sb) && sb[i] != b[i] {
	// 			t.Logf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
	// 			// fmt.Printf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
	// 			break
	// 		}
	// 	}
	// }
}

func saveTableToFile(t *testing.T, fname string, s string) error {
	// save for later inspection if anything goes wrong
	f, err := os.Create(fname)
	if nil != err {
		t.Errorf("smoke_test: Error creating file %s: %s\n", fname, err.Error())
		// fmt.Printf("smoke_test: Error creating file: %s\n", err.Error())
	}
	defer f.Close()
	fmt.Fprintf(f, "%s", s)
	return err
}

func DoMultiTableTextOutput(t *testing.T, m *[]Table) {

	fname := "smoke_multi_table_test.txt"
	f, err := os.Create(fname)
	if nil != err {
		t.Errorf("smoke_test: Error creating file %s: %s\n", fname, err.Error())
		// fmt.Printf("smoke_test: Error creating file: %s\n", err.Error())
	}

	if err := MultiTableTextPrint(*m, f); err != nil {
		t.Errorf("smoke_test: Error creating TEXT output: %s\n", err.Error())
	}
	// close file after operation
	f.Close()

	// now compare what we have to the known-good output
	b, _ := ioutil.ReadFile("./testdata/smoke_multi_table_test.txt")
	sb, _ := ioutil.ReadFile("./smoke_multi_table_test.txt")

	if len(b) != len(sb) {
		// fmt.Printf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
		t.Errorf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
	}
	if len(sb) > 0 && len(b) > 0 {
		for i := 0; i < len(b); i++ {
			if i < len(sb) && sb[i] != b[i] {
				t.Errorf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
				// fmt.Printf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
				break
			}
		}
	}
}

func DoMultiTableCSVOutput(t *testing.T, m *[]Table) {

	fname := "smoke_multi_table_test.csv"
	f, err := os.Create(fname)
	if nil != err {
		t.Errorf("smoke_test: Error creating file %s: %s\n", fname, err.Error())
		// fmt.Printf("smoke_test: Error creating file: %s\n", err.Error())
	}

	if err := MultiTableCSVPrint(*m, f); err != nil {
		t.Errorf("smoke_test: Error creating CSV output: %s\n", err.Error())
	}
	// close file after operation
	f.Close()

	// now compare what we have to the known-good output
	b, _ := ioutil.ReadFile("./testdata/smoke_multi_table_test.csv")
	sb, _ := ioutil.ReadFile("./smoke_multi_table_test.csv")

	if len(b) != len(sb) {
		// fmt.Printf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
		t.Errorf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
	}
	if len(sb) > 0 && len(b) > 0 {
		for i := 0; i < len(b); i++ {
			if i < len(sb) && sb[i] != b[i] {
				t.Errorf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
				// fmt.Printf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
				break
			}
		}
	}
}

func DoMultiTableHTMLOutput(t *testing.T, m *[]Table) {

	fname := "smoke_multi_table_test.html"
	f, err := os.Create(fname)
	if nil != err {
		t.Errorf("smoke_test: Error creating file %s: %s\n", fname, err.Error())
		// fmt.Printf("smoke_test: Error creating file: %s\n", err.Error())
	}

	if err := MultiTableHTMLPrint(*m, f); err != nil {
		t.Errorf("smoke_test: Error creating HTML output: %s\n", err.Error())
	}
	// close file after operation
	f.Close()

	// now compare what we have to the known-good output
	b, _ := ioutil.ReadFile("./testdata/smoke_multi_table_test.html")
	sb, _ := ioutil.ReadFile("./smoke_multi_table_test.html")

	if len(b) != len(sb) {
		// fmt.Printf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
		t.Errorf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
	}
	if len(sb) > 0 && len(b) > 0 {
		for i := 0; i < len(b); i++ {
			if i < len(sb) && sb[i] != b[i] {
				t.Errorf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
				// fmt.Printf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
				break
			}
		}
	}
}

func DoMultiTablePDFOutput(t *testing.T, m *[]Table) {

	fname := "smoke_multi_table_test.pdf"
	f, err := os.Create(fname)
	if nil != err {
		t.Errorf("smoke_test: Error creating file %s: %s\n", fname, err.Error())
		// fmt.Printf("smoke_test: Error creating file: %s\n", err.Error())
	}

	if err := MultiTablePDFPrint(*m, f, pdfProps); err != nil {
		t.Errorf("smoke_test: Error creating PDF output: %s\n", err.Error())
	}
	// close file after operation
	f.Close()

	// now compare what we have to the known-good output
	// b, _ := ioutil.ReadFile("./testdata/smoke_multi_table_test.pdf")
	sb, _ := ioutil.ReadFile("./smoke_multi_table_test.pdf")

	if len(sb) == 0 {
		t.Errorf("smoke_test: Expected some content in PDF output file,  found len = 0")
	}

	// if len(b) != len(sb) {
	// 	// fmt.Printf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
	// 	t.Errorf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
	// }
	// if len(sb) > 0 && len(b) > 0 {
	// 	for i := 0; i < len(b); i++ {
	// 		if i < len(sb) && sb[i] != b[i] {
	// 			t.Logf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
	// 			// fmt.Printf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
	// 			break
	// 		}
	// 	}
	// }
}
