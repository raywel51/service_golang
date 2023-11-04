package utils

import (
	"fmt"
	"math/rand"
)

func GenerateThaiLicensePlate() string {
	// Define the Thai alphabet characters that can appear on license plates
	thaiAlphabet := "กขคฆงจฉชซฌญฎฏฐฑฒณดตถทธนบปผฝพฟภมยรลวสหฬอฮ"

	// Generate a random numeric part (4 digits)
	numericPart := fmt.Sprintf("%04d", rand.Intn(10000))

	// Generate a random prefix (2 Thai alphabet characters)
	prefix := ""
	for i := 0; i < 2; i++ {
		prefix += string(thaiAlphabet[rand.Intn(len(thaiAlphabet))])
	}

	// Generate a random numeric prefix (1 digit or no prefix)
	var numericPrefix string
	if rand.Float64() < 0.5 {
		numericPrefix = ""
	} else {
		numericPrefix = fmt.Sprintf("%d", rand.Intn(10))
	}

	// Combine the numeric prefix, prefix, numeric part
	licensePlate := numericPrefix + prefix + numericPart

	return licensePlate
}
