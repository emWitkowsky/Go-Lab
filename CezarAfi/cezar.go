// Copyright by Michał Wojciech Witkowski index: 278862

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func caesarCipher(s string, shift int) string {
	encrypted := ""
	for _, r := range s {
		switch {
		case 'a' <= r && r <= 'z':
			encrypted += string('a' + (r-'a'+rune(shift))%26)
		case 'A' <= r && r <= 'Z':
			encrypted += string('A' + (r-'A'+rune(shift))%26)
		default:
			encrypted += string(r)
		}
	}
	return encrypted
}

func affineCipher(s string, a int, b int) string {
	encrypted := ""
	for _, r := range s {
		switch {
		case 'a' <= r && r <= 'z':
			encrypted += string('a' + (rune(a)*(r-'a')+rune(b))%26)
		case 'A' <= r && r <= 'Z':
			encrypted += string('A' + (rune(a)*(r-'A')+rune(b))%26)
		default:
			encrypted += string(r)
		}
	}
	return encrypted
}

func caesarEncrypt() {
	file, err := os.Open("key.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	key, _, err := reader.ReadRune()
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	textBytes, err := os.ReadFile("plain.txt")
	if err != nil {
		panic(err)
	}

	shift := key - '0'
	text := string(textBytes)
	encryptedText := caesarCipher(text, int(shift))

	writeOut := os.WriteFile("crypto.txt", []byte(encryptedText), 0644)
	if writeOut != nil {
		panic(writeOut)
	}
}

func caesarDecrypt() {
	file, err := os.Open("key.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	key, _, err := reader.ReadRune()
	if err != nil {
		panic(err)
	}

	textBytes, err := os.ReadFile("crypto.txt")
	if err != nil {
		panic(err)
	}

	shift := -(key - '0')
	text := string(textBytes)
	decryptedText := caesarCipher(text, int((shift+26)%26))

	writeOut := os.WriteFile("decrypt.txt", []byte(decryptedText), 0644)
	if writeOut != nil {
		panic(writeOut)
	}
}

func caesarCryptanalysisWithKnownText() {
	textBytes1, err1 := os.ReadFile("crypto.txt")
	if err1 != nil {
		panic(err1)
	}

	text1 := string(textBytes1)

	textBytes2, err := os.ReadFile("extra.txt")
	if err != nil {
		panic(err)
	}

	text2 := string(textBytes2)

	var shift int

	for i := range text2 {
		r := text1[i]
		r2 := text2[i]

		switch {
		case 'a' <= r && r <= 'z':
			shift = int((r - r2 + 26) % 26)
		case 'A' <= r && r <= 'Z':
			shift = int((r - r2 + 26) % 26)
		default:
			shift = int((r - r2 + 26) % 26)
		}
	}

	writeOut := os.WriteFile("key-found.txt", []byte(fmt.Sprintln(shift)), 0644)
	if writeOut != nil {
		panic(writeOut)
	}
}

func caesarCryptanalysisWithoutKnownText() {
	textBytes, err := os.ReadFile("crypto.txt")
	if err != nil {
		panic(err)
	}

	text := string(textBytes)

	var output string

	for shift := 1; shift < 26; shift++ {
		decryptedText := caesarCipher(text, 26-shift)
		output += decryptedText + "\n"
	}
	writeOut := os.WriteFile("decrypt.txt", []byte(output), 0644)
	if writeOut != nil {
		panic(writeOut)
	}
}

func NWD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func affineEncrypt() {
	file, err := os.Open("key.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	key, _, err := reader.ReadLine()
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	keys := strings.Fields(string(key))
	if len(keys) != 2 {
		fmt.Println("Wrong key")
	}

	// var aVal = []int{1, 3, 5, 7, 9, 11, 15, 17, 19, 21, 23, 25}
	shiftA, err := strconv.Atoi(keys[0])
	if err != nil {
		panic(err)
	}

	if NWD(shiftA, 26) != 1 {
		panic("Wrong key: a is not prime to 26")

	}

	shiftB, err := strconv.Atoi(keys[1])
	if err != nil {
		panic(err)
	}

	// fmt.Println(string(key))
	textBytes, err := os.ReadFile("plain.txt")
	if err != nil {
		panic(err)
	}

	text := string(textBytes)
	encryptedText := affineCipher(text, int(shiftA), int(shiftB))

	writeOut := os.WriteFile("crypto.txt", []byte(encryptedText), 0644)
	if writeOut != nil {
		panic(writeOut)
	}
}

func multiplicativeInverse(a, m int) int {
	a = a % m
	for x := 1; x < m; x++ {
		if (a*x)%m == 1 {
			return x
		}
	}
	return -1
}

func affineDecipher(s string, a int, b int) string {
	decrypted := ""
	aInv := multiplicativeInverse(a, 26)
	for _, r := range s {
		switch {
		case 'a' <= r && r <= 'z':
			val := (r - 'a') - rune(b)
			if val < 0 {
				val += 26
			}
			val = (rune(aInv) * val) % 26
			decrypted += string('a' + val)
		case 'A' <= r && r <= 'Z':
			val := (r - 'A') - rune(b)
			if val < 0 {
				val += 26
			}
			val = (rune(aInv) * val) % 26
			decrypted += string('A' + val)
		default:
			decrypted += string(r)
		}
	}
	return decrypted
}

func affineDecrypt() {
	file, err := os.Open("key.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	key, _, err := reader.ReadLine()
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// check if key has two values
	keys := strings.Fields(string(key))
	if len(keys) != 2 {
		fmt.Println("Wrong key")
	}

	shiftA, err := strconv.Atoi(keys[0])
	if err != nil {
		panic(err)
	}

	shiftB, err := strconv.Atoi(keys[1])
	if err != nil {
		panic(err)
	}

	textBytes, err := os.ReadFile("crypto.txt")
	if err != nil {
		panic(err)
	}

	text := string(textBytes)
	decryptedText := affineDecipher(text, shiftA, shiftB)

	writeOut := os.WriteFile("decrypt.txt", []byte(decryptedText), 0644)
	if writeOut != nil {
		panic(writeOut)
	}
}

func removeSpacesAndNewlines(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, " ", "")
	return s
}

func affineCryptanalysisWithKnownText() string {
	cryptoFile, err := os.Open("crypto.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer cryptoFile.Close()

	extraFile, err := os.Open("extra.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer extraFile.Close()

	scanner := bufio.NewScanner(cryptoFile)
	scanner.Scan()
	cryptogram := scanner.Text()
	cryptogram = removeSpacesAndNewlines(cryptogram)

	scanner = bufio.NewScanner(extraFile)
	scanner.Scan()
	extra := scanner.Text()
	extra = removeSpacesAndNewlines(extra)

	var key string
	for x := 0; x < len(extra)-1; x++ {
		y1 := int(cryptogram[x])
		y2 := int(cryptogram[x+1])
		x1 := int(extra[x])
		x2 := int(extra[x+1])

		if (x1 >= 65 && x1 <= 90) && (x2 >= 97 && x2 <= 122) || (x2 >= 65 && x2 <= 90) && (x1 >= 97 && x1 <= 122) {
			continue
		}
		if (x1 < 65 || x1 > 90) && (x1 < 97 || x1 > 122) {
			continue
		}
		if (x2 < 65 || x2 > 90) && (x2 < 97 || x2 > 122) {
			continue
		}

		for x1 > 25 {
			x1 -= 26
		}
		for x2 > 25 {
			x2 -= 26
		}
		for y1 > 25 {
			y1 -= 26
		}
		for y2 > 25 {
			y2 -= 26
		}

		keyA := float64(y1-y2) / float64(x1-x2)

		if int(keyA) == int(keyA) && int(keyA) > 0 {
			// fmt.Print(x1, x2, y1, y2, keyA, " ")
			keyB := 0
			support := (x1*int(keyA) + keyB) % 26
			for support != y1 {
				support++
				keyB++
				if support > 25 {
					support -= 26
				}
			}
			// fmt.Println(keyB)
			keyFile, err := os.Create("key-found.txt")
			if err != nil {
				log.Fatal(err)
			}
			defer keyFile.Close()
			keyFile.WriteString(fmt.Sprintf("%d %d", int(keyA), int(keyB)))
			// return fmt.Sprintf("b: %d, a: %d", keyB, int(keyA))
		}
	}
	return key
}

func affineCryptanalysisWithoutKnownText() {
	textBytes, err := os.ReadFile("crypto.txt")
	if err != nil {
		panic(err)
	}

	text := string(textBytes)

	var output string

	var ang = []int{1, 3, 5, 7, 9, 11, 15, 17, 19, 21, 23, 25}

	for _, shiftA := range ang {
		for shiftB := 0; shiftB < 26; shiftB++ {
			decryptedText := affineDecipher(text, shiftA, shiftB)
			output += decryptedText + "\n"
		}
	}

	writeOut := os.WriteFile("decrypt.txt", []byte(output), 0644)
	if writeOut != nil {
		panic(writeOut)
	}
}

func main() {
	cipher := flag.Bool("c", false, "Wybierz szyfr: cezar (-c) lub afiniczny (-a)")
	affine := flag.Bool("a", false, "Afiniczny szyfr (-a)")
	encrypt := flag.Bool("e", false, "Szyfrowanie (-e)")
	decrypt := flag.Bool("d", false, "Odszyfrowywanie (-d)")
	knownText := flag.Bool("j", false, "Kryptoanaliza z tekstem jawnym (-j)")
	onlyCryptogram := flag.Bool("k", false, "Kryptoanaliza wyłącznie w oparciu o kryptogram (-k)")
	flag.Parse()

	if *cipher {
		if *encrypt {
			caesarEncrypt()
		} else if *decrypt {
			caesarDecrypt()
		} else if *knownText {
			caesarCryptanalysisWithKnownText()
		} else if *onlyCryptogram {
			caesarCryptanalysisWithoutKnownText()
		} else {
			fmt.Println("Nieprawidłowa opcja")
		}
	} else if *affine {
		if *encrypt {
			affineEncrypt()
		} else if *decrypt {
			affineDecrypt()
		} else if *knownText {
			affineCryptanalysisWithKnownText()
		} else if *onlyCryptogram {
			affineCryptanalysisWithoutKnownText()
		} else {
			fmt.Println("Nieprawidłowa opcja")
		}
	} else {
		fmt.Println("Brak wybranego szyfru")
	}

}
