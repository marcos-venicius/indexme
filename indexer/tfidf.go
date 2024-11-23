package indexer

func getFrequency(tokens []string) map[string]int {
	frequency := make(map[string]int)

	for _, token := range tokens {
		if count, ok := frequency[token]; ok {
			frequency[token] = count + 1
		} else {
			frequency[token] = 1
		}
	}

	return frequency
}

// updateFrequencyTable returns the updated tokens based on the global one
func updatedFrequencyTableTokens(folderTermsFrequency map[string]int, documentTermsFrequency map[string]int) map[string]int {
	output := make(map[string]int)

	for token, count := range documentTermsFrequency {
		if currentCount, ok := folderTermsFrequency[token]; ok {
			if outputCount, ok := output[token]; ok {
				output[token] = currentCount + count + outputCount
			} else {
				output[token] = currentCount + count
			}
		} else {
			if outputCount, ok := output[token]; ok {
				output[token] = currentCount + count + outputCount
			} else {
				output[token] = currentCount + count
			}
		}
	}

	return output
}
