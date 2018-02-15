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
	'*': 3,
	'.': 2,
	'|': 1,
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

func main() {
	//Loop a string
	for _, c := range "*.|0" {
		//Write out precedences
		fmt.Println(getPrecedence(c))
	}
}
