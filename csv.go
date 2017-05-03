package gotable

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/dustin/go-humanize"
)

// CSVTable struct used to prepare table in html version
type CSVTable struct {
	*Table
	buf *csv.Writer
}

func (ct *CSVTable) writeTableOutput(w io.Writer) error {

	// vars
	var (
		err error
	)

	// get new writer for w io.Writer and assign it to ct.buf
	ct.buf = csv.NewWriter(w)

	// write title
	ct.writeTitle()

	// write section 1
	ct.writeSection1()

	// write section 2
	ct.writeSection2()

	// write section 3
	ct.writeSection3()

	// append headers and rows
	if headers, err := ct.formatHeaders(); err != nil {
		errHeaderRow := []string{err.Error()}
		ct.buf.Write(errHeaderRow)
	} else {

		// append rows
		if rows, err := ct.formatRows(); err != nil {
			errDataRow := []string{err.Error()}
			ct.buf.Write(errDataRow)
		} else {
			// write one header row
			ct.buf.Write(headers)

			// rows holds 2d slices
			ct.buf.WriteAll(rows)
		}
	}

	// // render errorlist
	// NOTE: if you enable this errorList feature then write them first on top,
	// then write headers, rows output
	// ct.buf.Write(ct.getErrorSection())

	// Write any buffered data to the underlying writer (standard output).
	ct.buf.Flush()

	if err = ct.buf.Error(); err != nil {
		return err
	}

	return err
}

func (ct *CSVTable) writeTitle() {
	var title []string
	if ct.Table.GetTitle() != "" {
		title = append(title, ct.Table.GetTitle())
		ct.buf.Write(title)
	}
}

func (ct *CSVTable) writeSection1() {
	var section1 []string
	if ct.Table.GetSection1() != "" {
		section1 = append(section1, ct.Table.GetSection1())
		ct.buf.Write(section1)
	}
}

func (ct *CSVTable) writeSection2() {
	var section2 []string
	if ct.Table.GetSection2() != "" {
		section2 = append(section2, ct.Table.GetSection2())
		ct.buf.Write(section2)
	}
}

func (ct *CSVTable) writeSection3() {
	var section3 []string
	if ct.Table.GetSection3() != "" {
		section3 = append(section3, ct.Table.GetSection3())
		ct.buf.Write(section3)
	}
}

// func (ct *CSVTable) getErrorSection() []string {
// 	var errSection []string

// 	errList := ct.Table.GetErrorList()

// 	if len(errList) > 0 {
// 		for _, errStr := range errList {
// 			errSection = append(errSection, errStr)
// 		}
// 	}

// 	// blank return
// 	return errSection
// }

func (ct *CSVTable) formatHeaders() ([]string, error) {
	var tHeaders []string

	// check for blank headers
	blankHdrsErr := ct.Table.HasHeaders()
	if blankHdrsErr != nil {
		return tHeaders, blankHdrsErr
	}

	// format headers
	for i := 0; i < len(ct.Table.ColDefs); i++ {
		tHeaders = append(tHeaders, ct.Table.ColDefs[i].ColTitle)
	}

	// remove last cellSep character from tHeaders
	// join slice of string by CellSep (default -> ',') character
	return tHeaders, nil
}

func (ct *CSVTable) formatRows() ([][]string, error) {
	var rowsOut [][]string

	// check for empty data table
	blankDataErr := ct.Table.HasData()
	if blankDataErr != nil {
		return rowsOut, blankDataErr
	}

	for i := 0; i < ct.Table.RowCount(); i++ {
		// for valid row, we will never get an error
		s, _ := ct.formatRow(i)
		rowsOut = append(rowsOut, s)
	}

	return rowsOut, nil
}

func (ct *CSVTable) formatRow(row int) ([]string, error) {

	// This method is only called by internal instance of TextTable
	// in formatRows method, so we should avoid following error check
	// unless we make it as export

	// // check that this passed row is valid or not
	// inValidRowErr := ct.Table.HasValidRow(row)
	// if inValidRowErr != nil {
	// 	return "", inValidRowErr
	// }

	// format table row
	var tRow []string

	for i := 0; i < len(ct.Table.Row[row].Col); i++ {

		switch ct.Table.Row[row].Col[i].Type {
		case CELLFLOAT:
			tRow = append(tRow, fmt.Sprintf(ct.Table.ColDefs[i].Pfmt, humanize.FormatFloat("#,###.##", ct.Table.Row[row].Col[i].Fval)))
		case CELLINT:
			tRow = append(tRow, fmt.Sprintf(ct.Table.ColDefs[i].Pfmt, ct.Table.Row[row].Col[i].Ival))
		case CELLSTRING:
			// FOR CSV, APPEND FULL STRING, THERE ARE NO MULTILINE STRING IN THIS
			tRow = append(tRow, ct.Table.Row[row].Col[i].Sval)
		case CELLDATE:
			tRow = append(tRow, fmt.Sprintf("%*.*s", ct.Table.ColDefs[i].Width, ct.Table.ColDefs[i].Width, ct.Table.Row[row].Col[i].Dval.Format(ct.Table.DateFmt)))
		case CELLDATETIME:
			tRow = append(tRow, fmt.Sprintf("%*.*s", ct.Table.ColDefs[i].Width, ct.Table.ColDefs[i].Width, ct.Table.Row[row].Col[i].Dval.Format(ct.Table.DateTimeFmt)))
		default:
			tRow = append(tRow, mkstr(ct.Table.ColDefs[i].Width, ' '))
		}
	}

	// return
	return tRow, nil
}

// MultiTableCSVPrint writes csv output from each table to w io.Writer
func MultiTableCSVPrint(m []Table, w io.Writer) error {
	funcname := "MultiTableCSVPrint"

	for i := 0; i < len(m); i++ {
		temp := bytes.Buffer{}
		err := m[i].CSVprintTable(&temp)
		if err != nil {
			errorLog("%s: Error while getting table output, title: %s, err: %s", funcname, m[i].Title, err.Error())
			return err
		}
		temp.WriteByte('\n')
		w.Write(temp.Bytes())
	}

	return nil
}
