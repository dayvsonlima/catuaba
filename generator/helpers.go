package generator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/jinzhu/inflection"
)

func Pluralize(in string) string {
	return inflection.Plural(in)
}

// goAcronyms maps common Go acronyms that should be fully uppercased.
var goAcronyms = map[string]string{
	"id": "ID", "url": "URL", "api": "API", "http": "HTTP",
	"https": "HTTPS", "html": "HTML", "css": "CSS", "json": "JSON",
	"xml": "XML", "sql": "SQL", "ssh": "SSH", "tcp": "TCP",
	"udp": "UDP", "ip": "IP", "uri": "URI", "uuid": "UUID",
	"uid": "UID", "cpu": "CPU", "gpu": "GPU", "os": "OS",
	"db": "DB", "io": "IO", "vm": "VM",
}

func Camelize(in string) string {
	parts := strings.Split(in, "_")
	var out strings.Builder

	for _, part := range parts {
		if part == "" {
			continue
		}
		if acronym, ok := goAcronyms[strings.ToLower(part)]; ok {
			out.WriteString(acronym)
		} else {
			out.WriteRune(unicode.ToUpper(rune(part[0])))
			out.WriteString(part[1:])
		}
	}

	return out.String()
}

func CamelizeVar(in string) string {
	parts := strings.Split(in, "_")
	var out strings.Builder

	for i, part := range parts {
		if part == "" {
			continue
		}
		if i == 0 {
			if acronym, ok := goAcronyms[strings.ToLower(part)]; ok {
				out.WriteString(strings.ToLower(acronym))
			} else {
				out.WriteRune(unicode.ToLower(rune(part[0])))
				out.WriteString(part[1:])
			}
		} else {
			if acronym, ok := goAcronyms[strings.ToLower(part)]; ok {
				out.WriteString(acronym)
			} else {
				out.WriteRune(unicode.ToUpper(rune(part[0])))
				out.WriteString(part[1:])
			}
		}
	}

	return out.String()
}

func Snakeze(str string) string {
	// Handle transitions like "ConversationID" -> "Conversation_ID"
	matchFirstCap := regexp.MustCompile("([A-Z]+)([A-Z][a-z])")
	matchLowerUpper := regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchLowerUpper.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func LowerPlural(in string) string {
	return strings.ToLower(Pluralize(in))
}

func GetAttributeName(in string) string {

	attribute := strings.Split(in, ":")
	attributeName := Camelize(attribute[0])

	return attributeName
}

func GetAttributeType(in string) string {
	attribute := strings.Split(in, ":")
	raw := attribute[1]

	// Map user-facing types to Go types
	switch raw {
	case "text":
		return "string"
	case "integer":
		return "int"
	case "float":
		return "float64"
	case "boolean":
		return "bool"
	case "datetime":
		return "time.Time"
	default:
		return raw
	}
}

// GetAttributeRawType returns the original type as specified by the user (e.g. "text", "string", "bool").
func GetAttributeRawType(in string) string {
	attribute := strings.Split(in, ":")
	return attribute[1]
}

func GetAttributeJson(in string) string {
	name := GetAttributeName(in)
	name = Snakeze(name)

	return fmt.Sprintf("`json:\"%s\"`", name)
}

func GetAttributeJsonBinding(in string) string {
	name := GetAttributeName(in)
	name = Snakeze(name)
	typ := GetAttributeType(in)

	// Booleans: binding:"required" rejects false values in Gin
	if typ == "bool" {
		return fmt.Sprintf("`json:\"%s\"`", name)
	}
	return fmt.Sprintf("`json:\"%s\" binding:\"required\"`", name)
}

func GetFormBinding(in string) string {
	name := GetAttributeName(in)
	name = Snakeze(name)

	return fmt.Sprintf("`form:\"%s\"`", name)
}

func GetFormLabel(in string) string {
	name := GetAttributeName(in)
	// Insert spaces before capitals: "PostId" -> "Post Id"
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")
	label := matchFirstCap.ReplaceAllString(name, "${1} ${2}")
	label = matchAllCap.ReplaceAllString(label, "${1} ${2}")
	return label
}

// GetSQLType maps a user-facing attribute type to a SQL column type.
// Uses generic SQL that works across SQLite, PostgreSQL, and MySQL.
func GetSQLType(in string) string {
	attribute := strings.Split(in, ":")
	typ := attribute[1]

	switch typ {
	case "string":
		return "VARCHAR(255)"
	case "text":
		return "TEXT"
	case "int", "integer":
		return "INTEGER"
	case "uint":
		return "INTEGER"
	case "float64", "float":
		return "REAL"
	case "bool", "boolean":
		return "BOOLEAN"
	case "datetime", "time.Time":
		return "TIMESTAMP"
	default:
		return "VARCHAR(255)"
	}
}

// GetSQLDefault returns the default value for a SQL column type.
func GetSQLDefault(in string) string {
	attribute := strings.Split(in, ":")
	typ := attribute[1]

	switch typ {
	case "string", "text":
		return "''"
	case "int", "integer", "uint", "float64", "float":
		return "0"
	case "bool", "boolean":
		return "0"
	default:
		return "''"
	}
}

// GetSQLColumnName returns the snake_case column name for an attribute.
func GetSQLColumnName(in string) string {
	name := GetAttributeName(in)
	return Snakeze(name)
}

func GetFormInputType(in string) string {
	attribute := strings.Split(in, ":")
	typ := attribute[1]

	switch typ {
	case "int", "integer", "uint", "float64", "float":
		return "number"
	case "bool", "boolean":
		return "checkbox"
	case "text":
		return "textarea"
	case "datetime", "time.Time":
		return "datetime-local"
	default:
		return "text"
	}
}
