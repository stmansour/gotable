package gotable

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"sort"
	"strconv"
	"text/template"

	"github.com/dustin/go-humanize"
	"github.com/kardianos/osext"
	"github.com/yosssi/gohtml"
)

// TABLECONTAINERCLASS et. al. are the constants used in the html version of table object
const (
	TABLECONTAINERCLASS = `rpt-table-container`
	TITLECLASS          = `title`
	SECTION1CLASS       = `section1`
	SECTION2CLASS       = `section2`
	SECTION3CLASS       = `section3`
	ERRORSSECTION       = `error-section`

	NOROWSCLASS    = `no-rows`
	NOHEADERSCLASS = `no-headers`

	// HEADERSCLASS        = `headers`
	// DATACLASS           = `data`
)

// HTMLTable struct used to prepare table in html version
type HTMLTable struct {
	*Table
	styleString bytes.Buffer
	buf         bytes.Buffer
}

// HTMLTemplateContext holds the context for table html template
type HTMLTemplateContext struct {
	FontSize                                    int
	HeadTitle, DefaultCSS, CustomCSS, TableHTML string
}

func (ht *HTMLTable) writeTableOutput(w io.Writer) error {

	// vars
	var (
		err error
	)

	// append title
	ht.buf.WriteString(ht.formatTitle())

	// append section 1
	ht.buf.WriteString(ht.formatSection1())

	// append section 2
	ht.buf.WriteString(ht.formatSection2())

	// append section 3
	ht.buf.WriteString(ht.formatSection3())

	// append headers
	var tableHdrsRows bytes.Buffer
	if headerStr, err := ht.formatHeaders(); err != nil {
		if cellCSSProps, ok := ht.getCSSPropertyList(NOHEADERSCLASS); ok {
			// get css string for section1
			ht.styleString.WriteString(`div.` + TABLECONTAINERCLASS + ` p`)
			ht.styleString.WriteString(ht.getCSSForClassSelector(NOHEADERSCLASS, cellCSSProps))
		}
		tableHdrsRows.WriteString(`<p class="` + NOHEADERSCLASS + `">` + err.Error() + `</p>`)
	} else {
		// if headers found then append rows
		if rowsStr, err := ht.formatRows(); err != nil {
			colSpan := strconv.Itoa(ht.Table.ColCount())
			if cellCSSProps, ok := ht.getCSSPropertyList(NOROWSCLASS); ok {
				// get css string for section1
				ht.styleString.WriteString(`div.` + TABLECONTAINERCLASS + ` table tbody tr td`)
				ht.styleString.WriteString(ht.getCSSForClassSelector(NOROWSCLASS, cellCSSProps))
			}
			noRowsTD := `<td colspan="` + colSpan + `" class="` + NOROWSCLASS + `">` + err.Error() + `</td>`
			tableHdrsRows.WriteString(`<tbody><tr>` + noRowsTD + `</tr></tbody>`)
		} else {
			// if rows exist, then only show headers
			tableHdrsRows.WriteString(headerStr)
			tableHdrsRows.WriteString(rowsStr)
		}

		// wrap headers and rows in a table
		tmpTableHdrsRows := tableHdrsRows.String()
		tableHdrsRows.Reset()
		tableHdrsRows.WriteString(`<table>` + tmpTableHdrsRows + `</table>`)
	}

	// // render errorlist
	// NOTE: if you enable this errorList feature then write them first on top,
	// then write headers, rows output
	// ct.buf.Write(ct.getErrorSection())

	// if it has some content then write it in buf
	if tableHdrsRows.Len() > 0 {
		ht.buf.Write(tableHdrsRows.Bytes())
	}

	// wrap it up in a div with a class
	tmpHTML := ht.buf.String()
	ht.buf.Reset()
	ht.buf.WriteString(`<div class="` + TABLECONTAINERCLASS + `">` + tmpHTML + `</div>`)

	// format and store html output in ht buf
	if err = ht.formatHTML(); err != nil {
		errorLog("formatHTML err: ", err.Error())
		return err
	}

	// after formatted output stored in ht.buf, write it to w
	w.Write(ht.buf.Bytes())
	return err
}

func (ht *HTMLTable) formatHTML() error {
	var err error

	// make context for template
	htmlContext := HTMLTemplateContext{FontSize: CSSFONTSIZE}
	htmlContext.HeadTitle = ht.Table.Title

	htmlContext.DefaultCSS, err = ht.getTableCSS()
	if err != nil {
		errorLog("While getting table css: ", err.Error())
		return err
	}
	htmlContext.DefaultCSS = `<style>` + htmlContext.DefaultCSS + `</style>`
	htmlContext.CustomCSS = `<style>` + ht.styleString.String() + `</style>`
	htmlContext.TableHTML = ht.buf.String()

	// get template string
	tmpl, err := ht.getHTMLTemplate()
	if err != nil {
		errorLog("While getting table template: ", err.Error())
		return err
	}

	// write html output in buffer
	// before writing it to buf, reset it
	ht.buf.Reset()
	err = tmpl.Execute(&ht.buf, htmlContext)
	if err != nil {
		errorLog("While template execution: ", err.Error())
		return err
	}

	// write buffered output after formatting html
	tmpHTMLString := ht.buf.String()
	ht.buf.Reset()
	// beautify html output, it is nice to have, not necessarY
	ht.buf.WriteString(gohtml.Format(tmpHTMLString))
	infoLog("HTML output has been formatted and written to internal buffer. :)")

	return nil
}

func (ht *HTMLTable) formatTitle() string {
	title := ht.Table.GetTitle()

	if title != "" {
		if cellCSSProps, ok := ht.getCSSPropertyList(TITLECLASS); ok {
			// get css string for title
			ht.styleString.WriteString(`div.` + TABLECONTAINERCLASS + ` p`)
			ht.styleString.WriteString(ht.getCSSForClassSelector(TITLECLASS, cellCSSProps))
		}
		return `<p class="` + TITLECLASS + `">` + title + `</p>`
	}

	// blank return
	return title
}

func (ht *HTMLTable) formatSection1() string {
	section1 := ht.Table.GetSection1()

	if section1 != "" {
		if cellCSSProps, ok := ht.getCSSPropertyList(SECTION1CLASS); ok {
			// get css string for section1
			ht.styleString.WriteString(`div.` + TABLECONTAINERCLASS + ` p`)
			ht.styleString.WriteString(ht.getCSSForClassSelector(SECTION1CLASS, cellCSSProps))
		}
		return `<p class="` + SECTION1CLASS + `">` + section1 + `</p>`
	}

	// blank return
	return section1
}

func (ht *HTMLTable) formatSection2() string {
	section2 := ht.Table.GetSection2()

	if section2 != "" {
		if cellCSSProps, ok := ht.getCSSPropertyList(SECTION2CLASS); ok {
			// get css string for section2
			ht.styleString.WriteString(`div.` + TABLECONTAINERCLASS + ` p`)
			ht.styleString.WriteString(ht.getCSSForClassSelector(SECTION2CLASS, cellCSSProps))
		}
		return `<p class="` + SECTION2CLASS + `">` + section2 + `</p>`
	}

	// blank return
	return section2
}

func (ht *HTMLTable) formatSection3() string {
	section3 := ht.Table.GetSection3()

	if section3 != "" {
		if cellCSSProps, ok := ht.getCSSPropertyList(SECTION3CLASS); ok {
			// get css string for section3
			ht.styleString.WriteString(`div.` + TABLECONTAINERCLASS + ` p`)
			ht.styleString.WriteString(ht.getCSSForClassSelector(SECTION3CLASS, cellCSSProps))
		}
		return `<p class="` + SECTION3CLASS + `">` + section3 + `</p>`
	}

	// blank return
	return section3
}

// func (ht *HTMLTable) getErrorSection() string {
// 	errSection := ""

// 	errList := ht.Table.GetErrorList()
// 	if len(errList) > 0 {
// 		for i, errStr := range errList {
// 			index := strconv.Itoa(i)
// 			errSection += `<p class="error-` + index + `">` + errStr + `</p>`
// 		}
// 		return `<div class="` + ERRORSSECTION + `">` + errSection + `</div>`
// 	}

// 	// blank return
// 	return errSection
// }

func (ht *HTMLTable) formatHeaders() (string, error) {

	// check for blank headers
	blankHdrsErr := ht.Table.HasHeaders()
	if blankHdrsErr != nil {
		return "", blankHdrsErr
	}

	// format headers
	var tHeaders bytes.Buffer

	for headerIndex := 0; headerIndex < len(ht.Table.ColDefs); headerIndex++ {

		headerCell := ht.Table.ColDefs[headerIndex]

		// css class for this header cell
		thClass := ht.Table.getCSSMapKeyForHeaderCell(headerIndex)

		// --------------------
		// Text Alignment
		// --------------------
		// decide align property
		alignProp := &CSSProperty{Name: "text-align"}
		if headerCell.Justify == COLJUSTIFYRIGHT {
			alignProp.Value = "right"
		} else if headerCell.Justify == COLJUSTIFYLEFT {
			alignProp.Value = "left"
		}

		// set align css for header cell
		ht.Table.SetHeaderCellCSS(headerIndex, []*CSSProperty{alignProp})

		// apply this property to all cells belong to this column
		ht.Table.SetColCSS(headerIndex, []*CSSProperty{alignProp})

		// --------------------
		// Column width
		// --------------------
		// NOTE: width calculatation should be done after alignment
		// width only needs to be set on header cells only not on all
		// cells belong to column
		var colWidthUnit string
		var colWidth int

		if headerCell.HTMLWidth != -1 {
			// calculate column width based on characters with font size
			colWidth = headerCell.HTMLWidth
		} else {
			// calculate column width based on characters with font size
			colWidth = ht.Table.ColDefs[headerIndex].Width
		}

		// if fontUnit is px then need to convert width in px
		switch ht.Table.fontUnit {
		case "px":
			colWidth = colWidth * CSSFONTSIZE
		}
		// TODO: put other units conversion switch cases too.....
		colWidthUnit = strconv.Itoa(colWidth) + ht.Table.fontUnit

		// set width css property on this header cell, no need to apply on each and every cell of this column
		ht.Table.SetHeaderCellCSS(headerIndex, []*CSSProperty{{Name: "width", Value: colWidthUnit}})

		// --------------------
		// apply css on each header cell
		// --------------------
		// get css props for this header cell in SORTED manner
		cellCSSProps, _ := ht.getCSSPropertyList(thClass)

		// get css string for headers
		ht.styleString.WriteString(`div.` + TABLECONTAINERCLASS + ` table thead tr th`)

		// ht.styleString.WriteString(`div.` + TABLECONTAINERCLASS + ` table thead.` + HEADERSCLASS + ` tr th`)
		ht.styleString.WriteString(ht.getCSSForClassSelector(thClass, cellCSSProps))

		// append each header cells in tHeaders
		tHeaders.WriteString(`<th class="` + thClass + `">` + headerCell.ColTitle + `</th>`)
	}

	return `<thead><tr>` + tHeaders.String() + `</tr></thead>`, nil
	// return `<thead class="` + HEADERSCLASS + `"><tr>` + tHeaders.WriteString() + `</tr></thead>`, nil
}

func (ht *HTMLTable) formatRows() (string, error) {

	// check for empty data table
	blankDataErr := ht.Table.HasData()
	if blankDataErr != nil {
		return "", blankDataErr
	}

	var rowsOut bytes.Buffer
	for i := 0; i < ht.Table.RowCount(); i++ {
		// for valid row, we will never get an error
		s, _ := ht.formatRow(i)
		rowsOut.WriteString(s)
	}

	// return with wrapping in tag tbody
	return `<tbody>` + rowsOut.String() + `</tbody>`, nil
}

func (ht *HTMLTable) formatRow(rowIndex int) (string, error) {

	// This method is only called by internal instance of TextTable
	// in formatRows method, so we should avoid following error check
	// unless we make it as export

	// // check that this passed rowIndex is valid or not
	// inValidRowErr := ht.Table.HasValidRow(rowIndex)
	// if inValidRowErr != nil {
	// 	return "", inValidRowErr
	// }

	// format table rows
	var tRow bytes.Buffer
	var trClass string

	if len(ht.Table.LineBefore) > 0 {
		j := sort.SearchInts(ht.Table.LineBefore, rowIndex)
		// line separator added in `LineAfter`??
		// If YES, then discard it
		sepExist := sort.SearchInts(ht.Table.LineAfter, rowIndex-1) < ht.Table.RowCount()
		if j < len(ht.Table.LineBefore) && rowIndex == ht.Table.LineBefore[j] && !sepExist {
			trClass += `top-line`
		}
	}

	// fill the content in rowTextList for the first line
	for colIndex := 0; colIndex < len(ht.Table.Row[rowIndex].Col); colIndex++ {

		var rowCell string
		// append content in TD
		switch ht.Table.Row[rowIndex].Col[colIndex].Type {
		case CELLFLOAT:
			rowCell = fmt.Sprintf(ht.Table.ColDefs[colIndex].Pfmt, humanize.FormatFloat("#,###.##", ht.Table.Row[rowIndex].Col[colIndex].Fval))
		case CELLINT:
			rowCell = fmt.Sprintf(ht.Table.ColDefs[colIndex].Pfmt, ht.Table.Row[rowIndex].Col[colIndex].Ival)
		case CELLSTRING:
			// ******************************************************
			// FOR HTML, APPEND FULL STRING, THERE ARE NO
			// MULTILINE TEXT IN THIS
			// ******************************************************
			rowCell = fmt.Sprintf("%s", ht.Table.Row[rowIndex].Col[colIndex].Sval)
		case CELLDATE:
			rowCell = fmt.Sprintf("%*.*s", ht.Table.ColDefs[colIndex].Width, ht.Table.ColDefs[colIndex].Width, ht.Table.Row[rowIndex].Col[colIndex].Dval.Format(ht.Table.DateFmt))
		case CELLDATETIME:
			rowCell = fmt.Sprintf("%*.*s", ht.Table.ColDefs[colIndex].Width, ht.Table.ColDefs[colIndex].Width, ht.Table.Row[rowIndex].Col[colIndex].Dval.Format(ht.Table.DateTimeFmt))
		default:
			rowCell = mkstr(ht.Table.ColDefs[colIndex].Width, ' ')
		}

		// format td cell with custom class if exists for it
		g := ht.Table.getCSSMapKeyForCell(rowIndex, colIndex)
		if cellCSSProps, ok := ht.getCSSPropertyList(g); ok {

			tdClass := `cell-row-` + strconv.Itoa(rowIndex) + `-col-` + strconv.Itoa(colIndex)

			// get css string for a row
			ht.styleString.WriteString(`div.` + TABLECONTAINERCLASS + ` table tbody tr td`)
			ht.styleString.WriteString(ht.getCSSForClassSelector(tdClass, cellCSSProps))

			tRow.WriteString(`<td class="` + tdClass + `">` + rowCell + `</td>`)
		} else {
			tRow.WriteString(`<td>` + rowCell + `</td>`)
		}
	}

	if len(ht.Table.LineAfter) > 0 {
		j := sort.SearchInts(ht.Table.LineAfter, rowIndex)
		if j < len(ht.Table.LineAfter) && rowIndex == ht.Table.LineAfter[j] {
			trClass += `bottom-line`
		}
	}

	if trClass != "" {
		return `<tr class="` + trClass + `">` + tRow.String() + `</tr>`, nil
	}
	return `<tr>` + tRow.String() + `</tr>`, nil
}

// getCSSForClassSelector returns css string for a class
func (ht *HTMLTable) getCSSForClassSelector(className string, cssList []*CSSProperty) string {
	var classCSS string

	// append notation for selector
	classCSS += `.` + className + `{`

	for _, cssProp := range cssList {
		// append css property name
		classCSS += cssProp.Name + `:` + cssProp.Value + `;`
	}

	// finally block ending sign
	classCSS += `}`

	// return class css string
	return classCSS
}

// getCSSForHTMLTag return css string for html tag element
// func (ht *HTMLTable) getCSSForHTMLTag(tagEl string, cssList []*CSSProperty) string {
// 	var classCSS string

// 	// append notation for selector
// 	classCSS += tagEl + `{`

// 	for _, cssProp := range cssList {
// 		// append css property name
// 		classCSS += cssProp.Name + `:` + cssProp.Value + `;`
// 	}

// 	// finally block ending sign
// 	classCSS += `}`

// 	// return class css string
// 	return classCSS
// }

// getTableCSS reads default css and return the content of it
func (ht *HTMLTable) getTableCSS() (string, error) {
	funcname := "getTableCSS"

	// 1. Get the content from custom css file if it exist
	cssPath := ht.Table.htmlTemplateCSS
	if ok, _ := isValidFilePath(cssPath); ok {
		cssString, err := ioutil.ReadFile(cssPath)
		if err != nil {
			errorLog(funcname, ": ", err.Error())
			return "", err
		}
		return string(cssString), nil
	}

	// 2. Get the content from default file, from within the execution path
	// in case first trial failed
	exDirPath, err := osext.ExecutableFolder()
	if err != nil {
		errorLog(funcname, ": ", err.Error())
		return "", err
	}
	cssPath = path.Join(exDirPath, "gotable.css")
	if ok, _ := isValidFilePath(cssPath); ok {
		cssString, err := ioutil.ReadFile(cssPath)
		if err != nil {
			errorLog(funcname, ": ", err.Error())
			return "", err
		}
		return string(cssString), nil
	}

	// 3. Get content from constant value defined in defaults.go
	// in case second trial failed
	return DCSS, nil
}

// getHTMLTemplate returns the *Template object, error
func (ht *HTMLTable) getHTMLTemplate() (*template.Template, error) {
	funcname := "getHTMLTemplate"

	// 1. Get the content from custom template file if it exist
	tmplPath := ht.Table.htmlTemplate
	if ok, _ := isValidFilePath(tmplPath); ok {

		// generates new template and parse content from html and returns it
		if tmpl, err := template.ParseFiles(tmplPath); err != nil {
			goto tmplexdir2
		} else {
			// if no error then return simply
			return tmpl, err
		}
	}

	// 2. Get the content from default file, from within the execution path
	// in case first trial failed
tmplexdir2:
	exDirPath, err := osext.ExecutableFolder()
	if err != nil {
		errorLog(funcname, ": ", err.Error())
		return nil, err
	}

	tmplPath = path.Join(exDirPath, "./tmpl/gotable.tmpl")
	if ok, _ := isValidFilePath(tmplPath); ok {

		// generates new template and parse content from html and returns it
		if tmpl, err := template.ParseFiles(tmplPath); err != nil {
			goto tmplconst3
		} else {
			// if no error then return simply
			return tmpl, err
		}
	}
tmplconst3:
	// 3. Get content from constant value defined in defaults.go
	// in case second trial failed
	tmpl, err := template.New("gotable.tmpl").Parse(DTEMPLATE)

	// finally return *Template, Error
	if err != nil {
		errorLog(funcname, ": ", err.Error())
	}
	return tmpl, err
}

// getCSSPropertyList returns the css property list from css map of table object
func (ht *HTMLTable) getCSSPropertyList(element string) ([]*CSSProperty, bool) {

	var ok bool
	var cellCSSProps []*CSSProperty

	if cssMap, ok := ht.Table.CSS[element]; ok {

		// sort list of css by its name
		cssNameList := []string{}
		for cssName := range cssMap {
			cssNameList = append(cssNameList, cssName)
		}
		sort.Strings(cssNameList)

		// list of css properties for this td cell
		for _, cssName := range cssNameList {
			cellCSSProps = append(cellCSSProps, cssMap[cssName])
		}

		// return
		return cellCSSProps, ok
	}

	// return
	return cellCSSProps, ok
}

// MultiTableHTMLPrint writes html output from each table to w io.Writer
func MultiTableHTMLPrint(m []Table, w io.Writer) error {
	funcname := "MultiTableHTMLPrint"

	for i := 0; i < len(m); i++ {

		// set custom template for reports
		if i == 0 {
			// set first table layout template
			m[i].SetHTMLTemplate("./tmpl/firstTable.tmpl")
		} else if i == len(m)-1 {
			// set last table layout template
			m[i].SetHTMLTemplate("./tmpl/lastTable.tmpl")
		} else {
			// set middle table layout template
			m[i].SetHTMLTemplate("./tmpl/middleTable.tmpl")
		}

		temp := bytes.Buffer{}
		err := m[i].HTMLprintTable(&temp)
		if err != nil {
			errorLog("%s: Error while getting table output, title: %s, err: %s", funcname, m[i].Title, err.Error())
			return err
		}
		w.Write(temp.Bytes())
	}

	return nil
}
