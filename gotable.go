package gotable

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Table is a simple skeletal row-column "class" for go that implements a few
// useful methods for building, maintaining, and printing tables of data.
// To use this table, you define all the columns first. Then call the AddRow
// method and begin adding and modifying values in the table cells.
//
// You can insert rows, append rows, sort all or selected rows by their column
// values.

// COLJUSTIFYLEFT et. al. are the constants used in the Table class
const (
	COLJUSTIFYLEFT  = 1
	COLJUSTIFYRIGHT = 2

	CELLINT      = 1
	CELLFLOAT    = 2
	CELLSTRING   = 3
	CELLDATE     = 4
	CELLDATETIME = 5

	TABLEOUTTEXT = 1
	TABLEOUTHTML = 2
	TABLEOUTPDF  = 3
	TABLEOUTCSV  = 4

	CSSFONTSIZE = 14
	NEWLINE     = "\n"
)

// Cell is the basic data value type for the Table class
type Cell struct {
	Type int       // int, float, or string enumeration
	Ival int64     // integer value
	Fval float64   // float value
	Sval string    // string value
	Dval time.Time // datetime value
}

// ColumnDef defines a Table column -- a column title, justification, and formatting
// information for cells in the column.
type ColumnDef struct {
	ColTitle  string   // the column title
	Width     int      // column width for TEXT
	Justify   int      // justification
	Pfmt      string   // printf-style formatting information for values in this column
	CellType  int      // type of data in this column
	Hdr       []string // multiple lines of column headers as needed -- based on width and Title
	Fdecimals int      // the number of decimal digits for floating point numbers. The default is 2
	HTMLWidth int
}

// Colset defines a set of Cells
type Colset struct {
	Col    []Cell // 1 row's worth of Cells, contains len(Col) number of Cells
	Height int    // height of row
}

// Rowset defines a set of rows to be operated on at a later time.
type Rowset struct {
	R []int // the row numbers of interest
}

// Table is a structure that defines a spreadsheet-like grid of cells and the
// operations that can be performed.
type Table struct {
	Title           string                             // table title
	Section1        string                             // another section for the title, in a different style
	Section2        string                             // a third section for the title, in a different style
	Section3        string                             // another section for extra usage
	ColDefs         []ColumnDef                        // table's column definitions, ordered 0..n left to right
	Row             []Colset                           // Each Colset forms a row
	maxHdrRows      int                                // maximum number of header rows across all ColDefs
	DateFmt         string                             // format for printing dates
	DateTimeFmt     string                             // format for datetime values
	LineAfter       []int                              // array of row numbers that have a horizontal line after they are printed
	LineBefore      []int                              // array of row numbers that have a horizontal line before they are printed
	RS              []Rowset                           // a list of rowsets
	CSS             map[string]map[string]*CSSProperty //CSS holds css property for title, section1, section2, headers, cells
	htmlTemplate    string                             // path of custom html template path
	htmlTemplateCSS string                             // path of custom css for html template
	// errorList       []string                           // stores the list of error in string format
}

// SetTitle sets the table's Title string to the supplied value.
// Note: Caller is expected to supply "\n" on strings if desired for
// text output.  The "\n" may not be desired or needed for the other
// formats
func (t *Table) SetTitle(s string) {
	t.Title = s
}

// GetTitle sets the table's Title string to the supplied value
func (t *Table) GetTitle() string {
	return t.Title
}

// SetSection1 sets the table's Section1 string to the supplied value
// Note: Caller is expected to supply "\n" on strings if desired for
// text output.  The "\n" may not be desired or needed for the other
// formats
func (t *Table) SetSection1(s string) {
	t.Section1 = s
}

// GetSection1 sets the table's Section1 string to the supplied value
func (t *Table) GetSection1() string {
	return t.Section1
}

// SetSection2 sets the table's Section2 string to the supplied value
// Note: Caller is expected to supply "\n" on strings if desired for
// text output.  The "\n" may not be desired or needed for the other
// formats
func (t *Table) SetSection2(s string) {
	t.Section2 = s
}

// GetSection2 sets the table's Section2 string to the supplied value
func (t *Table) GetSection2() string {
	return t.Section2
}

// SetSection3 sets the table's Section3 string to the supplied value
// Note: Caller is expected to supply "\n" on strings if desired for
// text output.  The "\n" may not be desired or needed for the other
// formats
func (t *Table) SetSection3(s string) {
	t.Section3 = s
}

// GetSection3 sets the table's Section3 string to the supplied value
func (t *Table) GetSection3() string {
	return t.Section3
}

// // AddErrorString adds error string in list of errors
// func (t *Table) AddErrorString(errS string) {
// 	t.errorList = append(t.errorList, errS)
// }

// // RemoveErrorString removes an error from list of errors
// func (t *Table) RemoveErrorString(errS string) {
// 	for i, errStr := range t.errorList {
// 		if errStr == errS {
// 			t.errorList = append(t.errorList[:i], t.errorList[i+1:]...)
// 		}
// 	}
// }

// // GetErrorList returns the list of erros
// func (t *Table) GetErrorList() []string {
// 	return t.errorList
// }

// RowCount returns the number of rows in the table
func (t *Table) RowCount() int {
	return len(t.Row)
}

// ColCount returns the number of columns in the table
func (t *Table) ColCount() int {
	return len(t.ColDefs)
}

// // TypeToString returns a string describing the data type of the cell.
// func (c *Cell) TypeToString() string {
// 	switch c.Type {
// 	case CELLSTRING:
// 		return "string"
// 	case CELLINT:
// 		return "int"
// 	case CELLFLOAT:
// 		return "float"
// 	case CELLDATE:
// 		return "date"
// 	case CELLDATETIME:
// 		return "datetime"
// 	}
// 	return "unknown"
// }

// Init sets internal formatting controls to their default values
func (t *Table) Init() {
	t.DateFmt = "01/02/2006"
	t.DateTimeFmt = "01/02/2006 15:04:00 MST"
	t.CSS = make(map[string]map[string]*CSSProperty)
}

// SetHTMLTemplate sets the path of custom html template
func (t *Table) SetHTMLTemplate(path string) error {
	if ok, _ := isValidFilePath(path); !ok {
		return fmt.Errorf("Provided path %s is not valid", path)
	}

	// set if path is valid
	t.htmlTemplate = path
	return nil
}

// SetHTMLTemplateCSS sets the path of custom css external file for html template
func (t *Table) SetHTMLTemplateCSS(path string) error {
	if ok, _ := isValidFilePath(path); !ok {
		return fmt.Errorf("Provided path %s is not valid", path)
	}

	// set if path is valid
	t.htmlTemplateCSS = path
	return nil
}

// AddLineAfter keeps track of the row numbers after which a line will be printed
func (t *Table) AddLineAfter(row int) {
	t.LineAfter = append(t.LineAfter, row)
	sort.Ints(t.LineAfter)
}

// AddLineBefore keeps track of the row numbers before which a line will be printed
func (t *Table) AddLineBefore(row int) {
	t.LineBefore = append(t.LineBefore, row)
	sort.Ints(t.LineBefore)
}

// CreateRowset creates a new rowset. You can add row indeces to it.  You can process the rows at those indeces later.
// The return value is the Rowset identifier; rsid.  Use it to refer to this rowset.
func (t *Table) CreateRowset() int {
	var a Rowset
	t.RS = append(t.RS, a)
	return len(t.RS) - 1
}

// AppendToRowset adds a new row index to the rowset rsid
func (t *Table) AppendToRowset(rsid, row int) {
	t.RS[rsid].R = append(t.RS[rsid].R, row)
}

// GetRowset returns an array of ints with the rows in the rowset
func (t *Table) GetRowset(rsid int) []int {
	if rsid < 0 || rsid >= len(t.RS) {
		return []int{}
	}
	return t.RS[rsid].R
}

// SumRowset computes the sum of the rows in rowset[rs] at the specified column index. It returns a Cell with the sum
func (t *Table) SumRowset(rsid, col int) Cell {
	var c Cell
	for i := 0; i < len(t.RS[rsid].R); i++ {
		row := t.RS[rsid].R[i]
		switch t.Row[row].Col[col].Type {
		case CELLINT:
			c.Type = CELLINT
			c.Ival += t.Row[row].Col[col].Ival
		case CELLFLOAT:
			c.Type = CELLFLOAT
			c.Fval += t.Row[row].Col[col].Fval
		}
	}
	return c
}

// AdjustFormatString can be called when the format string is null or when the column width changes
// to set a proper formatting string
func (t *Table) AdjustFormatString(cd *ColumnDef) {
	lft := ""
	if cd.Justify == COLJUSTIFYLEFT {
		lft += "-"
	}
	switch cd.CellType {
	case CELLINT:
		cd.Pfmt = fmt.Sprintf("%%%s%dd", lft, cd.Width)
	case CELLFLOAT:
		cd.Pfmt = fmt.Sprintf("%%%d.%ds", cd.Width, cd.Width)
	case CELLSTRING:
		cd.Pfmt = fmt.Sprintf("%%%s%d.%ds", lft, cd.Width, cd.Width)
	}
}

// AddColumn adds a new ColumnDef to the table
func (t *Table) AddColumn(title string, width, celltype int, justification int) {
	var cd = ColumnDef{
		ColTitle: title, Width: width,
		CellType: celltype, Justify: justification,
		Fdecimals: 2, HTMLWidth: -1,
	}
	t.AdjustColumnHeader(&cd)
	t.AdjustFormatString(&cd)
	t.ColDefs = append(t.ColDefs, cd)
}

// AdjustColumnHeader will break up the header into multiple lines if necessary to
// make the title fit.  If necessary, it will force the width of the column to be
// wide enough to fit the longest word in the title.
func (t *Table) AdjustColumnHeader(cd *ColumnDef) {
	a, maxColWidth := getMultiLineText(cd.ColTitle, cd.Width)
	if maxColWidth > cd.Width { // if the length of the column title is greater than the user-specified width
		cd.Width = maxColWidth //increase the column width to hold the column title
	}
	cd.Hdr = a
}

// AdjustAllColumnHeaders formats the column names for printing. It will attempt to break up the column headers
// into multiple lines if necessary.
func (t *Table) AdjustAllColumnHeaders() {
	//----------------------------------
	// Which column has the most rows?
	//----------------------------------
	t.maxHdrRows = 0
	for i := 0; i < len(t.ColDefs); i++ {
		j := len(t.ColDefs[i].Hdr)
		if j > t.maxHdrRows {
			t.maxHdrRows = j
		}
	}

	//---------------------------------------------
	// Set all columns to that number of rows...
	//---------------------------------------------
	for i := 0; i < len(t.ColDefs); i++ {
		n := make([]string, t.maxHdrRows)
		lenOrig := len(t.ColDefs[i].Hdr)
		iStart := t.maxHdrRows - lenOrig
		// Create a new Hdr array, n.
		// Add any initial blank lines...
		if iStart > 0 {
			for j := 0; j < iStart; j++ {
				n[j] = ""
			}
		}
		// now add the remaining strings
		for j := iStart; j < t.maxHdrRows; j++ {
			n[j] = standardizeSpaces(t.ColDefs[i].Hdr[j-iStart])
		}
		t.ColDefs[i].Hdr = n // replace the old hdr with the new one
	}
}

// Get returns the cell at the supplied row,col.  If the supplied
// row or col is outside the table's boundaries, then an empty cell
// is returned
func (t *Table) Get(row, col int) Cell {
	if row >= len(t.Row) || col >= len(t.ColDefs) {
		var c Cell
		return c
	}
	return t.Row[row].Col[col]
}

// Geti returns the int at the supplied row,col.  If the supplied
// row or col is outside the table's boundaries, then 0 is returned
func (t *Table) Geti(row, col int) int64 {
	if row >= len(t.Row) || col >= len(t.ColDefs) {
		return int64(0)
	}
	return t.Row[row].Col[col].Ival
}

// Getf returns the floatval at the supplied row,col.  If the supplied
// row or col is outside the table's boundaries, then 0
// is returned
func (t *Table) Getf(row, col int) float64 {
	if row >= len(t.Row) || col >= len(t.ColDefs) {
		return float64(0)
	}
	return t.Row[row].Col[col].Fval
}

// Gets returns the strinb value at the supplied row,col.  If the supplied
// row or col is outside the table's boundaries, then ""
// is returned
func (t *Table) Gets(row, col int) string {
	if row >= len(t.Row) || col >= len(t.ColDefs) {
		return ""
	}
	return t.Row[row].Col[col].Sval
}

// Getd returns the date at the supplied row,col.  If the supplied
// row or col is outside the table's boundaries, then a 0 date
func (t *Table) Getd(row, col int) time.Time {
	if row >= len(t.Row) || col >= len(t.ColDefs) {
		return time.Date(0, time.January, 0, 0, 0, 0, 0, time.UTC)
	}
	return t.Row[row].Col[col].Dval
}

// Type returns the data type for the cell at the supplied row,col.
// If the supplied row or col is outside the table's boundaries, then 0
// is returned
func (t *Table) Type(row, col int) int {
	if row >= len(t.Row) || col >= len(t.ColDefs) {
		return 0
	}
	return t.Row[row].Col[col].Type
}

// Puti updates the Cell at row,col with the int64 value v
// and sets its type to CELLINT. If row or col is out of
// bounds the return value is false. Otherwise, the return
// value is true
func (t *Table) Puti(row, col int, v int64) bool {
	if row >= len(t.Row) || col >= len(t.ColDefs) {
		return false
	}
	if row < 0 {
		row = len(t.Row) - 1
	}
	t.Row[row].Col[col].Type = CELLINT
	t.Row[row].Col[col].Ival = v
	return true
}

// Putf updates the Cell at row,col with the float64 value v
// and sets its type to CELLFLOAT.
// if row < 0 then row is set to the last row of the table.
// If row or col is out of
// bounds the return value is false. Otherwise, the return
// value is true.
func (t *Table) Putf(row, col int, v float64) bool {
	if row >= len(t.Row) || col >= len(t.ColDefs) {
		return false
	}
	if row < 0 {
		row = len(t.Row) - 1
	}
	t.Row[row].Col[col].Type = CELLFLOAT
	t.Row[row].Col[col].Fval = v
	return true
}

// Puts updates the Cell at row,col with the string value v
// and sets its type to CELLSTRING. If row or col is out of
// bounds the return value is false. Otherwise, the return
// value is true
func (t *Table) Puts(row, col int, v string) bool {
	if row >= len(t.Row) || col >= len(t.ColDefs) {
		return false
	}
	if row < 0 {
		row = len(t.Row) - 1
	}
	t.Row[row].Col[col].Type = CELLSTRING
	t.Row[row].Col[col].Sval = standardizeSpaces(v)

	// Need to check width of column everytime when we adding new content
	// if it is updatable or not
	cd := t.ColDefs[col]
	_, cellWidth := getMultiLineText(v, cd.Width)
	if cellWidth > cd.Width { // if the length of the column title is greater than the user-specified width
		cd.Width = cellWidth //increase the column width to hold the column title
		t.AdjustFormatString(&cd)
		t.ColDefs[col] = cd
	}

	return true
}

// Putd updates the Cell at row,col with the date value v
// and sets its type to CELLDATE. If row or col is out of
// bounds the return value is false. Otherwise, the return
// value is true
func (t *Table) Putd(row, col int, v time.Time) bool {
	return t.putdint(row, col, v, CELLDATE)
}

// Putdt updates the Cell at row,col with the datetimv value v
// and sets its type to CELLDATETIME. If row or col is out of
// bounds the return value is false. Otherwise, the return
// value is true
func (t *Table) Putdt(row, col int, v time.Time) bool {
	return t.putdint(row, col, v, CELLDATETIME)
}

func (t *Table) putdint(row, col int, v time.Time, x int) bool {
	if row >= len(t.Row) || col >= len(t.ColDefs) {
		return false
	}
	if row < 0 {
		row = len(t.Row) - 1
	}
	t.Row[row].Col[col].Type = x
	t.Row[row].Col[col].Dval = v
	return true
}

// Put places Cell c at location row,col
func (t *Table) Put(row, col int, c Cell) {
	if row < 0 {
		row = len(t.Row) - 1
	}
	t.Row[row].Col[col] = c
}

// createColSet creates a new colset with cells, total number of Headers
func (t *Table) createColSet(c *Colset) {
	for i := 0; i < len(t.ColDefs); i++ {
		var cell Cell
		c.Col = append(c.Col, cell)
	}
	c.Height = 1
}

// Sum computes the sum of the rows at the specified column index. It returns a Cell
func (t *Table) Sum(col int) Cell {
	return t.SumRows(col, 0, len(t.Row)-1)
}

// SumRows computes the sum of rows 0 thru row at the specified column index. It returns a Cell
func (t *Table) SumRows(col, from, to int) Cell {
	var c Cell
	if from < 0 {
		from = 0
	}
	if to >= len(t.Row) {
		to = len(t.Row) - 1
	}
	for i := from; i <= to; i++ {
		switch t.Row[i].Col[col].Type {
		case CELLINT:
			c.Type = CELLINT
			c.Ival += t.Row[i].Col[col].Ival
		case CELLFLOAT:
			c.Type = CELLFLOAT
			c.Fval += t.Row[i].Col[col].Fval
		}
	}
	return c
}

// InsertSumRow inserts a new Row at index row, it then sums the specified columns in the Row range: from,to
// and sets the newly inserted row values at the specified columns to the sums.
func (t *Table) InsertSumRow(row, from, to int, cols []int) {
	t.InsertRow(row)
	for i := 0; i < len(cols); i++ {
		c := t.SumRows(cols[i], from, to)
		t.Put(row, cols[i], c)
	}
}

// Sort sorts rows (from,to) by column col ascending
func (t *Table) Sort(from, to, col int) {
	// fmt.Printf("Table.Sort:  from = %d, to = %d, col = %d,  len(t.Row) = %d\n", from, to, col, len(t.Row))
	var swap bool
	for i := from; i < to; i++ {
		for j := i + 1; j <= to; j++ {
			switch t.Row[i].Col[col].Type {
			case CELLINT:
				swap = t.Row[i].Col[col].Ival > t.Row[j].Col[col].Ival
			case CELLFLOAT:
				swap = t.Row[i].Col[col].Fval > t.Row[j].Col[col].Fval
			case CELLSTRING:
				swap = strings.ToLower(t.Row[i].Col[col].Sval) > strings.ToLower(t.Row[j].Col[col].Sval)
			case CELLDATE, CELLDATETIME:
				swap = t.Row[i].Col[col].Dval.After(t.Row[j].Col[col].Dval)
			}
			if swap {
				t.Row[i], t.Row[j] = t.Row[j], t.Row[i]
			}
		}
	}
}

// InsertSumRowsetCols sums the values for the specified rowset and appends it at the specified row
// rsid = the RowSet on which to perform the sum
// row  = a row will be inserted at this index, and the totals will be added to this row
// cols = an array of column numbers to total
func (t *Table) InsertSumRowsetCols(rsid, row int, cols []int) {
	t.InsertRow(row)
	for i := 0; i < len(cols); i++ {
		c := t.SumRowset(rsid, cols[i])
		t.Put(row, cols[i], c)
	}
}

// AddRow appends a new Row to the table. Initially, all cells are empty
// It returns the row number just added.
func (t *Table) AddRow() {
	var c Colset
	t.createColSet(&c)
	t.Row = append(t.Row, c)
}

// InsertRow adds a new Row at the specified index.
func (t *Table) InsertRow(row int) {
	if row >= len(t.Row) || row < 0 {
		t.AddRow()
		return
	}
	var c Colset
	t.createColSet(&c)
	t.Row = append(t.Row[:row+1], t.Row[row:]...)
	t.Row[row] = c

	// Adjust LineAfter
	for i := 0; i < len(t.LineAfter); i++ {
		if t.LineAfter[i] >= row {
			t.LineAfter[i]++
		}
	}
	// Adjust RowSets
	for i := 0; i < len(t.RS); i++ {
		// adjust the existing rows...
		for j := 0; j < len(t.RS[i].R); j++ {
			if t.RS[i].R[j] >= row {
				t.RS[i].R[j]++
			}
		}
		// add in the new row...
		t.RS[i].R = append(t.RS[i].R, row)
	}
}

// DeleteRow removes the table row at the specified index. All rowsets and LineAfter sets are adjusted.
// Cleanup on LineAfter and RowSets does not work if row == 0. I was just too lazy at the time to add this
// code because I know how/where delete will be used and it will not affect row 0.
func (t *Table) DeleteRow(row int) {
	t.Row = t.Row[:row+copy(t.Row[row:], t.Row[row+1:])] // this removes t.Row[row]
	// Clean up LineAfter
	for i := 0; i < len(t.LineAfter); i++ {
		if t.LineAfter[i] > row {
			t.LineAfter[i]--
		}
	}
	// Clean up RowSets
	for i := 0; i < len(t.RS); i++ {
		l := len(t.RS[i].R) // length of jth rowset
		// First remove the row from the rowset
		for j := 0; j < l; j++ {
			if t.RS[i].R[j] == row {
				t.RS[i].R = append(t.RS[i].R[:j], t.RS[i].R[j+1:]...) // this append statement removes element j from the slice
				break
			}
		}
		// Now adjust the row numbers of the remaining rows
		l = len(t.RS[i].R)
		for j := 0; j < l; j++ {
			if t.RS[i].R[j] >= row {
				t.RS[i].R[j]--
			}
		}
	}
}

// TightenColumns goes through all values in STRING columns and determines the maximum length in characters (max).
// If this length is less than the column width the column width is reduced to max.  This is
// mostly useful for text formatting.
func (t *Table) TightenColumns() {
	for i := 0; i < len(t.ColDefs); i++ {
		if t.ColDefs[i].CellType != CELLSTRING {
			continue
		}
		max := 0
		for j := 0; j < len(t.ColDefs[i].Hdr); j++ { // first, find the max len of the col hdrs
			l := len(t.ColDefs[i].Hdr[j])
			if max < l {
				max = l
			}
		}
		for j := 0; j < len(t.Row); j++ { // continue by find the max width of cell values in this col
			if t.Row[j].Col[i].Type == CELLSTRING {
				l := len(t.Row[j].Col[i].Sval)
				if max < l {
					max = l
				}
			}
		}
		if max < t.ColDefs[i].Width { // if the max width is less than the column width, contract the column width
			t.ColDefs[i].Width = max
		}
		cd := t.ColDefs[i]
		t.AdjustFormatString(&cd)
		t.ColDefs[i] = cd
	}
}

// HasData checks that table has actually data or not
func (t *Table) HasData() error {
	// if there are no rows in table
	if t.RowCount() < 1 {
		return fmt.Errorf("No Records Found")
	}
	return nil
}

// HasHeaders checks headers are present or not
func (t *Table) HasHeaders() error {
	if len(t.ColDefs) < 1 {
		return fmt.Errorf("No Header Columns Found")
	}
	return nil
}

// HasValidRow checks that rowIndex is valid or not
func (t *Table) HasValidRow(rowIndex int) error {
	if rowIndex < 0 {
		return fmt.Errorf("Row number is less than zero, row: %d", rowIndex)
	}
	if rowIndex >= t.RowCount() {
		return fmt.Errorf("Row number > no of rows in table, row: %d", rowIndex)
	}
	return nil
}

// HasValidColumn checks that colIndex is valid or not
func (t *Table) HasValidColumn(colIndex int) error {
	if colIndex < 0 {
		return fmt.Errorf("Column number is less than zero, column: %d", colIndex)
	}
	if colIndex >= t.ColCount() {
		return fmt.Errorf("Column number > no of columns in table, column: %d", colIndex)
	}
	return nil
}

// ==========================
// Table Export Output
// ==========================

// TableExportType each export output format must satisfy this interface
type TableExportType interface {
	writeTableOutput(w io.Writer) error
	formatTitle() string
	formatSection1() string
	formatSection2() string
	formatSection3() string
	// getErrorSection() string
	formatHeaders() (string, error)
	formatRows() (string, error)
	formatRow(row int) (string, error)
}

// String is the "stringer" method implementation for gotable so that you can simply
// print(t)
func (t Table) String() string {
	var temp bytes.Buffer
	err := t.FprintTable(&temp)
	if err != nil {
		return err.Error()
	}
	return temp.String()
}

// SprintTable renders the entire table to a string for text output
func (t *Table) SprintTable() (string, error) {
	var temp bytes.Buffer
	err := t.FprintTable(&temp)
	if err != nil {
		return "", err
	}
	return temp.String(), nil
}

// FprintTable renders the entire table for io.Writer object for text output
func (t *Table) FprintTable(w io.Writer) error {
	var tout TableExportType = &TextTable{Table: t, TextColSpace: 2}
	return tout.writeTableOutput(w)
}

// TextprintTable renders the entire table for text output, alias for FprintTable
func (t *Table) TextprintTable(w io.Writer) error {
	return t.FprintTable(w)
}

// CSVprintTable renders the entire table for csv output
func (t *Table) CSVprintTable(w io.Writer) error {
	var tout = &CSVTable{Table: t}
	return tout.writeTableOutput(w)
}

// HTMLprintTable renders the entire table for html output
func (t *Table) HTMLprintTable(w io.Writer) error {
	var tout TableExportType = &HTMLTable{Table: t}
	return tout.writeTableOutput(w)
}

// PDFprintTable renders the entire table for pdf output
func (t *Table) PDFprintTable(w io.Writer) error {
	var tout = &PDFTable{Table: t}
	return tout.writeTableOutput(w)
}

// ==========================
// METHODs for HTML output //
// ==========================

// CSSProperty holds css property to be used as inline css
type CSSProperty struct {
	Name, Value string
}

// String is the "stringer" method implementation for CSSProperty
func (cp CSSProperty) String() string {
	return `"` + cp.Name + `:` + cp.Value + `;"`
}

// getCSSMapKeyForCell format and returns key for cell for css properties usage
func (t *Table) getCSSMapKeyForCell(rowIndex, colIndex int) string {
	return `row:` + strconv.Itoa(rowIndex) + `-col:` + strconv.Itoa(colIndex)
}

// getCSSMapKeyForHeaderCell format and returns key for eader cell for css properties usage
func (t *Table) getCSSMapKeyForHeaderCell(colIndex int) string {
	return `header-` + strconv.Itoa(colIndex)
}

// SetTitleCSS sets css for title row
func (t *Table) SetTitleCSS(cssList []*CSSProperty) {
	// css property map
	cssMap, ok := t.CSS[TITLECLASS]
	if !ok {
		cssMap = make(map[string]*CSSProperty)
	}

	// map it in style of html table
	for _, cssProp := range cssList {
		cssMap[cssProp.Name] = cssProp
	}

	t.CSS[TITLECLASS] = cssMap
}

// SetSection1CSS sets css for section1 row
func (t *Table) SetSection1CSS(cssList []*CSSProperty) {
	// css property map
	cssMap, ok := t.CSS[SECTION1CLASS]
	if !ok {
		cssMap = make(map[string]*CSSProperty)
	}

	// map it in style of html table
	for _, cssProp := range cssList {
		cssMap[cssProp.Name] = cssProp
	}

	t.CSS[SECTION1CLASS] = cssMap
}

// SetSection2CSS sets css for section2 row
func (t *Table) SetSection2CSS(cssList []*CSSProperty) {
	// css property map
	cssMap, ok := t.CSS[SECTION2CLASS]
	if !ok {
		cssMap = make(map[string]*CSSProperty)
	}

	// map it in style of html table
	for _, cssProp := range cssList {
		cssMap[cssProp.Name] = cssProp
	}

	t.CSS[SECTION2CLASS] = cssMap
}

// SetSection3CSS sets css for section3 row
func (t *Table) SetSection3CSS(cssList []*CSSProperty) {
	// css property map
	cssMap, ok := t.CSS[SECTION3CLASS]
	if !ok {
		cssMap = make(map[string]*CSSProperty)
	}

	// map it in style of html table
	for _, cssProp := range cssList {
		cssMap[cssProp.Name] = cssProp
	}

	t.CSS[SECTION3CLASS] = cssMap
}

// // SetErrorSectionCSS sets css for error section row
// func (t *Table) SetErrorSectionCSS(cssList []*CSSProperty) {
// 	// css property map
// 	cssMap, ok := t.CSS[ERRORSSECTION]
// 	if !ok {
// 		cssMap = make(map[string]*CSSProperty)
// 	}

// 	// map it in style of html table
// 	for _, cssProp := range cssList {
// 		cssMap[cssProp.Name] = cssProp
// 	}

// 	t.CSS[ERRORSSECTION] = cssMap
// }

// SetHeaderCellCSS sets css for only headers cell
func (t *Table) SetHeaderCellCSS(colIndex int, cssList []*CSSProperty) error {
	// check row is valid or not
	if err := t.HasValidColumn(colIndex); err != nil {
		return err
	}

	// header class
	thClass := t.getCSSMapKeyForHeaderCell(colIndex)
	// css property map
	cssMap, ok := t.CSS[thClass]
	if !ok {
		cssMap = make(map[string]*CSSProperty)
	}

	// map it in style of html table
	for _, cssProp := range cssList {
		cssMap[cssProp.Name] = cssProp
	}

	t.CSS[thClass] = cssMap

	return nil
}

// SetHeaderCSS sets css for headers row
func (t *Table) SetHeaderCSS(cssList []*CSSProperty) {
	for colIndex := 0; colIndex < t.ColCount(); colIndex++ {
		t.SetHeaderCellCSS(colIndex, cssList)
	}
}

// SetColHTMLWidth sets the column width for table
func (t *Table) SetColHTMLWidth(colIndex int, width uint, unit string) error {

	// fix the bug of unit coversion
	// unit will be not in effect, only one unit for all cells in table will be applied

	// TODO: conversion from different units of font to `ch` unit with body font base size
	// so that width has value with `px` unit value

	if err := t.HasValidColumn(colIndex); err != nil {
		return err
	}

	t.ColDefs[colIndex].HTMLWidth = int(width)
	return nil
}

// SetCellCSS sets css properties for Table Cells
func (t *Table) SetCellCSS(rowIndex, colIndex int, cssList []*CSSProperty) error {

	// check row is valid or not
	if err := t.HasValidRow(rowIndex); err != nil {
		return err
	}

	// check row is valid or not
	if err := t.HasValidColumn(colIndex); err != nil {
		return err
	}

	// css property map
	g := t.getCSSMapKeyForCell(rowIndex, colIndex)
	cssMap, ok := t.CSS[g]
	if !ok {
		cssMap = make(map[string]*CSSProperty)
	}

	// map it in style of html table
	for _, cssProp := range cssList {
		cssMap[cssProp.Name] = cssProp
	}

	t.CSS[g] = cssMap

	return nil
}

// SetAllCellCSS sets css properties for all Table Cells
func (t *Table) SetAllCellCSS(cssList []*CSSProperty) {

	// convert it into cells attributes
	for colIndex := 0; colIndex < t.ColCount(); colIndex++ {
		for rowIndex := 0; rowIndex < t.RowCount(); rowIndex++ {
			// will never meet an error from below function
			t.SetCellCSS(rowIndex, colIndex, cssList)
		}
	}
}

// SetRowCSS sets css properties for Table Rows
func (t *Table) SetRowCSS(rowIndex int, cssList []*CSSProperty) error {

	// check row is valid or not
	if err := t.HasValidRow(rowIndex); err != nil {
		return err
	}

	// convert it into cells attributes
	for colIndex := 0; colIndex < t.ColCount(); colIndex++ {
		// for valid rowIndex set css for all cells belongs to rowIndex row
		t.SetCellCSS(rowIndex, colIndex, cssList)
	}

	return nil
}

// SetColCSS sets css properties for Table Columns
func (t *Table) SetColCSS(colIndex int, cssList []*CSSProperty) error {

	// check row is valid or not
	if err := t.HasValidColumn(colIndex); err != nil {
		return err
	}

	// convert it into cells attributes
	for rowIndex := 0; rowIndex < t.RowCount(); rowIndex++ {
		// for valid colIndex set css for all cells belongs to colIndex column
		t.SetCellCSS(rowIndex, colIndex, cssList)
	}

	return nil
}

// SetNoRowsCSS sets css on html tag which has no rows msg
func (t *Table) SetNoRowsCSS(cssList []*CSSProperty) {
	// css property map
	cssMap, ok := t.CSS[NOROWSCLASS]
	if !ok {
		cssMap = make(map[string]*CSSProperty)
	}

	// map it in style of html table
	for _, cssProp := range cssList {
		cssMap[cssProp.Name] = cssProp
	}

	t.CSS[NOROWSCLASS] = cssMap
}

// SetNoHeadersCSS sets css on html tag which has no headers/columns msg
func (t *Table) SetNoHeadersCSS(cssList []*CSSProperty) {
	// css property map
	cssMap, ok := t.CSS[NOHEADERSCLASS]
	if !ok {
		cssMap = make(map[string]*CSSProperty)
	}

	// map it in style of html table
	for _, cssProp := range cssList {
		cssMap[cssProp.Name] = cssProp
	}

	t.CSS[NOHEADERSCLASS] = cssMap
}
