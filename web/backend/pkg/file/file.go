package file

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/google/uuid"
)

type File struct {
	Filename string
	Reader   io.Reader
}

type ParecerDataHTML struct {
	User     string
	CreciStr string
	DateStr  string
	Content  string
}

func GenerateParecerHTML(user, creci, content string, date time.Time) (*File, error) {
	if user == "" || creci == "" || content == "" {
		return nil, fmt.Errorf("missing data to generate parecer")
	}

	dateStr := date.Format("01-02-2006")

	input := ParecerDataHTML{
		User:     user,
		CreciStr: creci,
		DateStr:  dateStr,
		Content:  content,
	}
	templatePath, err := filepath.Abs("templates/parecer.html")
	if err != nil {
		return nil, fmt.Errorf("error getting absolute path for template: %v", err)
	}

	templateParecer, err := template.ParseFiles(templatePath)

	if err != nil {
		return nil, fmt.Errorf("error parsing template: %v", err)
	}

	var buf bytes.Buffer

	writer := bufio.NewWriter(&buf)

	if err := templateParecer.Execute(writer, input); err != nil {
		return nil, fmt.Errorf("error executing template: %v", err)
	}

	if err := writer.Flush(); err != nil {
		return nil, fmt.Errorf("error flushing writer: %v", err)
	}

	filename := fmt.Sprintf("parecer-%s-%s", dateStr, uuid.New().String())

	htmlReader := bytes.NewReader(buf.Bytes())

	return &File{
		Filename: filename,
		Reader:   htmlReader,
	}, nil
}

func GeneratePDF(file *File) (*File, error) {
	page := wkhtmltopdf.NewPageReader(file.Reader)

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, fmt.Errorf("error creating PDF generator: %v", err)
	}

	pdfg.AddPage(page)

	if err = pdfg.Create(); err != nil {
		return nil, fmt.Errorf("error generating PDF: %v", err)
	}

	filename := fmt.Sprintf("%s.pdf", file.Filename)

	pdfReader := bytes.NewReader(pdfg.Bytes())

	return &File{
		Filename: filename,
		Reader:   pdfReader,
	}, nil
}
