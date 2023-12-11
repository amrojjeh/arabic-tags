package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/amrojjeh/arabic-tags/internal/speech"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		fmt.Println("Input error!")
		return
	}
	text = strings.TrimSpace(text)
	words, err := speech.Disambiguate(text)
	if err != nil {
		if errors.Is(err, speech.ErrRequest) {
			fmt.Printf("%v\n", text)
			return
		}
		fmt.Printf("Error! %v\n", err.Error())
		return
	}

	for _, w := range words {
		fmt.Printf("%v ", w)
	}
	fmt.Println()
}
