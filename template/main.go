package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"
	"time"
)

func main1() {
	// Define a template.
	const letter = `
Dear ${.Name},
${if .Attended}
It was a pleasure to see you at the wedding.
${else}
It is a shame you couldn't make it to the wedding.
${end}
${with .Gift }
Thank you for the lovely ${.}.
${end}
Best wishes,
Josie
`

	// Prepare some data to insert into the template.
	type Recipient struct {
		Name, Gift string
		Attended   bool
	}
	var recipients = []Recipient{
		{"Aunt Mildred", "bone china tea set", true},
		{"Uncle John", "moleskin pants", false},
		{"Cousin Rodney", "", false},
	}

	// Create a new template and parse the letter into it.
	t := template.Must(template.New("letter").Delims("${", "}").Parse(letter))

	// Execute the template for each recipient.
	for _, r := range recipients {
		err := t.Execute(os.Stdout, r)
		if err != nil {
			log.Println("executing template:", err)
		}
	}
}

var now time.Time

func year() string {
	return fmt.Sprintf("%04d", now.Year())
}

func mon() string {
	return fmt.Sprintf("%02d", now.Month())
}

func day() string {
	return fmt.Sprintf("%02d", now.Day())
}

func hour() string {
	return fmt.Sprintf("%02d", now.Hour())
}

func min() string {
	return fmt.Sprintf("%02d", now.Minute())
}

func sec() string {
	return fmt.Sprintf("%02d", time.Now().Second())
}

func main() {
	// First we create a FuncMap with which to register the function.
	funcMap := template.FuncMap{
		"year": year,
		"mon":  mon,
		"day":  day,
		"hour": hour,
		"min":  min,
		"sec":  sec,
	}

	// A simple template definition to test our function.
	// We print the input text several ways:
	// - the original
	// - title-cased
	// - title-cased and then printed with %q
	// - printed with %q and then title-cased.
	const templateText = `kodo-parquet/date=$( year  )-$(mon)-$( day)/hour=$(hour)/min=$(min)/$(sec)-`

	// Create a template, add the function map, and parse the text.
	tmpl, err := template.New("titleTest").Delims("$(", ")").Funcs(funcMap).Parse(templateText)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}

	now = time.Now()

	// Run the template to verify the output.
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, "")
	if err != nil {
		log.Fatalf("execution: %s", err)
	}
	fmt.Println(buf.String())

	time.Sleep(time.Second * 10)

	buf.Reset()
	now = time.Now()
	err = tmpl.Execute(&buf, "")
	if err != nil {
		log.Fatalf("execution: %s", err)
	}
	fmt.Println(buf.String())
}
