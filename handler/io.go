package handler

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type IOHandler interface {
	Println(a ...interface{})
	Printf(format string, a ...interface{})
	ReadLine(scanner *bufio.Scanner, prompt string) string
}

type IOHandlerImpl struct {
	Reader io.Reader
	Writer io.Writer
}

func NewIOHandler() *IOHandlerImpl {
	return &IOHandlerImpl{
		Reader: os.Stdin,
		Writer: os.Stdout,
	}
}

func (ioh *IOHandlerImpl) Println(a ...interface{}) {
	fmt.Fprintln(ioh.Writer, a...)
}

func (ioh *IOHandlerImpl) Printf(format string, a ...interface{}) {
	fmt.Fprintf(ioh.Writer, format, a...)
}

func (ioh *IOHandlerImpl) ReadLine(scanner *bufio.Scanner, prompt string) string {
	ioh.Printf("%s", prompt)
	scanner.Scan()
	return scanner.Text()
}
