SIMPLE BPE LOGIC

BPE (Byte-Pair Encoding) adalah teknik kompresi data yang digunakan dalam pemrosesan bahasa alami. BPE bekerja dengan menggabungkan dua karakter yang paling sering muncul bersama-sama.
Proses ini dilakukan secara berulang hingga mencapai jumlah token yang diinginkan.

BPE ini juga dapat digunakan untuk membangun kamus kata yang lebih besar. Semisal, kita tidak memili kamus kata
pada dataset yang kita gunakan, dengan BPE kita dapat membuat kamus kata yang lebih besar.

``` go
package main

import (
	"fmt"
	"strings"
)

func main() {

	text := "Sometime is a some and time"

	text = strings.ToLower(text)

	words := strings.Fields(text)
	var (
		result      []string
		frequency   = make(map[string]int)
		maxLendWord = 15
	)

	for i := 2; i <= maxLendWord; i++ {

		for _, word := range words {
			for j := 0; j < len(word); j++ {
				if j+i > len(word) {
					continue
				}

				subWord := word[j : j+i]
				result = append(result, subWord)

				if v, ok := frequency[subWord]; ok {
					frequency[subWord] = v + 1
				} else {
					frequency[subWord] = 1
				}

			}
		}
	}

	fmt.Println(result)
	fmt.Println(frequency)

}
```
