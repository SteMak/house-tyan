package util

import (
	"errors"
	"strconv"
	"strings"
)

type Hex string

type RGB struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

func (hex Hex) GetClearHexValue() (Hex, error) {
	result := strings.TrimPrefix(string(hex), "#")
	if len(result) != 6 {
		return Hex("ffffff"), errors.New("wrong hex length")
	}
	return Hex(result), nil
}

func (rgb RGB) GetDecimalRGB() int {
	return int(rgb.Red)*256*256 + int(rgb.Green)*256 + int(rgb.Blue)
}

func Hex2RGB(hex Hex) (RGB, error) {
	var rgb RGB
	hex, err := hex.GetClearHexValue()
	if err != nil {
		return RGB{}, err
	}

	values, err := strconv.ParseUint(string(hex), 16, 32)
	if err != nil {
		return RGB{}, err
	}

	rgb = RGB{
		Red:   uint8(values >> 16),
		Green: uint8((values >> 8) & 0xFF),
		Blue:  uint8(values & 0xFF),
	}

	return rgb, nil
}
