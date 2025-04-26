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
	graph, ok := doc.([]interface{})
	if !ok {
		return ctx
	}
	for _, node := range graph {
		nodeMap, ok := node.(map[string]interface{})
		if !ok {
			continue
		}
		for k := range nodeMap {
			if strings.HasPrefix(k, "http") {
				short := shorten(k)
				ctx[short] = k
			}
		}
	}
	return ctx
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
