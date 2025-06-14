package handler

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	Bold        = "\033[1m"
)

type IOHandler interface {
	Println(a ...interface{})
	Printf(format string, a ...interface{})
	PrintlnColored(color string, a ...interface{})
	PrintfColored(color string, format string, a ...interface{})
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

func (ioh *IOHandlerImpl) PrintlnColored(color string, a ...interface{}) {
	fmt.Fprint(ioh.Writer, color)
	ioh.Println(a...)
	fmt.Fprint(ioh.Writer, ColorReset)
}

func (ioh *IOHandlerImpl) PrintfColored(color string, format string, a ...interface{}) {
	fmt.Fprint(ioh.Writer, color)
	ioh.Printf(format, a...)
	fmt.Fprint(ioh.Writer, ColorReset)
}

func (ioh *IOHandlerImpl) ReadLine(scanner *bufio.Scanner, prompt string) string {
	ioh.Printf("%s", prompt)
	scanner.Scan()
	return scanner.Text()
}
