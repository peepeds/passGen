package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

func generateCharset(useUpper, useLower, useNumber, useSymbol bool) []byte {
	
	var charset []byte
	ranges := []struct{ from, to int }{}

	if useUpper {
		ranges = append(ranges, struct{ from, to int }{65, 90}) // A-Z
	}
	if useLower {
		ranges = append(ranges, struct{ from, to int }{97, 122}) // a-z
	}
	if useNumber {
		ranges = append(ranges, struct{ from, to int }{48, 57}) // 0-9
	}
	if useSymbol {
		ranges = append(ranges,
			struct{ from, to int }{33, 47},
			struct{ from, to int }{58, 64},
			struct{ from, to int }{91, 96},
			struct{ from, to int }{123, 126},
		)
	}

	for _, r := range ranges {
		for i := r.from; i <= r.to; i++ {
			charset = append(charset, byte(i))
		}
	}

	return charset
}

func generatePassword(length int, charset []byte) (string, error) {
	if length < 8 {
		return "", errors.New("minimum length is 8")
	}
	if len(charset) == 0 {
		return "", errors.New("illegal choices")
	}

	password := make([]byte, length)
	for i := range password {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))

		if err != nil {
			return "", err
		}
		password[i] = charset[idx.Int64()]
	}

	return string(password), nil
}

func main() {
	var length int
	var upper, lower, number, symbol bool

	pflag.IntVarP(&length, "length", "L", 8, "Password Length")
	pflag.BoolVarP(&upper, "upper", "u", false, "Using Upper case (A-Z)")
	pflag.BoolVarP(&lower, "lower", "l", false, "Using Lower case (a-z)")
	pflag.BoolVarP(&number, "number", "n", false, "Using Number (0-9)")
	pflag.BoolVarP(&symbol, "symbol", "s", false, "Using Symbol (!@#$...)")

	// Help message custom
	pflag.Usage = func() {
		var b strings.Builder
		fmt.Fprintln(&b, "Password Generator CLI")
		fmt.Fprintln(&b, "Create random password based on ASCII character.\n")
		fmt.Fprintln(&b, "Examples:")
		fmt.Fprintln(&b, "  passgen -L 16 -u -l -n -s")
		fmt.Fprintln(&b, "  passgen --length 20 --upper --lower\n")
		fmt.Fprintln(&b, "Opsi:")
		fmt.Fprint(os.Stderr, b.String())
		pflag.PrintDefaults()
	}

	pflag.Parse()

	charset := generateCharset(upper, lower, number, symbol)

	password, err := generatePassword(length, charset)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	fmt.Println("Generated password:", password)
}
