// Copyright (c) 2020 Alexandru Catrina <alex@codeissues.net>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package main

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// Buffer size
	BUFF_8BIT_24 = 3
	BUFF_6BIT_24 = 4

	// Filename suffix
	FILE_FST = ".sfa"
	FILE_SND = ".sfb"
)

// Scatter the content of a file by convering chunks of 24-bit information from
// raw binary to base64 representation. The original content will be saved into
// two different files.
func Scatter(filepath string) (err error) {
	var fd, sfa, sfb *os.File
	var buffa, buffb *bufio.Writer

	if fd, err = os.Open(filepath); err != nil {
		return err
	}
	defer fd.Close()

	if sfa, err = os.Create(filepath + FILE_FST); err != nil {
		return err
	}
	buffa = bufio.NewWriter(sfa)
	defer sfa.Close()

	if sfb, err = os.Create(filepath + FILE_SND); err != nil {
		return err
	}
	buffb = bufio.NewWriter(sfb)
	defer sfb.Close()

	for {
		buf8bit := make([]byte, BUFF_8BIT_24)
		buf6bit := make([]byte, BUFF_6BIT_24)

		var bytez int
		if bytez, err = fd.Read(buf8bit); err != nil {
			break
		} else {
			base64.StdEncoding.Encode(buf6bit, buf8bit[:bytez])
		}

		for i := 0; i < BUFF_6BIT_24; i += 2 {
			buffa.WriteByte(buf6bit[i])
			buffb.WriteByte(buf6bit[i+1])
		}
	}

	buffa.Flush()
	buffb.Flush()

	return nil
}

// Format the content of the original file by formatting back chunks of 24-bit
// information from base64 representation to binary. It is impossible to probe
// the provided scattered files and validate the compatibility with each other
// since there is no way to know if these files are the result of the original
// file.
func Format(firstFile, secondFile string) (err error) {
	var fd, sfa, sfb *os.File
	var buffd *bufio.Writer

	if sfa, err = os.Open(firstFile); err != nil {
		return err
	}
	defer sfa.Close()

	if sfb, err = os.Open(secondFile); err != nil {
		return err
	}
	defer sfb.Close()

	filename := strings.TrimSuffix(firstFile, FILE_FST)
	if _, err := os.Stat(filename); err == nil {
		return errors.New("File already exists")
	}

	if fd, err = os.Create(filename); err != nil {
		return err
	}
	buffd = bufio.NewWriter(fd)
	defer fd.Close()

loop:
	for {
		buf8bit := make([]byte, BUFF_8BIT_24)
		buf6bit := make([]byte, BUFF_6BIT_24)

		var tmp []byte
		for i := 0; i < BUFF_6BIT_24; i += 2 {
			tmp = make([]byte, 1)
			if _, err = sfa.Read(tmp); err != nil {
				break loop
			} else {
				buf6bit[i] = tmp[0]
			}

			tmp = make([]byte, 1)
			if _, err = sfb.Read(tmp); err != nil {
				break loop
			} else {
				buf6bit[i+1] = tmp[0]
			}
		}

		n, err := base64.StdEncoding.Decode(buf8bit, buf6bit)
		if err != nil {
			return err
		}

		if n < BUFF_8BIT_24 {
			buffd.Write(buf8bit[:n])
		} else {
			buffd.Write(buf8bit)
		}
	}

	buffd.Flush()

	return nil
}

// Launch program and scan stdin and command line arguments.
func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Cannot get working directory: %s\n", err)
		os.Exit(1)
	}

	finfo, err := os.Stdin.Stat()
	if err != nil {
		fmt.Printf("Error reading stdin: %s\n", err)
		os.Exit(1)
	}

	if os.ModeNamedPipe&finfo.Mode() == 0 {
		pairs := os.Args[1:]

		if len(pairs) == 0 {
			fmt.Printf("Usage: %s file.a file.b\n", os.Args[0])
			os.Exit(0)
		}

		if len(pairs)%2 != 0 {
			fmt.Println("Provided an odd number of arguments")
			os.Exit(1)
		}

		for i := 0; i < len(pairs); i += 2 {
			fst, snd := pairs[i], pairs[i+1]
			if err := Format(fst, snd); err != nil {
				msg := "Error formatting %s and %s: %s\n"
				fmt.Printf(msg, fst, snd, err)
			}
		}

		return
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fp, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		file := filepath.Join(cwd, strings.TrimSpace(fp))
		if err := Scatter(file); err != nil {
			fmt.Printf("Error scattering %s: %s\n", fp, err)
		}
	}
}
