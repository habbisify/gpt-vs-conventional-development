package main

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
)

type Node struct {
	Id     int32  `csv:"ID"` // .csv column headers
	Name   string `csv:"NODE"`
	Parent int32  `csv:"PARENT"`
}

type TreeNode struct {
	Id          int32
	Name        string
	Parent      int32
	NextSibling int32
	FirstChild  int32
}

func (n TreeNode) AddNextSibling(s int32) TreeNode {
	n.NextSibling = s
	return n
}

func (n TreeNode) AddFirstChild(c int32) TreeNode {
	n.FirstChild = c
	return n
}

func itemInSlice(item TreeNode, list []TreeNode) bool {
	for _, listItem := range list {
		if item == listItem {
			return true
		}
	}
	return false
}

func main() {
	// fmt.Println("Elapsed time: 2:55:00")

	// Open csv file
	in, err := os.Open("advanced_test.csv")
	if err != nil {
		panic(err)
	}
	defer in.Close()

	// Create Nodes to unpack contents of CSV file
	Nodes := []*Node{}

	// Unpack Nodes into Nodes
	if err := gocsv.UnmarshalFile(in, &Nodes); err != nil {
		panic(err)
	}

	// Create a non-binary tree from Nodes using TreeNode:s for map Tree
	var Tree = map[int32]TreeNode{}

	var rootId = int32(-1)

	for _, Node := range Nodes {
		//fmt.Println("Hello, ", Node.Id, Node.Name, Node.Parent)
		// Assign this new element to Tree
		Tree[Node.Id] = TreeNode{Node.Id, Node.Name, Node.Parent, -1, -1}

		if Node.Parent == 0 {
			rootId = Node.Id
		} else {
			if Tree[Node.Parent].FirstChild == -1 {
				// Modify Parent's FirstChild if not any yet
				Tree[Node.Parent] = Tree[Node.Parent].AddFirstChild(Node.Id)
			} else {
				// Else search last of childs using NextSibling:s
				prevId := Tree[Node.Parent].FirstChild
				for {
					thisId := Tree[prevId].NextSibling
					if thisId == -1 {
						Tree[prevId] = Tree[prevId].AddNextSibling(Node.Id)
						break
					}
					prevId = thisId
				}
			}
		}
	}

	var endKey = "Up"
	// DFS
	Item := Tree[rootId]
	Visited := []TreeNode{}

	for {
		if itemInSlice(Item, Visited) {
			if Item.NextSibling != -1 {
				Visited = Visited[:len(Visited)-1]
				Item = Tree[Item.NextSibling]

			} else {
				Visited = Visited[:len(Visited)-1]
				Item = Tree[Item.Parent]
			}
		} else {
			fmt.Println(Item.Name)
			Visited = append(Visited, Item)

			if Item.Name == endKey || Item.Id == 0 {
				break
			} else if Item.FirstChild != -1 {
				Item = Tree[Item.FirstChild]

			} else if Item.NextSibling != -1 {
				Visited = Visited[:len(Visited)-1]
				Item = Tree[Item.NextSibling]

			} else {
				Visited = Visited[:len(Visited)-1]
				Item = Tree[Item.Parent]
			}
		}

		//fmt.Println(Visited)
	}

	fmt.Println("")
	// BFS
	Item = Tree[rootId]
	ThisLayer := []TreeNode{Item}
	NextLayer := []TreeNode{}
	EndKeyFlag := false

	for {
		for _, Item := range ThisLayer {
			fmt.Println(Item.Name)

			if Item.Name == endKey {
				EndKeyFlag = true
				break
			} else if Item.FirstChild == -1 {
				continue
			} else {
				NextLayer = append(NextLayer, Tree[Item.FirstChild])
				prevId := Tree[Item.Id].FirstChild
				for {
					thisId := Tree[prevId].NextSibling
					if thisId == -1 {
						break
					} else {
						NextLayer = append(NextLayer, Tree[thisId])
					}
					prevId = thisId
				}
			}
		}
		// fmt.Println(ThisLayer, NextLayer)
		if len(NextLayer) == 0 || EndKeyFlag {
			break
		}
		ThisLayer = NextLayer
		NextLayer = []TreeNode{}
	}

}
