package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type TreeNode struct {
	ID       string
	Node     string
	ParentID string
	Children []*TreeNode
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter CSV file path: ")
	filePath, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	filePath = strings.TrimSpace(filePath)

	tree, err := buildTreeFromCSV(filePath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Enter search method (DFS/BFS): ")
	method, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	method = strings.TrimSpace(strings.ToUpper(method))

	if method != "DFS" && method != "BFS" {
		log.Fatal("Invalid search method. Please choose DFS or BFS.")
	}

	fmt.Print("Enter search key (leave blank for complete tree): ")
	key, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	key = strings.TrimSpace(key)

	if method == "DFS" {
		result := DFS(tree, key)
		fmt.Println(strings.Join(result, ", "))
	} else {
		result := BFS(tree, key)
		fmt.Println(strings.Join(result, ", "))
	}
}

func buildTreeFromCSV(filename string) (*TreeNode, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// Read and validate headers
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}
	if len(headers) < 3 || headers[0] != "ID" || headers[1] != "NODE" || headers[2] != "PARENT" {
		return nil, fmt.Errorf("invalid CSV format")
	}

	treeMap := make(map[string]*TreeNode)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(row) < 3 {
			return nil, fmt.Errorf("invalid CSV format")
		}

		id := strings.TrimSpace(row[0])
		node := strings.TrimSpace(row[1])
		parentID := strings.TrimSpace(row[2])

		if !isValidID(id) || !isValidNode(node) || !isValidID(parentID) {
			return nil, fmt.Errorf("invalid characters in CSV data")
		}

		treeNode := &TreeNode{
			ID:       id,
			Node:     node,
			ParentID: parentID,
		}

		treeMap[id] = treeNode
	}

	// Build the tree structure based on the parent-child relationships
	var root *TreeNode
	for _, node := range treeMap {
		if node.ParentID == "0" {
			root = node
		} else {
			parentNode, ok := treeMap[node.ParentID]
			if !ok {
				return nil, fmt.Errorf("parent node with ID '%s' not found", node.ParentID)
			}
			parentNode.Children = append(parentNode.Children, node)
		}
	}

	return root, nil
}

func isValidID(id string) bool {
	return id != ""
}

func isValidNode(node string) bool {
	return node != ""
}

func DFS(root *TreeNode, key string) []string {
	if root == nil {
		return nil
	}

	result := make([]string, 0)
	stack := []*TreeNode{root}

	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		result = append(result, node.Node)

		if key != "" && strings.Contains(node.Node, key) {
			break
		}

		for i := len(node.Children) - 1; i >= 0; i-- {
			stack = append(stack, node.Children[i])
		}
	}

	return result
}

func BFS(root *TreeNode, key string) []string {
	if root == nil {
		return nil
	}

	result := make([]string, 0)
	queue := []*TreeNode{root}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		result = append(result, node.Node)

		if key != "" && strings.Contains(node.Node, key) {
			break
		}

		for _, child := range node.Children {
			queue = append(queue, child)
		}
	}

	return result
}
