package file

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"parecer-gen/pkg/date"
	"path/filepath"

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

func GenerateParecerHTML(user, creci, dateStr, content string) (*File, error) {
	if user == "" || creci == "" || dateStr == "" || content == "" {
		return nil, fmt.Errorf("missing data to generate parecer")
	}

	dateTime := date.StringToTime(dateStr)

	dateStrFormatted := date.TimeToBRString(dateTime)

	input := ParecerDataHTML{
		User:     user,
		CreciStr: creci,
		DateStr:  dateStrFormatted,
		Content:  content,
	}
	templatePath, err := filepath.Abs("web/templates/parecer.html")
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
