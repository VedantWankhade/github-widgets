package services

import (
	"fmt"
	"os"
)

func GetCard(cardName string) (*os.File, error) {
	var f *os.File
	var err error
	switch cardName {
	case "ping":
		f, err = pingCard()
	default:
		f, err = nil, fmt.Errorf("invalid card name")
	}
	return f, err
}

func pingCard() (*os.File, error) {
	f, err := os.Open("gsd.png")
	if err != nil {
		return nil, err
	}
	return f, nil
}
