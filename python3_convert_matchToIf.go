package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) != 2 {
		log.Fatal("usage: pythonMatchToIf <python file>")
	}
	filename := os.Args[1]

	var fileContents []string

	// read file into fileContents array
	if f, err := os.Open(filename); err != nil {
		log.Fatalf(err.Error())
	} else {
		o, err := io.ReadAll(f)
		f.Close()
		if err != nil {
			log.Fatalf(err.Error())
		}
		fileContents = strings.Split(string(o), "\n")
	}

	// function to log error with line information before exiting
	logError := func(i int, s string) {
		log.Fatal(fmt.Sprintf("%s: line %d\n\t%s", s, i+1, strings.TrimSpace(fileContents[i])))
	}

	// read line from fileContent
	readLine := func(i int) (string, string) {
		s := fileContents[i]
		l := strings.TrimSpace(s)
		indent := s[:len(s)-len(l)]
		return l, indent
	}

	for i := 0; i < len(fileContents)-1; i++ {
		l, indent := readLine(i)

		if strings.HasPrefix(l, "match") {
			matchBlock := true
			matchIndent := indent

			// parse match statement
			flds := strings.Fields(l)
			if len(flds) < 2 {
				logError(i, "malformed match statement")
			}

			// operator will be used to create the if statement. trailing colon
			// is chopped off
			operator := strings.TrimSuffix(flds[1], ":")
			if len(operator) == len(flds[1]) {
				logError(i, "malformed match statement")
			}

			// get indent for first case statement. this will be used to detect
			// the end of individual case blocks
			i++
			_, caseIndent := readLine(i)

			// an if statement will be output in place of the first case
			// statement. for subsequent cases, an elif statement will be output
			firstCase := true

			for ; matchBlock && i < len(fileContents); i++ {
				l, currentIndent := readLine(i)

				// only other case statements should be at this indentation level
				if currentIndent == caseIndent {
					if !strings.HasPrefix(l, "case ") {
						logError(i, "expecting a case statement")
					}

					// operand will be used to create the if statement
					operand, statement, ok := strings.Cut(l[len("case "):], ":")
					if !ok {
						logError(i, "malformed case statement")
					}

					// output if/elif statement in place of case statement
					if firstCase {
						fmt.Printf("%sif %s == %s:\n", matchIndent, operator, operand)
						firstCase = false
					} else {
						fmt.Printf("%selif %s == %s:\n", matchIndent, operator, operand)
					}

					// output first line of if statement if case line includes it
					if len(statement) > 0 {
						fmt.Printf("%s    %s\n", matchIndent, statement)
					}
				} else if len(currentIndent) <= len(matchIndent) {
					// end of match block
					fmt.Println(fileContents[i])
					matchBlock = false
				} else {
					// print out lines inside case block
					fmt.Printf("%s    %s\n", matchIndent, fileContents[i])
				}
			}
		}

		fmt.Println(fileContents[i])
	}
}
