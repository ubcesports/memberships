package mailer

import (
	"bytes"
	_ "embed"
	"html/template"
	"os"
	"strings"
)

//go:embed templates/base.html
var baseTemplateSrc string

var baseTemplate = template.Must(template.New("base").Parse(baseTemplateSrc))

// Row is a single label/value line in an email's info table
// (eg. "Tier" / "Lounge Membership").
type Row struct {
	Label  string
	Value  string
	IsLast bool
}

// EmailData is the content that fills the shared base email template.
// Rows and CTA are optional: a nil/empty Rows omits the info table, and an
// empty CTAText omits the button. LogoURL is filled in by RenderEmail and
// does not need to be set by callers.
type EmailData struct {
	Title      string
	Heading    string
	Subheading string
	Rows       []Row
	CTAText    string
	CTAURL     string
	LogoURL    string
}

// NewRows builds a Row slice from ordered label/value pairs and marks the
// last row so the template can skip its bottom divider.
//
// Example usage:
//
//	mailer.NewRows(
//		"Tier", "Lounge Membership",
//		"Amount paid", "$25.00 CAD",
//	)
func NewRows(labelsAndValues ...string) []Row {
	rows := make([]Row, 0, len(labelsAndValues)/2)
	for i := 0; i+1 < len(labelsAndValues); i += 2 {
		rows = append(rows, Row{Label: labelsAndValues[i], Value: labelsAndValues[i+1]})
	}
	if len(rows) > 0 {
		rows[len(rows)-1].IsLast = true
	}
	return rows
}

// RenderEmail fills the shared base template with the given content and
// returns the resulting HTML.
func RenderEmail(data EmailData) (string, error) {
	data.LogoURL = LogoURL()

	var buf bytes.Buffer
	if err := baseTemplate.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// FrontendURL returns the primary frontend origin (no trailing slash), used
// to build links and asset URLs inside emails.
//
// Requirements from .env file:
//   - FRONTEND_URL: comma separated list of allowed frontend origins; the
//     first one is treated as the canonical app URL.
func FrontendURL() string {
	first, _, _ := strings.Cut(os.Getenv("FRONTEND_URL"), ",")
	return strings.TrimRight(strings.TrimSpace(first), "/")
}

// LogoURL returns the absolute URL of the UBCEA logo, for use in email headers.
func LogoURL() string {
	return FrontendURL() + "/ubcea_logo.jpg"
}
