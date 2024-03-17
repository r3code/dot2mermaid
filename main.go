package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/awalterschulze/gographviz"
	"github.com/emicklei/dot"
)

func СonvertDOTToMermaid(dotInput string) (string, error) {
	gParsed, err := gographviz.Parse([]byte(dotInput))
	if err != nil {
		return "", fmt.Errorf("error parsing dot input: %w", err)
	}

	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(gParsed, graph); err != nil {
		return "", fmt.Errorf("error analysing graph: %w", err)
	}

	dGraph := dot.NewGraph(dot.Directed)
	// Используем map для хранения узлов по ID, но без указателей
	nodesMap := make(map[string]dot.Node)

	for _, node := range graph.Nodes.Nodes {
		nodeID := node.Name
		createdNode := dGraph.Node(nodeID)
		nodesMap[nodeID] = createdNode // Прямое сохранение узла в мапе
	}

	for _, edge := range graph.Edges.Edges {
		srcNode, srcExists := nodesMap[edge.Src]
		dstNode, dstExists := nodesMap[edge.Dst]
		if srcExists && dstExists { // Проверка, что узлы существуют в мапе
			// Теперь используем узлы напрямую
			label := edge.Attrs["label"]
			formattedLabel := formatLabel(label)
			dGraph.Edge(srcNode, dstNode).Attr("label", formattedLabel)
		}
	}

	mermaid := dot.MermaidGraph(dGraph, dot.MermaidTopToBottom)
	return mermaid, nil
}

func formatLabel(label string) string {
	// В MermaidJS пустые метки могут быть полностью опущены или должны содержать пробел.
	if label == "" {
		return " " // Использовать пробельный символ внутри метки для явного отображения пустой метки.
	}
	return fmt.Sprintf("%s", label) // Возвращаем метку, заключенную в кавычки, для не пустых меток.
}

func main() {
	stdinFlag := flag.Bool("i", false, "Read from STDIN")
	flag.Parse()

	var input string
	if *stdinFlag {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			input += scanner.Text() + "\n"
		}
	} else {
		fmt.Println("Reading directly from a file is not supported in this version. Please pipe file content into the program.")
		os.Exit(1)
	}

	mermaidOutput, err := СonvertDOTToMermaid(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to convert DOT to Mermaid: %v", err)
		os.Exit(1)
	}

	fmt.Println(mermaidOutput)
}
