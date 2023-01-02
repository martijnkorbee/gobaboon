package util

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func PrintInfo(msg string) {
	white := color.New(color.FgHiWhite).SprintFunc()
	fmt.Fprintf(os.Stdout, "%s\t%s\n", white("INF"), msg)
}

func PrintWarning(msg string) {
	yellow := color.New(color.FgHiYellow).SprintFunc()
	fmt.Fprintf(os.Stdout, "%s\t%s\n", yellow("WRN"), msg)
}

func PrintError(msg string, err error) {
	red := color.New(color.FgRed).SprintFunc()
	fmt.Fprintf(os.Stdout, "%s\t%s > error='%s'\n", red("ERR"), msg, err.Error())
}

func PrintFatal(msg string, err error) {
	red := color.New(color.FgHiRed).SprintFunc()
	fmt.Fprintf(os.Stdout, "%s\t%s > error='%s'\n", red("FTL"), msg, err.Error())

	os.Exit(1)
}

func PrintSuccess(msg string) {
	green := color.New(color.FgHiGreen).SprintFunc()
	fmt.Fprintf(os.Stdout, "%s\t%s\n", green("OK"), msg)
}

func PrintResult(msg string, v ...any) {
	green := color.New(color.FgHiGreen).SprintFunc()
	blue := color.New(color.FgHiBlue).SprintFunc()
	fmt.Fprintf(os.Stdout, "%s\t%s: %v\n", green("OK"), blue(msg), v)
}
