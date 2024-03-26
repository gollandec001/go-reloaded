package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func replaceHex(input string) string {
	re := regexp.MustCompile(`^\(hex\)(\s|$)`)
	if re.MatchString(input) {
		return re.ReplaceAllString(input, "")
	}
	reHex := regexp.MustCompile(`(\b[0-9A-Fa-f]+)\s*\(hex\)`)
	return reHex.ReplaceAllStringFunc(input, func(match string) string {
		hexNumber := reHex.FindStringSubmatch(match)[1]
		decimalNumber, _ := strconv.ParseInt(hexNumber, 16, 64)
		return strconv.FormatInt(decimalNumber, 10)
	})
}

func replaceBin(input string) string {
	re := regexp.MustCompile(`^\(bin\)(\s|$)`)
	if re.MatchString(input) {
		return re.ReplaceAllString(input, "")
	}
	reBin := regexp.MustCompile(`(\b[01]+)\s*\(bin\)`)
	return reBin.ReplaceAllStringFunc(input, func(match string) string {
		binNumber := reBin.FindStringSubmatch(match)[1]
		decimalNumber, _ := strconv.ParseInt(binNumber, 2, 64)
		return strconv.FormatInt(decimalNumber, 10)
	})
}

func modifyCaseSensitive(input string) string {
	re := regexp.MustCompile(`^\((up|low|cap)(?:, (\d+))?\)$`)
	if re.MatchString(input) {
		return re.ReplaceAllString(input, "")
	}
	modifiers := make([]int, 0)
	reModifiers := regexp.MustCompile(`(up|low|cap)(?:, (\d+))?`)
	reModifiers.ReplaceAllStringFunc(input, func(match string) string {
		count, _ := strconv.Atoi(reModifiers.FindStringSubmatch(match)[2])
		modifiers = append(modifiers, count)
		return strconv.Itoa(count)
	})
	processed := input
	for {
		modified := processed
		for i := 0; i < len(modifiers); i++ {
			if modifiers[i] >= 0 {
				var b int
				if modifiers[i] > 0 {
					b = modifiers[i] - 1
				}
				reMod := regexp.MustCompile(fmt.Sprintf(`((?:\S*\b[\w']+\b[^\w']*){0,%d}\b[\w']+\b)\s*\((up|low|cap)(?:, (\d+))?\)`, b))
				processed = reMod.ReplaceAllStringFunc(processed, func(match string) string {
					word := reMod.FindStringSubmatch(match)[1]
					mod := reMod.FindStringSubmatch(match)[2]
					count, _ := strconv.Atoi(reMod.FindStringSubmatch(match)[3])
					var result string
					if count == modifiers[i] {
						switch mod {
						case "up":
							result = strings.ToUpper(word)
						case "low":
							result = strings.ToLower(word)
						case "cap":
							result = replaceCap(word)
						}
						return result
					}
					return match
				})
			}
		}
		if modified == processed {
			break
		}
	}
	return processed
}

func replaceCap(input string) string {
	words := strings.Fields(input)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}
	result := strings.Join(words, " ")
	return result
}

	func fixArticles(input string) string {
	re := regexp.MustCompile(`\b([aA])\s+(\b\w+\b)?`)
	return re.ReplaceAllStringFunc(input, func(match string) string {
		firstChar := re.FindStringSubmatch(match)[1]
		nextWord := re.FindStringSubmatch(match)[2]
		if len(nextWord) > 0 && strings.Contains("aeiouhAEIOUH", string(nextWord[0])) {
			return firstChar + "n" + " " + nextWord
		}
		return firstChar + " " + nextWord
	})
}

func formatPunctuation(input string) string {
	rePunctuation := regexp.MustCompile(`(?: [.,!?:;])`)
	input = rePunctuation.ReplaceAllStringFunc(input, func(match string) string {
		punctuation := rePunctuation.FindStringSubmatch(match)[0]
		return strings.TrimSpace(punctuation)
	})
	reComma := regexp.MustCompile(`,([^ ])`)
	input = reComma.ReplaceAllString(input, ", $1")
	reSingleQuotes := regexp.MustCompile(`' \s*([^']+)\s* '`)
	input = reSingleQuotes.ReplaceAllString(input, "'$1'")
	reSingleQuotes = regexp.MustCompile(`'\s*([^']+)\s* '`)
	input = reSingleQuotes.ReplaceAllString(input, "'$1'")
	reSingleQuotes = regexp.MustCompile(`' \s*([^']+)\s*'`)
	input = reSingleQuotes.ReplaceAllString(input, "'$1'")
	return input
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go input.txt output.txt")
		os.Exit(1)
	}

	inputFileName := os.Args[1]
	outputFileName := os.Args[2]

	inputFile, err := os.Open(inputFileName)
	if err != nil {
		fmt.Println("Error opening input file:", err)
		os.Exit(1)
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	var modifiedLines []string
	for _, text := range lines {
		text = replaceHex(text)
		text = replaceBin(text)
		text = modifyCaseSensitive(text)
		text = fixArticles(text)
		text = formatPunctuation(text)
		modifiedLines = append(modifiedLines, text)
	}

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	for _, modifiedLine := range modifiedLines {
		fmt.Fprintln(outputFile, modifiedLine)
	}
}
