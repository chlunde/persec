package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func nonDigits(str string) string {
	var buffer bytes.Buffer
	for _, ch := range str {
		if !unicode.IsDigit(ch) && !unicode.IsSpace(ch) {
			buffer.WriteRune(ch)
		}
	}
	return buffer.String()
}

func collectDigits(a string) (int, int) {
	var ret bytes.Buffer
	for _, ch := range a {
		if unicode.IsDigit(ch) {
			ret.WriteRune(ch)
		} else {
			break
		}
	}
	n, err := strconv.Atoi(ret.String())
	if err != nil {
		panic("digit-only string could not be parsed: " + ret.String())
	}
	return n, ret.Len()
}

func delta(a, b string, elapsed time.Duration, w io.Writer) {
	var ai, bi int
	for {
		for ai < len(a) && !unicode.IsDigit(rune(a[ai])) {
			w.Write([]byte{a[ai]})
			ai++
		}

		for bi < len(b) && !unicode.IsDigit(rune(b[bi])) {
			bi++
		}

		if ai == len(a) {
			break
		}

		numA, lenA := collectDigits(a[ai:])
		ai += lenA

		numB, lenB := collectDigits(b[bi:])
		bi += lenB

		// TODO: look at the next words, what's the unit?
		// TODO: colors
		// TODO: KB/MB
		fmt.Fprintf(w, "%.2f/s",
			float64(numB-numA)/(float64(elapsed)/float64(time.Second*1)))
	}
	fmt.Fprintln(w)
}

func PerSec(then, now []string, elapsed time.Duration, didPrint map[string]int, out io.Writer) {
	// TODO len(then) != len(now)
	for i, j := 0, 0; i < len(then); i++ {
		idA := nonDigits(then[i])
		idB := nonDigits(now[j])
		if idA == idB && (then[i] != now[j] || didPrint[idA] > 0) {
			delta(then[i], now[j], elapsed, out)
			if then[i] != now[j] {
				didPrint[idA] = 10
			} else {
				didPrint[idA]--
			}
		}
		j++
	}
}

func main() {
	// TODO: flags: duration, count, colors

	childCommand := os.Args[1:]
	var then []string
	//delay := time.Second * 1
	delay := time.Millisecond * 1500
	didPrint := make(map[string]int)
	for {
		cmd := exec.Command(childCommand[0], childCommand[1:]...)
		out, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}

		now := strings.Split(string(out), "\n")
		if then != nil {
			fmt.Println("\033[2J\033[H")
			PerSec(then, now, delay, didPrint, os.Stdout)
		}
		time.Sleep(delay)
		then = now
	}
}
