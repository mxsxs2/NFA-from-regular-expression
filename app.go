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

func main() {

	fmt.Printf("From a.(b.b)*.a to %s", convertInfixToPostfix("a.(b.b)*.a"))
	fmt.Printf("\nFrom 0.0.(1.1)*.0 to %s", convertInfixToPostfix("0.0.(1.1)*.0"))
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
