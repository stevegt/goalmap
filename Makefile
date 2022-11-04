f=examples/food

all: 
	go run main.go < $(f).md > $(f).dot               
	dot -Tsvg -o $f.svg $f.dot                          

