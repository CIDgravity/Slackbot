package utils

import "twinQuasarAppV2/backend/customTypes"

func CheckIfExist(element customTypes.UniqueBlockMessage, array []customTypes.UniqueBlockMessage) (result bool, index int) {
	result = false
	index = -1

	for i, e := range array {
		if e.Nonce == element.Nonce && e.From == element.From {
			result = true
			index = i
			break
		}
	}

	return result, index
}
