package utils

import (
	"github.com/schollz/progressbar/v3"
)

type ProgressBar struct {
	bar *progressbar.ProgressBar
}

func InitBar(length int) *ProgressBar {
	return &ProgressBar{
		bar: progressbar.NewOptions(length,
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionShowCount(),
			progressbar.OptionSetWidth(50),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "[green]=[reset]",
				SaucerHead:    "[green]>[reset]",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			})),
	}
}

func (b *ProgressBar) Begin() {
	err := b.bar.RenderBlank()
	if err != nil {
		panic(err)
	}
}

func (b *ProgressBar) Add() {
	err := b.bar.Add(1)
	if err != nil {
		panic(err)
	}
}

func (b *ProgressBar) Finish() {
	err := b.bar.Finish()
	if err != nil {
		panic(err)
	}
}
