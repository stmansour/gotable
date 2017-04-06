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

// PDFProperty struct used to hold wkhtmltopdf pdf properties
type PDFProperty struct {
	Option, Value string // Value could be optional
}

func (pt *PDFTable) writeTableOutput(w io.Writer, pdfProps []*PDFProperty) error {

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
		errorLog("Unable to write html output of table to buffer: ", err.Error())
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
		errorLog("Unable to create temporary html file for wkhtmltopdf stdin: ", err.Error())
		return err
	}
	// write html string to file
	tempHTMLFile.WriteString(htmlString)
	tempHTMLFile.Close()

	// remove this temp file after operation
	defer os.Remove(tempHTMLFile.Name())

	// return output file path
	if err = pt.writePDFBuffer(filePath, pdfProps); err != nil {
		errorLog("writePDFBuffer error : ", err.Error())
		return err
	}

	// write output to passed io.Writer interface object
	w.Write(pt.buf.Bytes())
	return err
}

func (pt *PDFTable) writePDFBuffer(inputFile string, pdfProps []*PDFProperty) error {

	htmlExportFile := inputFile + ".html"

	// pdfOpts holds only options which does not require any value
	pdfOpts := []string{}

	// pdfOptsV holds options which has a value
	pdfOptsV := []string{}

	for _, prop := range pdfProps {
		if prop.Option != "" && prop.Value != "" {
			// option which has value
			pdfOptsV = append(pdfOptsV, prop.Option)
			pdfOptsV = append(pdfOptsV, prop.Value)
		} else if prop.Option != "" && prop.Value == "" {
			// option which does not require value
			pdfOpts = append(pdfOpts, prop.Option)
		}
	}

	// make cmdArgs list from pdfOpts and pdfOptsV
	cmdArgs := []string{}
	// first append option which has no option
	for _, opt := range pdfOpts {
		cmdArgs = append(cmdArgs, opt)
	}
	// later append option with value
	for _, optV := range pdfOptsV {
		cmdArgs = append(cmdArgs, optV)
	}

	// append input and output finally
	cmdArgs = append(cmdArgs, []string{htmlExportFile, "-"}...)
	debugLog("Command line arguments for wkhtmltopdf:\n", cmdArgs, "\n\n")

	// prepare command
	wkhtmltopdf := exec.Command(WKHTMLTOPDFCMD, cmdArgs...)

	// REF: https://github.com/aodin/go-pdf-server/blob/master/pdf_server.go

	// get output pipeline
	infoLog("wkhtmltopdf exec.Command > Getting Stdout pipe... ")
	output, err := wkhtmltopdf.StdoutPipe()
	if err != nil {
		errorLog("wkhtmltopdf exec.command Stdout Pipe err: ", err.Error())
		return err
	}

	// Begin the command
	infoLog("wkhtmltopdf exec.Command > Starting...")
	if err = wkhtmltopdf.Start(); err != nil {
		errorLog("wkhtmltopdf exec.command Start err: ", err.Error())
		return err
	}

	// Read the generated PDF from std out
	infoLog("wkhtmltopdf exec.Command > Reading output...")
	b, err := ioutil.ReadAll(output)
	if err != nil {
		errorLog("wkhtmltopdf ReadAll of output err: ", err.Error())
		return err
	}

	// End the command
	infoLog("wkhtmltopdf exec.Command > Waiting for command to exit...")
	if err = wkhtmltopdf.Wait(); err != nil {
		errorLog("wkhtmltopdf exec.Command Wait err: ", err.Error())
		return err
	}

	// write output to buffer
	pt.buf.Write(b)
	debugLog("Length of buffer for pdf output: ", pt.buf.Len(), " bytes")
	infoLog("wkhtmltopdf output has been written to buffer")

	return nil
}
