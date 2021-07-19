package filecoin

import "strconv"

func CalculateFilRewardFromAtoFil(reward string) float64 {
	convertedReward, _ := strconv.ParseFloat(reward, 64)
	divideBy5, _ := strconv.ParseFloat("5", 64)
	divideBy1e18, _ := strconv.ParseFloat("1e18", 64)
	convertedReward = convertedReward / divideBy5 / divideBy1e18
	return convertedReward
}
