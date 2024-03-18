package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Heiko-san/mermaidgen/flowchart"
	gographviz "github.com/awalterschulze/gographviz"
)

const (
	NShapeCylinder         = `[("%s")]`
	NShapeAsymmetric       = `>"%s"]`
	NShapeStadium          = `(["%s"])`
	NShapeHexagon          = `{{"%s"}}`
	NShapeTrapezoid        = `[/"%s"\]`
	NShapeTrapezoidAlt     = `[\"%s"/]`
	NShapeParallelogram    = `[/"%s"/]`
	NShapeParallelogramAlt = `[\"%s"\]`
	NShapeDoubleCircle     = `((("%s")))`
)

// ConvertDOTToMermaid преобразует входную строку в формате DOT и возвращает её в формате MermaidJS.
func ConvertDOTToMermaid(dotInput string) (string, error) {
	parsedGraph, err := gographviz.Parse([]byte(dotInput))
	if err != nil {
		return "", err
	}

	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(parsedGraph, graph); err != nil {
		return "", err
	}

	fc := flowchart.NewFlowchart()

	for _, node := range graph.Nodes.Nodes {
		nodeName := node.Name
		newNode := fc.AddNode(nodeName)
		// Применение стиля на основе атрибутов узла
		shape := node.Attrs["shape"]
		switch shape {
		case "cylinder":
			newNode.Shape = NShapeCylinder
		case "rarrow":
			newNode.Shape = NShapeAsymmetric // В MermaidJS нет точного эквивалента, это приблизительно
		case "octagon":
			newNode.Shape = NShapeHexagon // Используется ромб вместо октагона
		case "rectangle":
			newNode.Shape = flowchart.NShapeRect
		default:
			newNode.Shape = NShapeStadium // По умолчанию используем овал
		}

		// Применять текст метки к узлу таким образом:
		if label, exists := node.Attrs["label"]; exists && label != "" {
			newNode.AddLines(label)
		}
	}

	for _, edge := range graph.Edges.Edges {
		fromNode := fc.GetNode(edge.Src)
		toNode := fc.GetNode(edge.Dst)
		newEdge := fc.AddEdge(fromNode, toNode)

		label := stripQuotes(edge.Attrs["label"])
		if label != "" {
			newEdge.AddLines(label)
		}
	}

	return fc.String(), nil
}

func stripQuotes(s string) string {
	// Удаляем кавычки в начале и конце строки, если они есть
	s = strings.TrimPrefix(s, "\"")
	s = strings.TrimSuffix(s, "\"")

	// Заменяем задвоенные кавычки на одинарные
	return strings.ReplaceAll(s, "\"\"", "\"")
}

func main() {
	stdinFlag := flag.Bool("i", false, "read from STDIN")
	flag.Parse()

	var input string
	if *stdinFlag {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			input += scanner.Text() + "\n"
		}
	} else if len(flag.Args()) > 0 {
		fileName := flag.Arg(0)
		fileData, err := os.ReadFile(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read file: %v\n", err)
			return
		}
		input = string(fileData)
	} else {
		fmt.Fprintf(os.Stderr, "No input provided. Use -i flag to read DOT from STDIN or pass dot filename as argument: me dot filename\n")
		return
	}

	mermaidOutput, err := ConvertDOTToMermaid(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "conversion failed: %v\n", err)
		return
	}

	fmt.Println(mermaidOutput)
}
