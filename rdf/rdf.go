package rdf

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/cayleygraph/quad"
	"github.com/cayleygraph/quad/nquads"
	"github.com/piprate/json-gold/ld"
	"gonum.org/v1/gonum/graph/formats/rdf"
)

func ReadRdf(filename string, baseUri string) (interface{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	dataset, err := parseNQuads(bufio.NewReader(file))
	if err != nil {
		return nil, err
	}
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions(baseUri)
	options.Format = "application/n-quads"
	expandedDoc, err := proc.FromRDF(dataset, options)
	if err != nil {
		return nil, err
	}
	context := generateContext(expandedDoc)
	return proc.Compact(expandedDoc, context, options)
}

func parseNQuads(r io.Reader) (string, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
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
