package main

import (
	"bufio"
	"strings"
	"fmt"
	"bytes"
	"io"
)

func main() {
	
	c := strings.NewReader(cipher)

	// m,t  := LetterCount(c)

	// for k, v := range m {
	// 	fmt.Printf("%v:%v <- %v %%\n", k,v, (float64(v) / float64(t)),  )
	// }

	scanner := bufio.NewScanner(c)
	scanner.Split(bufio.ScanRunes)

	for scanner.Scan() {
		l := scanner.Text()		
		// if bytes.Runes([]byte(word))[0] == 13 {continue}
		// if bytes.Runes([]byte(word))[0] == 32 {continue}		
		switch l {
			case "b" : fmt.Printf("t")
			case "p" : fmt.Printf("h")
			case "r" : fmt.Printf("e")
			case "m" : fmt.Printf("a")
			case "j" : fmt.Printf("o")			
			case "w" : fmt.Printf("i")
			case "i" : fmt.Printf("s")
			case "k" : fmt.Printf("n")
			case "l" : fmt.Printf("b")
			case "n" : fmt.Printf("u")
			case "v" : fmt.Printf("c")
			case "t" : fmt.Printf("y")
			case "q" : fmt.Printf("k")
			case "x" : fmt.Printf("f")
			case "u" : fmt.Printf("r")
			case "o" : fmt.Printf("g")
			case "y" : fmt.Printf("m")
			case "e" : fmt.Printf("v")
			case "s" : fmt.Printf("p")
			case "a" : fmt.Printf("x")
			case "h" : fmt.Printf("l")
			case "c" : fmt.Printf("w")
			case "g" : fmt.Printf("z")
			
				default: fmt.Printf("%v", l)
		}
		
	}

}



func LetterCount(rdr io.Reader) (map[string]int, int) {
	var total int
	counts := map[string]int{}
	scanner := bufio.NewScanner(rdr)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		word := scanner.Text()		
		if bytes.Runes([]byte(word))[0] == 13 {continue}
		if bytes.Runes([]byte(word))[0] == 32 {continue}
		word = strings.ToLower(word)
		counts[word]++
		total++
	}
	return counts, total
}

var cipher = `lrvmnir bpr sumvbwvr jx bpr lmiwv yjeryrkbi jx qmbm wi 
bpr xjvni mkd ymibrut jx irhx wi bpr riirkvr jx 
ymbinlmtmipw utn qmumbr dj w ipmhh but bj rhnvwdmbr bpr
yjeryrkbi jx bpr qmbm mvvjudwko bj yt wkbrusurbmbwjk
lmird jk xjubt trmui jx ibndt

wb wi kjb mk rmit bmiq bj rashmwk rmvp yjeryrkb mkd wbi 
iwokwxwvmkvr mkd ijyr ynib urymwk nkrashmwkrd bj ower m
vjyshrbr rashmkmbwjk jkr cjnhd pmer bj lr fnmhwxwrd mkd 
wkiswurd bj invp mk rabrkb bpmb pr vjnhd urmvp bpr ibmbr
jx rkhwopbrkrd ywkd vmsmlhr jx urvjokwgwko ijnkdhrii
ijnkd mkd ipmsrhrii ipmsr w dj kjb drry ytirhx bpr xwkmh
mnbpjuwbt lnb yt rasruwrkvr cwbp qmbm pmi hrxb kj djnlb
bpmb bpr xjhhjcwko wi bpr sujsru msshwvmbwjk mkd
wkbrusurbmbwjk w jxxru yt bprjuwri wk bpr pjsr bpmb bpr
riirkvr jx jqwkmcmk qmumbr cwhh urymwk wkbmvb`


func Ceasar(cipherText string) {


	//Will print all versions of the possible keys as the key space is only 26
	for i := 0; i < 26; i++{

		scanner := bufio.NewScanner(strings.NewReader(cipherText))

		scanner.Split(bufio.ScanRunes)

		for scanner.Scan() {
			v := scanner.Text()

			rns := bytes.Runes([]byte(v))

			if rns[0] == 32{fmt.Printf(" "); continue}
			if rns[0] == 13{fmt.Printf("\n"); continue}

			rNum := rune(rns[0]-97)
			modded := (rNum + rune(i)) % 26
			clear := modded + 97
			// fmt.Printf("%v == %v\t", clear, rns[0])
			fmt.Printf("%v", string(clear))
			// fmt.Println(string(rune(((rns[0])-65)%26)+65) )
		}
		fmt.Println()
	}
}