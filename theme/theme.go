package theme

import (
	color "github.com/fatih/color"
)

var title = color.New(color.FgWhite).Add(color.Bold, color.Underline)

func Title(text string) {
	title.Println(text)
}

func Warning(text string) {
	color.Yellow(text)
}

func Error(text string) {
	color.Red(text)
}

func Info(text string) {
	color.HiBlack(text)
}

func Success(text string) {
	color.Green(text)
}
