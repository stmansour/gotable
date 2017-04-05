package gotable

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

// WKHTMLTOPDFCMD command : html > pdf
const (
	WKHTMLTOPDFCMD = "wkhtmltopdf"
	TEMPSTORE      = "."
	DATETIMEFMT    = "_2 Jan 2006 3:04 PM IST"
)

// PDFTable struct used to prepare table in pdf version
type PDFTable struct {
	*Table
	buf bytes.Buffer
}

func (pt *PDFTable) writeTableOutput(w io.Writer) error {

	// get html output first
	var temp bytes.Buffer

	// copy table object so that we can override properties over table
	// so it won't affect original table
	var pdfTable = *pt.Table
	var ht = &HTMLTable{Table: &pdfTable}

	// set custom values over ht
	ht.SetCSSFontUnit("px")

	var tout TableExportType = ht

	if err := tout.writeTableOutput(&temp); err != nil {
		return err
	}

	htmlString := temp.String()

	timeCharReplacer := strings.NewReplacer(":", "-", ".", "", "T", "")
	currentTime := timeCharReplacer.Replace(time.Now().Format(time.RFC3339Nano))

	// create temp file
	filePath := path.Join(TEMPSTORE, "tablePDF_"+currentTime)

	// only works with html file extension
	// be careful, must append it
	tempHTMLFile, err := os.Create(filePath + ".html")
	if err != nil {
		return err
	}
	// write html string to file
	tempHTMLFile.WriteString(htmlString)
	tempHTMLFile.Close()

	// remove this temp file after operation
	defer os.Remove(tempHTMLFile.Name())

	// return output file path
	if err = pt.writePDFBuffer(filePath); err != nil {
		return err
	}

	// write output to passed io.Writer interface object
	w.Write(pt.buf.Bytes())
	return err
}

func (pt *PDFTable) writePDFBuffer(inputFile string) error {

	pdfExportTime := time.Now().Format(DATETIMEFMT)
	htmlExportFile := inputFile + ".html"

	cmdArgs := []string{
		"--disable-smart-shrinking",
		// top margin
		"-T", "15",
		// header center content
		"--header-center", pt.Table.GetTitle(),
		// header font size
		"--header-font-size", "7",
		// header font
		"--header-font-name", "opensans",
		// header spacing
		"--header-spacing", "3",
		// bottom margin
		"-B", "15",
		// footer spacing
		"--footer-spacing", "5",
		// footer font
		"--footer-font-name", "opensans",
		// footer font size
		"--footer-font-size", "7",
		// footer left content
		"--footer-left", pdfExportTime,
		// footer right content
		"--footer-right", "Page [page] of [toPage]",
		// page size
		"--page-size", "Letter",
		// orientation
		"--orientation", "Landscape",
		// input, output
		htmlExportFile, "-",
	}

	// prepare command
	wkhtmltopdf := exec.Command(WKHTMLTOPDFCMD, cmdArgs...)

	// REF: https://github.com/aodin/go-pdf-server/blob/master/pdf_server.go

	// get output pipeline
	output, err := wkhtmltopdf.StdoutPipe()
	if err != nil {
		return err
	}

	// Begin the command
	if err = wkhtmltopdf.Start(); err != nil {
		return err
	}

	// Read the generated PDF from std out
	b, err := ioutil.ReadAll(output)
	if err != nil {
		return err
	}

	// End the command
	if err = wkhtmltopdf.Wait(); err != nil {
		return err
	}

	// write output to buffer
	pt.buf.Write(b)

	return nil
}
