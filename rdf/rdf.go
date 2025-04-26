package rdf

import (
	"bufio"
	"fmt"
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
	fmt.Println(expandedDoc)
	context := generateContext(expandedDoc)
	fmt.Println(context)
	return proc.Compact(expandedDoc, context, options)
}

func generateContext(doc interface{}) map[string]interface{} {
	ctx := make(map[string]interface{})
	scanForPredicates(doc, ctx)
	return ctx
}

func scanForPredicates(value interface{}, ctx map[string]interface{}) {
	switch v := value.(type) {
	case map[string]interface{}:
		for k, val := range v {
			if isIRI(k) {
				short := shorten(k)
				if _, exists := ctx[short]; !exists {
					ctx[short] = k
				}
			}
			scanForPredicates(val, ctx)
		}
	case []interface{}:
		for _, item := range v {
			scanForPredicates(item, ctx)
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
	parts := strings.Split(iri, "/")
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
