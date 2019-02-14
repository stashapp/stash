package utils

import (
	"encoding/base64"
	"fmt"
	"regexp"
)

func ProcessBase64Image(imageString string) (string, []byte, error) {
	if imageString == "" {
		return "", nil, fmt.Errorf("empty image string")
	}

	regex := regexp.MustCompile(`^data:.+\/(.+);base64,(.*)$`)
	matches := regex.FindStringSubmatch(imageString)
	var encodedString string
	if len(matches) > 2 {
		encodedString = regex.FindStringSubmatch(imageString)[2]
	} else {
		encodedString = imageString
	}
	imageData, err := GetDataFromBase64String(encodedString)
	if err != nil {
		return "", nil, err
	}

	return MD5FromBytes(imageData), imageData, nil
}

func GetDataFromBase64String(encodedString string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encodedString)
}

func GetBase64StringFromData(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)

	// Really slow
	//result = regexp.MustCompile(`(.{60})`).ReplaceAllString(result, "$1\n")
	//if result[len(result)-1:] != "\n" {
	//	result += "\n"
	//}
	//return result
}
