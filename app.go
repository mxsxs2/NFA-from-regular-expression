package main

import (
	"fmt"
)

/* Recommended structure
type nfa struct {
...
}
func regexcompile(r string) nfa {
...
return n
}
func (n nfa) regexmatch(n nfa, r sting) bool {
...
return ismatch
}
func main() {
n := regexcompile("01*0")
t := n.regexmatch("01110")
f := n.regexmatch("1000001")
}
*/
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

type nfa struct {
	initial *nfaNode
	accept  *nfaNode
}

//Function converts a positfoxed regular expression string into an NFA
func postfixToNFA(postfix string) *nfa {
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
	return nfaStack[0]
}

func main() {

	fmt.Printf("From a.(b.b)*.a to %s", convertInfixToPostfix("a.(b.b)*.a"))
	fmt.Printf("\nFrom 0.0.(1.1)*.0 to %s", convertInfixToPostfix("0.0.(1.1)*.0"))
	fmt.Printf("\n a.b.c* to %s", convertInfixToPostfix("a.b.c*"))
	//Examples from lecture video
	fmt.Printf("\n (a.(b|d))* to %s", convertInfixToPostfix("(a.(b|d))*"))
	fmt.Printf("\n a.(b|d).c* to %s", convertInfixToPostfix("a.(b|d).c*"))
	fmt.Printf("\n a.(b.b)+.c to %s\n", convertInfixToPostfix("a.(b.b)+.c"))

	nfa := postfixToNFA("ab.c*|")
	fmt.Println(nfa)
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
