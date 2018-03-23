package main

import (
	"fmt"
	"strings"
)

//Map of the special characters with their precedence
var precedenceMap = map[rune]int{
	'+': 3, //One or more
	'*': 3, //Kleene star
	'?': 3, //Zero or one
	'.': 2, //Concatenate
	'|': 1, //Alternate (either side but not both)
	'(': 0,
}

//Structure of a node of an NFA
type nfaNode struct {
	symbol rune
	edge1  *nfaNode
	edge2  *nfaNode
}

//Stroage for the start and end node of the NFA linked list
type nfa struct {
	initial *nfaNode
	accept  *nfaNode
}

func main() {
	//Accept all strings of 0's and 1;s that begin with two zeros
	n := regexcompile("0.0.(0|1)*")
	fmt.Println("Regex: 0.0.(0|1)*")
	fmt.Println("001110 passes:", n.regexmatch("001110"))
	fmt.Println("01110 passes:", n.regexmatch("01110"))
	fmt.Println("1000001 passes:", n.regexmatch("1000001"))
	fmt.Println("001000001 passes:", n.regexmatch("001000001"))
}

//Compiles a regex string into an NFA linked list
func regexcompile(r string) nfa {
	//Declare the nfa
	var nfa nfa
	//Trim the white spaces from the strng
	nr := strings.TrimSpace(r)
	//Make sure the string contains something
	if len(nr) > 0 {
		//Convert the infixed regex to postfix
		nr = convertInfixToPostfix(nr)
		//Create the NFA linked list from the postifxed regex
		nfa = postfixToNFA(nr)
	}
	//Return the nfa
	return nfa
}

//Function used to match a string to a regex(nfa structure)
func (nfa nfa) regexmatch(input string) bool {
	//State of the macth
	ismatch := false
	//Create the current node array with the starting nodes
	currentNodes := []*nfaNode{}
	//Nodes that the current node points to
	nextNodes := []*nfaNode{}

	//Fill the current nodes from the nfa's starting node and adding the accept node
	currentNodes = addState(currentNodes[:], nfa.initial, nfa.accept)

	//Loop the input string
	for _, c := range input {
		//Loop the current nodes array
		for _, currentNode := range currentNodes {
			//If the current node's sybmol is the same as the current character from the inout string
			if currentNode.symbol == c {
				//Get the nodes which the current node points to
				nextNodes = addState(nextNodes[:], currentNode.edge1, nfa.accept)
			}
		}
		//Swap the current nodes with the next nodes and create an empty array for the next nodes
		currentNodes, nextNodes = nextNodes, []*nfaNode{}
	}

	//Loop the current nodes to determine if any of them is an accept node
	for _, currentNode := range currentNodes {
		//If the current node is an accept node
		if currentNode == nfa.accept {
			//Set the accept state of this regex matching
			ismatch = true
			break
		}
	}

	//Return the result state
	return ismatch
}

//Function converts a positfoxed regular expression string into an NFA
func postfixToNFA(postfix string) nfa {
	//Create empty nfa stack
	nfaStack := []*nfa{}

	//Loop the regular expression
	for _, r := range postfix {
		switch r {
		case '.':
			//Pop the second fragment form the NFA stack
			frag2 := nfaStack[len(nfaStack)-1]
			nfaStack = nfaStack[:len(nfaStack)-1]
			//Pop the fisrt fragment form the NFA stack
			frag1 := nfaStack[len(nfaStack)-1]
			nfaStack = nfaStack[:len(nfaStack)-1]
			//Link(concatenate) the two nfas to each other.
			frag1.accept.edge1 = frag2.initial
			//Push the new fragmenst back to the NFA stack
			nfaStack = append(nfaStack, &nfa{initial: frag1.initial, accept: frag2.accept})
		case '|':
			//Pop the second fragment form the NFA stack
			frag2 := nfaStack[len(nfaStack)-1]
			nfaStack = nfaStack[:len(nfaStack)-1]
			//Pop the fisrt fragment form the NFA stack
			frag1 := nfaStack[len(nfaStack)-1]
			nfaStack = nfaStack[:len(nfaStack)-1]
			//Create an empty accept node
			accept := nfaNode{}
			//Jin the two nodes into one node
			initial := nfaNode{edge1: frag1.initial, edge2: frag2.initial}
			//Set the edges of both fragments to the accept node as ether of em should be accepted
			frag1.accept.edge1 = &accept
			frag2.accept.edge1 = &accept
			//Push the new initial and accept nodes to the NFA stack
			nfaStack = append(nfaStack, &nfa{initial: &initial, accept: &accept})
		case '*':
			//Pop the last fragment form the NFA stack
			frag := nfaStack[len(nfaStack)-1]
			nfaStack = nfaStack[:len(nfaStack)-1]
			//Create an empty accept node
			accept := nfaNode{}
			initial := nfaNode{edge1: frag.initial, edge2: &accept}
			frag.accept.edge1 = frag.initial
			frag.accept.edge2 = &accept
			//Push the new initial and accept nodes to the NFA stack
			nfaStack = append(nfaStack, &nfa{initial: &initial, accept: &accept})
		default:
			//Create an empty accept node
			accept := nfaNode{}
			//Add the current non special character to a new node
			initial := nfaNode{symbol: r, edge1: &accept}
			//Push the new nodes to the stack
			nfaStack = append(nfaStack, &nfa{initial: &initial, accept: &accept})
		}
	}
	//Return the final nfa
	return *nfaStack[0]
}

//Gets the nodes which "s" current node points to if the s is not the same as the "a" accept node
func addState(l []*nfaNode, s *nfaNode, a *nfaNode) []*nfaNode {
	//Add the current node to the inital list of nodes
	l = append(l, s)
	//Determine if the node is symbol node or not and the current state is not the same as accept state
	if s != a && s.symbol == 0 {
		//Add the node from the edge 1
		l = addState(l, s.edge1, a)
		if s.edge2 != nil {
			l = addState(l, s.edge2, a)
		}
	}
	//Return the nodes list
	return l
}

// Finds the precedence of a character
func getPrecedence(c rune) int {
	//Try to get the value from the map
	value, found := precedenceMap[c]
	//If the character is in the map
	if found {
		//Return the assigned precedence
		return value
	}
	//If the value was not in the map then return the highest precedence
	return len(precedenceMap) + 1
}

//Function used to convert infixed regex string to postfixed one
//The algorithm is using stacks for the conversion
func convertInfixToPostfix(input string) string {
	//Slice for the output string
	output := []rune{}
	//Slice for storing the current characters
	stack := []rune{}

	//Loop the input string
	for _, char := range input {
		//If the character is an open bracket
		if char == '(' {
			//Add the character to the stack
			stack = append(stack, char)
		} else if char == ')' { //If the character is a closing bracket
			//Loop the stack until the opening bracket is found
			for len(stack)-1 >= 0 && stack[len(stack)-1] != '(' {
				//Add to the output
				output = append(output, stack[len(stack)-1])
				//Remove the top item from the stack
				stack = stack[:len(stack)-1]
			}
			//Remove the opening bracket from the stack
			//The closing bracket was never added to the stack
			if len(stack)-1 >= 0 {
				stack = stack[:len(stack)-1]
			}

		} else { //If the character is not an opening or closing bracket
			//while there are tokens to be read
			for i := len(stack) - 1; i >= 0; i-- {
				//Get the top of the stack
				var top = stack[len(stack)-1]
				//Check if the current character has a lower precedence than the one on the top of the stack
				if getPrecedence(top) >= getPrecedence(char) {
					//Add to the output
					output = append(output, top)
					//Remove the top item from the stack
					stack = stack[:len(stack)-1]
				} else {
					break
				}
			}
			//Add the character to the stack
			stack = append(stack, char)
		}

	}
	//Add the remainder of the stack
	for i := len(stack) - 1; i >= 0; i-- {
		//Add to the output
		output = append(output, stack[len(stack)-1])
		//Remove the top item from the stack
		stack = stack[:len(stack)-1]
	}

	//Convert the rune slice to a string and return it
	return string(output)
}
