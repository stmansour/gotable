package gotable

import (
	"bytes"
	"fmt"
	"io"
	"sort"

	"github.com/dustin/go-humanize"
)

// TextTable struct used to prepare table in text version
type TextTable struct {
	*Table
	TextColSpace int
	buf          bytes.Buffer
}

func (tt *TextTable) writeTableOutput(w io.Writer) error {

	// vars
	var (
		err error
	)

	// append title
	tt.buf.WriteString(tt.formatTitle())

	// append section 1
	tt.buf.WriteString(tt.formatSection1())

	// append section 2
	tt.buf.WriteString(tt.formatSection2())

	// append section 3
	tt.buf.WriteString(tt.formatSection3())

	// append headers
	if headerStr, err := tt.formatHeaders(); err != nil {
		tt.buf.WriteString(stringln(err.Error()))
	} else {
		// append rows
		if rowsStr, err := tt.formatRows(); err != nil {
			tt.buf.WriteString(stringln(err.Error()))
		} else {
			// if rows exist, then only show headers
			tt.buf.WriteString(headerStr)
			tt.buf.WriteString(rowsStr)
		}
	}

	// // render errorlist
	// NOTE: if you enable this errorList feature then write them first on top,
	// then write headers, rows output
	// ct.buf.Write(ct.getErrorSection())

	// write output to passed io.Writer interface object
	w.Write(tt.buf.Bytes())
	return err
}

func (tt *TextTable) formatTitle() string {
	title := tt.Table.GetTitle()
	if title != "" {
		return stringln(title)
	}
	return title
}

func (tt *TextTable) formatSection1() string {
	section1 := tt.Table.GetSection1()
	if section1 != "" {
		return stringln(section1)
	}
	return section1
}

func (tt *TextTable) formatSection2() string {
	section2 := tt.Table.GetSection2()
	if section2 != "" {
		return stringln(section2)
	}
	return section2
}

func (tt *TextTable) formatSection3() string {
	section3 := tt.Table.GetSection3()
	if section3 != "" {
		return stringln(section3)
	}
	return section3
}

// func (tt *TextTable) getErrorSection() string {
// 	errSection := ""

// 	errList := tt.Table.GetErrorList()
// 	if len(errList) > 0 {
// 		for _, errStr := range errList {
// 			errSection += stringln(errStr)
// 		}
// 		errSection = NEWLINE + errSection + NEWLINE
// 	}

// 	// blank return
// 	return errSection
// }

// SprintColHdrsText formats the column headers as text and returns the string
func (tt *TextTable) formatHeaders() (string, error) {

	// check for blank headers
	blankHdrsErr := tt.Table.HasHeaders()
	if blankHdrsErr != nil {
		return "", blankHdrsErr
	}

	tt.Table.AdjustAllColumnHeaders()

	var s bytes.Buffer

	for j := 0; j < len(tt.Table.ColDefs[0].Hdr); j++ {
		for i := 0; i < len(tt.Table.ColDefs); i++ {
			sf := ""
			lft := ""
			if tt.Table.ColDefs[i].Justify == COLJUSTIFYLEFT {
				lft += "-"
			}
			sf += fmt.Sprintf("%%%s%ds", lft, tt.Table.ColDefs[i].Width)
			s.WriteString(fmt.Sprintf(sf, tt.Table.ColDefs[i].Hdr[j]))
			s.WriteString(mkstr(tt.TextColSpace, ' '))
		}

		// remove last textColSpace from s
		tmp := s.Bytes()
		tmp = tmp[0 : len(tmp)-tt.TextColSpace]
		s.Reset()
		s.Write(tmp)

		// append new line after first line of grid
		s.WriteByte('\n')
	}

	// finally append separator with line
	s.WriteString(tt.sprintLineText())

	return s.String(), nil
}

func (tt *TextTable) formatRows() (string, error) {

	// check for empty data table
	blankDataErr := tt.Table.HasData()
	if blankDataErr != nil {
		return "", blankDataErr
	}

	var rowsOut bytes.Buffer
	for i := 0; i < tt.Table.RowCount(); i++ {
		// for valid row, we will never get an error
		s, _ := tt.formatRow(i)
		rowsOut.WriteString(s)
	}

	return rowsOut.String(), nil
}

func (tt *TextTable) formatRow(row int) (string, error) {

	// This method is only called by internal instance of TextTable
	// in formatRows method, so we should avoid following error check
	// unless we make it as export

	// check that this passed row is valid or not
	// inValidRowErr := tt.Table.HasValidRow(row)
	// if inValidRowErr != nil {
	// 	return "", inValidRowErr
	// }

	// format table row
	var s bytes.Buffer

	if len(tt.Table.LineBefore) > 0 {
		j := sort.SearchInts(tt.Table.LineBefore, row)
		// line separator added in `LineAfter`??
		// If YES, then discard it
		sepExist := sort.SearchInts(tt.Table.LineAfter, row-1) < tt.Table.RowCount()

		if j < len(tt.Table.LineBefore) && row == tt.Table.LineBefore[j] && !sepExist {
			s.WriteString(tt.sprintLineText())
		}
	}

	rowColumns := tt.Table.ColCount()

	// columns string chunk map, each column holds list of string
	// used for multi line text
	colMultiLineTextMap := map[int][]string{}

	// get Height of row that require to fit the content of max cell string content
	// by default table has no all the data in string format, so that we need to add
	// logic here only, to support multi line functionality
	for gridColIndex := 0; gridColIndex < rowColumns; gridColIndex++ {
		if tt.Table.Row[row].Col[gridColIndex].Type == CELLSTRING {
			cd := tt.Table.ColDefs[gridColIndex]

			// get multi line text
			a, _ := getMultiLineText(tt.Table.Row[row].Col[gridColIndex].Sval, cd.Width)

			// store multi line text list in column multi line text map
			colMultiLineTextMap[gridColIndex] = a

			// if greater value found then store it
			if len(a) > tt.Table.Row[row].Height {
				tt.Table.Row[row].Height = len(a)
			}
		}
	}

	rowHeight := tt.Table.Row[row].Height

	// rowGrid holds grid for row with multi line text
	// NOTE: Non constant bound array error
	// cannot create with runtime variable value
	rowGrid := [][]string{}

	// fill grid with empty whitespace value so that it can hold proper spacing
	// to fit the row in table text output
	emptyCols := make([]string, rowColumns)
	for gridColIndex := 0; gridColIndex < rowColumns; gridColIndex++ {
		// assign default string with length of column width
		emptyCols[gridColIndex] = mkstr(tt.Table.ColDefs[gridColIndex].Width, ' ')
	}
	// fit these prepared empty column list in rowGrid for each row
	for gridRowIndex := 0; gridRowIndex < rowHeight; gridRowIndex++ {
		rowGrid = append(rowGrid, emptyCols)
	}

	// for the first line in grid fill all type of data in it
	// for string type take it from col multi line text map first chunk
	// FIRST LINE OF ROW GRID
	for gridColIndex := 0; gridColIndex < rowColumns; gridColIndex++ {
		switch tt.Table.Row[row].Col[gridColIndex].Type {
		case CELLFLOAT:
			s.WriteString(fmt.Sprintf(tt.Table.ColDefs[gridColIndex].Pfmt, humanize.FormatFloat("#,###.##", tt.Table.Row[row].Col[gridColIndex].Fval)))
		case CELLINT:
			s.WriteString(fmt.Sprintf(tt.Table.ColDefs[gridColIndex].Pfmt, tt.Table.Row[row].Col[gridColIndex].Ival))
		case CELLSTRING:
			s.WriteString(fmt.Sprintf(tt.Table.ColDefs[gridColIndex].Pfmt, colMultiLineTextMap[gridColIndex][0]))
		case CELLDATE:
			s.WriteString(fmt.Sprintf("%*.*s", tt.Table.ColDefs[gridColIndex].Width, tt.Table.ColDefs[gridColIndex].Width, tt.Table.Row[row].Col[gridColIndex].Dval.Format(tt.Table.DateFmt)))
		case CELLDATETIME:
			s.WriteString(fmt.Sprintf("%*.*s", tt.Table.ColDefs[gridColIndex].Width, tt.Table.ColDefs[gridColIndex].Width, tt.Table.Row[row].Col[gridColIndex].Dval.Format(tt.Table.DateTimeFmt)))
		default:
			s.WriteString(mkstr(tt.Table.ColDefs[gridColIndex].Width, ' '))
		}

		// append text col whitespace
		s.WriteString(mkstr(tt.TextColSpace, ' '))
	}

	// remove last textColSpace from s
	tmp := s.Bytes()
	tmp = tmp[0 : len(tmp)-tt.TextColSpace]
	s.Reset()
	s.Write(tmp)

	// append new line after first line of grid
	s.WriteByte('\n')

	// now proceed with rest of the line in row grid
	// for multi line text
	for gridRowIndex := 1; gridRowIndex < rowHeight; gridRowIndex++ {

		for gridColIndex := 0; gridColIndex < rowColumns; gridColIndex++ {

			if tt.Table.Row[row].Col[gridColIndex].Type == CELLSTRING {
				if gridRowIndex >= len(colMultiLineTextMap[gridColIndex]) {
					rowGrid[gridRowIndex][gridColIndex] = fmt.Sprintf(tt.Table.ColDefs[gridColIndex].Pfmt, "")
				} else {
					rowGrid[gridRowIndex][gridColIndex] = fmt.Sprintf(tt.Table.ColDefs[gridColIndex].Pfmt, colMultiLineTextMap[gridColIndex][gridRowIndex])
				}
			}

			// write string in the cell
			s.WriteString(rowGrid[gridRowIndex][gridColIndex])

			// append text col whitespace
			s.WriteString(mkstr(tt.TextColSpace, ' '))
		}

		// remove last textColSpace from s
		tmp := s.Bytes()
		tmp = tmp[0 : len(tmp)-tt.TextColSpace]
		s.Reset()
		s.Write(tmp)

		// append new line
		s.WriteByte('\n')
	}

	if len(tt.Table.LineAfter) > 0 {
		j := sort.SearchInts(tt.Table.LineAfter, row)
		if j < len(tt.Table.LineAfter) && row == tt.Table.LineAfter[j] {
			s.WriteString(tt.sprintLineText())
		}
	}
	return s.String(), nil
}

// SprintLineText returns a line across all rows in the table
func (tt *TextTable) sprintLineText() string {
	var s string
	for i := 0; i < len(tt.Table.ColDefs); i++ {
		// draw line with hyphen `-` char
		s += mkstr(tt.Table.ColDefs[i].Width, '-')

		// separate text columns
		s += mkstr(tt.TextColSpace, ' ')
	}
	// remove last textColSpace from s
	s = s[0 : len(s)-tt.TextColSpace]
	return stringln(s)
}
