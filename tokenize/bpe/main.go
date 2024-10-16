// this fil is used to tokenize the text using BPE algorithm
// it's allow model to breakdown the text into into smaller subword units or
// individual characters
package main

type bpe struct {
}

func NewEncoder() *bpe {
	return &bpe{}
}
