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
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

var TestFile = "/tmp/test_sf_document"
var Content = []byte{
	0x00, 0x10, 0x20, 0x30,
	0x40, 0x50, 0x60, 0x70,
	0x80, 0x90, 0x0a, 0x0b,
	0x0c, 0x0d, 0x0e, 0x0f,
}

func TestScatter(t *testing.T) {
	var err error

	err = ioutil.WriteFile(TestFile, Content, 0600)
	if err != nil {
		t.Error(err)
	}

	err = Scatter(TestFile)
	if err != nil {
		t.Error(err)
	}

	if _, err = os.Stat(TestFile + FILE_FST); os.IsNotExist(err) {
		t.Fail()
	}

	if _, err = os.Stat(TestFile + FILE_SND); os.IsNotExist(err) {
		t.Fail()
	}
}

func TestFormat(t *testing.T) {
	var err error
	var data []byte

	if err = os.Remove(TestFile); err != nil {
		t.Error(err)
	}

	if err = Format(TestFile+FILE_FST, TestFile+FILE_SND); err != nil {
		t.Error(err)
	}

	if _, err = os.Stat(TestFile); os.IsNotExist(err) {
		t.Fail()
	}

	if data, err = ioutil.ReadFile(TestFile); err != nil {
		t.Error(err)
	} else {
		n := bytes.Compare(Content, data)
		if n != 0 {
			t.Fail()
		}
	}
}
