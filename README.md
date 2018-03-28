# NFA-from-regular-expression
A simple regural expression processor written in golang for Graph Theory 2018 Project

## What is it?
The application is a regular expression processor. The processor takes the alphabet of ```∑ = {0, 1}``` and matches it to a set of rules which is represented by a regular expression string. The result of the matching is either ```true``` or ```false```. If the result is ```true```, the input string matches the regular expression rules. This indicates the acceptance by the underlying Nondeterministic Finite Automata(NFA), otherwise it is rejected by the NFA.
<br />
### Suppoted symbols(Required by project description)
#### Concatenation ```.```
Concatenates two characters for example ```0.1``` is ```01```
#### Alternate (OR) ```|```
Either side of the ```|``` sybmbol, but noth both. For example ```0|1``` means, either 0 or 1 but not both.
#### Kleene Star ```*```
Zero or more of the previous character. For example ```0*``` means, any amount of zeros.
### Additional supported symbols
#### One or more ```+```
One or more of the previous character but not zero amount. For example ```0+``` means, either a single zero or more of it.
#### Zero or one ```?```
It is the opposite of the ```+``` symbol. Either exactly one of the previous character or none of it.


### How to install and run GO

To install, simply go to GO's website and download the installer and run it: https://golang.org/

When the installation is done, the "go" command is going to be avaialable in terminal.(You have to restart an opened terminal)

### How to use this repository
To build the application, navigate to the desired folder and type this command: 
```
go build 
```

The previous command will compile the go project into a runnable.

Once the runnable is created then it can be run in terminal, for example: 
```
./NFA-from-regular-expression.exe
```

#### How to use the regular expression code
Every regular expression has to be compiled first. This guarantees better performance over compiling it every single time before use. Once the regex is compiled an NFA will be returned along an empty error message. If the message is not empty that means the regex could not be compiled.
<br />
Example usage:
```
	//Accept all strings of 0's and 1;s that begin with two zeros
	if n, err := regexcompile("0.0.(0|1)*"); err == nil {
		fmt.Println("Regex: 0.0.(0|1)*")
		fmt.Println("001110 passes:", n.regexmatch("001110"))
		fmt.Println("01110 passes:", n.regexmatch("01110"))
	} else {
		fmt.Println(err)
	}
```
The example will return:
```
Regex: 0.0.(0|1)
001110 passes:true
01110 passes:false
```

## How it works
The regular expression processor is built up from three main components
### 1. Infixed regular expression to postfixed
The first step of the processing is to convert infixed regular experssions to postfixed as a postfixed regular expression is easier to process programatically.
<br />
The conversion is done by the use of Shunting yard algorithm. This algorithm is based on "special character" precedence and two stacks.<br />
The infixed string is looped through character by character. Each of these characters are added to a temporary stack (expect "```)```") and to an output stack, depending on the precedence of the current character and the one on the top of the stack. If the character on the top of the temporary stack has a higher precedence then it is added to the output stack. 
<br />
Once all characters are checked then the output stack is returned as a new string. This string is the postfix notated version of the original string.
<br />
The character precedence is as follows:
```
//Map of the special characters with their precedence
var precedenceMap = map[rune]int{
	'+': 3, //One or more
	'*': 3, //Kleene star
	'?': 3, //Zero or one
	'.': 2, //Concatenate
	'|': 1, //Alternate (either side but not both)
	'(': 0,
}
```
### 2. Convert postfixed string into Nondeterministic Finite Automata(NFA)
A new structure is introduced to store each state of the NFA. Each structure has two edges to lead to another state, it also has a symbol which is only used if the character in the postfixed string is not a special character. The algorithm is based on a stack. Each stack item has an initial state and an accept state. During the algorithm these items are joined together.
<br />
The same special characters are used in this algorithm as in the infix to postfix converter.
<br />

#### Concatenation ```.```
Two items are popped from the top of the stack the "higher" item's accept state is connected to the "lower" items accept state, then as a new conjoined item, it is pushed back to the stack.

#### Alternate (OR) ```|```
Two items are popped from the top of the stack.
<br />
A new initial state is created with the two edges connected to the two stack items's initial states. An additional empty accept state is created and connected to the two popped items's accept state's edge.
<br />
A new item is pushed to the stack with the new initial state and the empty accept state.

#### Kleene Star ```*```
<br />
One item is popped from the top of the stack.
<br />
A new initial state is created. This state has the popped node's initial state as the first edge and an empty accept state as second edge.
The popped item's accept state's first edge is joined back to the original initial state and the second edge is joined to the empty accept state.
<br />
A new item is pushed to the stack with the new initial state and the empty accept state.

#### One or more ```+```
<br />
One item is popped from the top of the stack.
<br />
A new initial state is created. This state has the popped node's initial state as the first edge and an empty accept state as second edge.
The popped item's accept state's first edge is joined to the new initial state and the second edge is joined to the empty accept state.
<br />
A new item is pushed to the stack with the original item's initial state and the empty accept state.

#### Zero or one ```?```
<br />
One item is popped from the top of the stack.
<br />
A new initial state is created. This state has an empty accept state as the first edge and the popped node's initial state as second edge.
The popped item's accept state's edges are both joined to the empty accept state.
<br />
A new item is pushed to the stack with the new initial state and the empty accept state.

#### None of the above
If the current character is not a special character then a new initial state is created which has this character as the symbol and an empty accept state
<br />
A new item is pushed to the stack with the new initial state and the empty accept state.

#### Done 
If everything goes well at the end of this process there should be only one item left in the stack. This item is the initial state of the whole NFA

### 3. Match a string to the NFA
This algorithm is maintaining two set of arrays. The first array is the "current nodes"(initialized from the initial starting node) and the second one is the "next nodes",
<br />
A helper function is introduced here which adds a given node to a list of nodes and checks if the current node is a symbol node. If it is a symbol node then it calls recursively on the two edges of this node. Otherwise it returns the list. This helper function is used to create the "current nodes" and "next nodes" arrays. 
<br />
The input string is looped through character by character. Each character is matched against the items in the current node array and mached to each node's symbol, if the symbol matches the character, the next batch of nodes are loaded into the next nodes array starting from the current nodes first edge. Once the loop is done on the current nodes they are replaced by the next nodes array and the main loop goes to the next character.
<br />
<br />
After the main loop is finished a new loop goes through the remainder of the current nodes. If any of them is the same as the original accept node, the algorithm returns true, which indicates the regex is accepted for the given string otherwise it is rejected.


## Extension possibilities
As this is a really simple regular expression parser, there is a possibility to add a number of new features:
<br />
* Character escaping: ```\*```
* Character ranges: ```[a-zA-z]``` ```[0-9]```
* Quantifiers: ```{1}``` ```{1,}```
And so on...

## References
* https://swtch.com/~rsc/regexp/regexp1.html
### Infix to postfix
* Infix to postfix - http://csis.pace.edu/~wolf/CS122/infix-postfix.htm
* Shunting yard algorithm - https://en.wikipedia.org/wiki/Shunting-yard_algorithm
### NFA
* Nondeterministic finite automota - https://www.tutorialspoint.com/automata_theory/non_deterministic_finite_automaton.htm
* Thompson's Construction - https://stackoverflow.com/questions/11819185/steps-to-creating-an-nfa-from-a-regular-expression
### GO Lecture videos
* Shunting yard algorithm - https://web.microsoftstream.com/video/9d83a3f3-bc4f-4bda-95cc-b21c8e67675e
* Thompson's Construction - https://web.microsoftstream.com/video/68a288f5-4688-4b3a-980e-1fcd5dd2a53b
