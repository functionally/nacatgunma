package rdf

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/cayleygraph/quad"
	"github.com/cayleygraph/quad/nquads"
	"github.com/piprate/json-gold/ld"
	"gonum.org/v1/gonum/graph/formats/rdf"
)

func ReadRdf(filename string, baseUri string, format string) (interface{}, error) {
	dataset, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions(baseUri)
	options.Format = format
	expandedDoc, err := proc.FromRDF(dataset, options)
	if err != nil {
		return nil, err
	}
	context := generateContext(expandedDoc)
	return proc.Compact(expandedDoc, context, options)
}

func generateContext(doc interface{}) map[string]interface{} {
	ctx := make(map[string]interface{})
	scanForPredicates(doc, ctx, make(map[string]bool))
	return ctx
}

func scanForPredicates(value interface{}, ctx map[string]interface{}, fnd map[string]bool) {
	switch v := value.(type) {
	case map[string]interface{}:
		for k, val := range v {
			if _, exists := fnd[k]; exists {
				break
			}
			if isIRI(k) {
				short := shorten(k)
				if _, exists := ctx[short]; !exists {
					for i := min(5, len(short)); i <= len(short); i++ {
						if _, exists = ctx[short[:i]]; !exists {
							fnd[k] = true
							ctx[short[:i]] = k
							break
						}
					}
				}
			}
			scanForPredicates(val, ctx, fnd)
		}
	case []interface{}:
		for _, item := range v {
			scanForPredicates(item, ctx, fnd)
		}
	}
}

func isIRI(k string) bool {
	if strings.HasPrefix(k, "@") {
		return false
	}
	return strings.Contains(k, ":")
}

func shorten(iri string) string {
	var sep string
	if strings.Contains(iri, "#") {
		sep = "#"
	} else {
		sep = "/"
	}
	parts := strings.Split(iri, sep)
	last := parts[len(parts)-1]
	if last == "" && len(parts) > 1 {
		last = parts[len(parts)-2]
	}
	return last
}

// Deprecated: work in progress on alternative workflow
func ReadNquads(filename string) ([]quad.Quad, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := nquads.NewReader(bufio.NewReader(file), true)
	var quads []quad.Quad
	for {
		q, err := reader.ReadQuad()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		quads = append(quads, q)
	}
	return quads, nil
}

// Deprecated: work in progress on alternative workflow
func ReadStatements(filename string) ([]rdf.Statement, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := rdf.NewDecoder(bufio.NewReader(file))
	var statements []rdf.Statement
	for {
		s, err := reader.Unmarshal()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		statements = append(statements, *s)
	}
	return statements, nil
}
