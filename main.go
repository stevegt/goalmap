package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/goccy/go-graphviz"
	. "github.com/stevegt/goadapt"

	// "github.com/goccy/go-graphviz"
	// "github.com/yuin/goldmark"
	"github.com/bbrks/wrap"
	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	gtext "github.com/yuin/goldmark/text"
)

func main() {
	// read the input file from stdin
	in, err := ioutil.ReadAll(os.Stdin)
	Ck(err)
	// convert the markdown to graphviz dot format
	out := md2graphviz(in)
	// write the output to stdout
	fmt.Println(string(out))
}

// md2graphviz converts a set of nested bullet lists from markdown to
// graphviz dot format
func md2graphviz(in []byte) (out []byte) {
	// create a new goldmark object
	markdown := goldmark.New()
	// convert the input buffer to a goldmark text.Reader
	reader := gtext.NewReader(in)
	// parse the markdown
	ast := markdown.Parser().Parse(reader)
	// convert the AST to an internal graph representation
	graph := NewGraph()
	ast2graph(in, ast, graph)
	// convert the graph to graphviz dot format
	out = graph2dot(graph)
	return
}

type Node struct {
	// the node's name
	name string
	// the list of tail nodes that point to this head node
	tails []*Node
}

type Graph struct {
	// map of node name to node
	Nodes map[string]*Node
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
	}
}

// AddNode adds a node to the graph
func (g *Graph) AddNode(name string) (node *Node) {
	// if the node already exists
	node, ok := g.Nodes[name]
	if ok {
		return
	}
	node = &Node{name: name}
	g.Nodes[name] = node
	return
}

// AddLink adds a link to the graph
func (g *Graph) AddLink(head, tail *Node) {
	// Pf("AddLink head '%v' tail '%v'\n", head.name, tail.name)
	// assert that the head and tail nodes exist
	_, ok := g.Nodes[head.name]
	Assert(ok, "head node does not exist", head.name)
	_, ok = g.Nodes[tail.name]
	Assert(ok, "tail node does not exist", tail.name)
	// add the link
	head.tails = append(head.tails, tail)
}

// nodeKind returns true if the node is of the specified kind
// XXX this is probably the wrong way to do this but for some reason
// the right way is not obvious from goldmark's documentation
func nodeKind(node gast.Node, kind string) bool {
	return node.Kind().String() == kind
}

func getText(src []byte, node gast.Node) string {
	// iterate over the node's children
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		// Pf("child: %v\n", child.Kind())
		// if the child is a text node
		if nodeKind(child, "TextBlock") {
			// get the text
			text := string(child.Text(src))
			// wrap name at 30 characters
			// XXX this should be a command line option
			text = wrap.Wrap(text, 30)
			return text
		}
	}
	return ""
}

// ast2graph converts a goldmark AST to an internal graph representation
func ast2graph(src []byte, ast gast.Node, graph *Graph) {
	// if the node is a list item
	if nodeKind(ast, "ListItem") {
		tailText := getText(src, ast)
		// add the node to the graph
		tail := graph.AddNode(tailText)
		// if the node has a grandparent
		if ast.Parent().Parent() != nil {
			// get the parent's text
			headText := getText(src, ast.Parent().Parent())
			if headText != "" {
				// add the parent to the graph
				head := graph.AddNode(headText)
				// add the link to the graph
				graph.AddLink(head, tail)
			}
		}
	}
	// for each child node
	for child := ast.FirstChild(); child != nil; child = child.NextSibling() {
		// recursively process the child node
		ast2graph(src, child, graph)
	}
}

// graph2dot converts an internal graph representation to graphviz dot format
func graph2dot(graph *Graph) (dot []byte) {
	// Pprint(graph)

	// create a new graphviz graph
	g := graphviz.New()
	gv, err := g.Graph()
	Ck(err)
	gv.SetRankDir("TB")
	// iterate over the nodes in the graph
	for _, inthead := range graph.Nodes {
		// create a new graphviz node
		head, err := gv.CreateNode(inthead.name)
		Ck(err)
		// iterate over the tails in the node
		for _, tail := range inthead.tails {
			// create a new graphviz node
			tail, err := gv.CreateNode(tail.name)
			// create a new graphviz edge
			_, err = gv.CreateEdge("", tail, head)
			Ck(err)
		}
	}
	// render the graph to dot format
	var buf bytes.Buffer
	err = g.Render(gv, "dot", &buf)
	Ck(err)
	dot = buf.Bytes()
	return
}
