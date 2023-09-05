package parsing

import (
	"fmt"
	"strings"

	"github.com/khulnasoft-lab/package-analysis/internal/staticanalysis/token"
	"github.com/khulnasoft-lab/package-analysis/internal/utils"
)

// Language represents a programming language used for parsing.
type Language string

const (
	NoLanguage Language = ""
	JavaScript Language = "JavaScript"
)

var allLanguages = []Language{JavaScript}

func SupportedLanguages() []Language {
	return allLanguages[:]
}

// tokenType denotes types of source code tokens collected during parsing (see token package).
type tokenType string

// statusType denotes a type of status reported by the parser about the parsing process.
type statusType string

const (
	identifier tokenType  = "Identifier" // source code identifier (variable, class, function name)
	literal    tokenType  = "Literal"    // source code data (string, integer, floating point literals)
	comment    tokenType  = "Comment"    // source code comments
	parseInfo  statusType = "Info"       // information about the parsing (e.g. number of bytes read by parser)
	parseError statusType = "Error"      // any error encountered by parser; some are recoverable and some are not
)

type parsedIdentifier struct {
	Type token.IdentifierType
	Name string
	Pos  token.Position
}

func (i parsedIdentifier) String() string {
	return fmt.Sprintf("%s %s [pos %d:%d]", i.Type, i.Name, i.Pos.Row(), i.Pos.Col())
}

type parsedLiteral[T any] struct {
	Type     string
	GoType   string
	Value    T
	RawValue string
	InArray  bool
	Pos      token.Position
}

func (l parsedLiteral[T]) String() string {
	s := fmt.Sprintf("%s (%s) %v (raw: %s) pos %d:%d", l.Type, l.GoType, l.Value, l.RawValue, l.Pos.Row(), l.Pos.Col())
	if l.InArray {
		s += " [array]"
	}
	return s
}

type parsedComment struct {
	Type string
	Data string
	Pos  token.Position
}

func (c parsedComment) String() string {
	return fmt.Sprintf("%s %s pos %d:%d", c.Type, c.Data, c.Pos.Row(), c.Pos.Col())
}

type parserStatus struct {
	Type    statusType
	Name    string
	Message string
	Pos     token.Position
}

func (s parserStatus) String() string {
	return fmt.Sprintf("[%s] %s: %s pos %d:%d", s.Type, s.Name, s.Message, s.Pos.Row(), s.Pos.Col())
}

// singleParseData holds package-internal data for a single file processed by a single language parser.
type singleParseData struct {
	ValidInput  bool
	Identifiers []parsedIdentifier
	Literals    []parsedLiteral[any]
	Comments    []parsedComment
	Info        []parserStatus
	Errors      []parserStatus
}

func (d singleParseData) String() string {
	identifiers := utils.Transform(d.Identifiers, func(pi parsedIdentifier) string { return pi.String() })
	literals := utils.Transform(d.Literals, func(pl parsedLiteral[any]) string { return pl.String() })
	comments := utils.Transform(d.Comments, func(c parsedComment) string { return c.String() })
	info := utils.Transform(d.Info, func(i parserStatus) string { return i.String() })
	errors := utils.Transform(d.Errors, func(e parserStatus) string { return e.String() })

	parts := []string{
		"== Identifiers ==",
		strings.Join(identifiers, "\n"),
		"== Literals ==",
		strings.Join(literals, "\n"),
		"== Comments ==",
		strings.Join(comments, "\n"),
		"== Info ==",
		strings.Join(info, "\n"),
		"== Errors ==",
		strings.Join(errors, "\n"),
	}

	return strings.Join(parts, "\n")
}
