// this file content how to create simple tokenize raw data
// step for llm is
// raw-text --> tokenized text --> token IDs --> token embeding --> tranformer --> post processing step --> output
// this file cover section
// raw-text --> tokenized text
package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/born2ngopi/llm/tokenize/bpe"
)

func main() {

	// open raw text file and read it to string
	buf, err := os.ReadFile("raw.txt")
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
		return
	}

	var dataset = string(buf)

	tokenizer := NewToken()

	// first we need to tokenize the text
	tokenizedText := tokenizer.TokenizedText(dataset)

	// next we need to convert token to token ID
	// token ID is unique number for each token
	tokenizer.SetTokenId(tokenizedText)

	// testing encode and decode
	var texts = []string{"Hello, do you like tea?", "In the sunlit terraces of the palace."}
	text := strings.Join(texts, " <|endoftext|> ")

	tokenIds := tokenizer.Encode(text)
	fmt.Println("Token IDs: ", tokenIds)
	t := tokenizer.Decode(tokenIds)
	fmt.Println("Decoded: ", t)

}

type Token struct {
	tokenId      map[string]int
	convertToken map[int]string
	reEncode     *regexp.Regexp
	reDecode     *regexp.Regexp
}

func NewToken() *Token {
	return &Token{
		tokenId:      make(map[string]int),
		convertToken: make(map[int]string),
		reEncode:     regexp.MustCompile(`([,.:;?_!"()\']|--|\s)`),
		reDecode:     regexp.MustCompile(`\s+([,.?!"()\'])`),
	}
}

func (t *Token) Encode(text string) []int {

	splits := t.reEncode.Split(text, -1)
	matches := t.reEncode.FindAllString(text, -1)

	var tokenized []string
	for i := 0; i < len(splits); i++ {
		if splits[i] != "" {
			tokenized = append(tokenized, splits[i])
		}
		if i < len(matches) && matches[i] != " " {
			tokenized = append(tokenized, matches[i])
		}
	}

	var result []int
	for _, token := range tokenized {
		if id, ok := t.tokenId[token]; ok {
			result = append(result, id)
			continue
		} else {
			// if token not found in token ID we use UNK token
			result = append(result, t.tokenId["<|unk|>"])
			continue
		}
	}

	return result
}

func (t *Token) Decode(Ids []int) string {
	var s []string
	for _, id := range Ids {
		s = append(s, t.convertToken[id])
	}

	txt := strings.Join(s, " ")

	txt = t.reDecode.ReplaceAllString(txt, "$1")

	return txt

}

func (t *Token) SetTokenId(tokenizedText []string) {
	// before we convert token to token ID we need to sort the token
	// so expected result is example
	// [! , -- -- . . a am Hello I ini is keren LLM mantap model my name transformer]
	// and the token ID is
	// [0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15]
	sort.Slice(tokenizedText, func(i, j int) bool {
		return strings.ToLower(tokenizedText[i]) < strings.ToLower(tokenizedText[j])
	})

	for i, token := range tokenizedText {
		t.tokenId[token] = i
		t.convertToken[i] = token
	}
}

// tokenizedText function to split the text into tokens
// example input: "Hello, my name is LLM. I am a transformer model. -- mantap! ini-- keren"
// example output: ["Hello", ",", "my", "name", "is", "LLM", ".", "I", "am", "a", "transformer", "model", ".", "--", "mantap", "!", "ini", "--", "keren"]
func (t *Token) TokenizedText(dataset string) []string {

	// Split the text based on the regular expression
	splits := t.reEncode.Split(dataset, -1)
	matches := t.reEncode.FindAllString(dataset, -1)

	// Combine the splits and matches
	var result []string
	uniqueItems := make(map[string]bool)
	for i := 0; i < len(splits); i++ {
		if splits[i] != "" {
			item := splits[i]
			if !uniqueItems[item] {
				result = append(result, splits[i])
				uniqueItems[item] = true
			}
		}
		if i < len(matches) && matches[i] != " " {
			item := matches[i]
			if !uniqueItems[item] {
				result = append(result, matches[i])
				uniqueItems[item] = true
			}
		}
	}

	_, frequencies := bpe.Tokenize(dataset)
	minFreq := 2

	for word, freq := range frequencies {
		if freq >= minFreq {
			if !uniqueItems[word] {
				// append only new items
				result = append(result, word)
				uniqueItems[word] = true
			}
		}
	}

	// Some tokenizers use special tokens to help the LLM with additional context
	// example:
	// - [BOS] (beginning of sequence)
	// - [EOS] (end of sequence)
	// - [PAD] (padding)
	// - [UNK] (unknown token)
	// - Note that GPT-2 does not need any of these tokens mentioned above but only uses an <|endoftext|> token to reduce complexity
	// - The <|endoftext|> is analogous to the [EOS] token mentioned above
	// - GPT also uses the <|endoftext|> for padding (since we typically use a mask when training on batched inputs, we would not attend padded tokens anyways, so it does not matter what these tokens are)
	result = append(result, "<|endoftext|>", "<|unk|>")

	return result
}
