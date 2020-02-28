package main

import (
	"bytes"
	"go/format"
	"strings"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func printDiff(expected, actual string) string {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(expected, actual, false)

	return dmp.DiffPrettyText(diffs)
}

func Test_generateInterface(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))

	err := generateInterface(buf, "Something", "Aaa", "Bbb", "Ccc")
	if err != nil {
		t.Errorf(err.Error())
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		t.Errorf(err.Error())
	}

	trimmed := strings.TrimSpace(string(formatted))

	expected := strings.TrimSpace(`
		type EnumSomething interface {
			something()
			String() string
		}

		type privateSomethingEnumType string

		func (s privateSomethingEnumType) something() {
			// unreachable method
		}

		func (s privateSomethingEnumType) String() string {
			return string(s)
		}
	`)

	if expected != trimmed {
		t.Errorf("incorrect generated output from generateInterface: \n %s\n", printDiff(expected, trimmed))
	}
}

func Test_generateConstants(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))

	err := generateConstants(buf, "Something", "Aaa", "Bbb", "Cc_xxxx_zzz")
	if err != nil {
		t.Errorf(err.Error())
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		t.Errorf(err.Error())
	}

	trimmed := strings.TrimSpace(string(formatted))

	expected := strings.TrimSpace(`
         const (
			strSomethingAaa       = "Aaa"
			strSomethingBbb       = "Bbb"
			strSomethingCcXxxxZzz = "Cc_xxxx_zzz"

			EnumSomethingAaa       = privateSomethingEnumType(strSomethingAaa)
			EnumSomethingBbb       = privateSomethingEnumType(strSomethingBbb)
			EnumSomethingCcXxxxZzz = privateSomethingEnumType(strSomethingCcXxxxZzz)
		)
	`)

	if expected != trimmed {
		t.Errorf("incorrect generated output from generateConstants: \n %s\n", printDiff(expected, trimmed))
	}
}

func Test_generateFuncFromStr(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))

	err := generateFuncFromStr(buf, "Something", "Aaa", "Bbb", "Cc_xxxx_zzz")
	if err != nil {
		t.Errorf(err.Error())
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		t.Errorf(err.Error())
	}

	trimmed := strings.TrimSpace(string(formatted))

	expected := strings.TrimSpace(`
         func FromString(s string) (EnumSomething, error) {
			switch s {
			case strSomethingAaa:
				return EnumSomethingAaa, nil
			case strSomethingBbb:
				return EnumSomethingBbb, nil
			case strSomethingCcXxxxZzz:
				return EnumSomethingCcXxxxZzz, nil
			}
			return nil, errors.New("unknown EnumSomething " + s)
		}
	`)

	if expected != trimmed {
		t.Errorf("incorrect generated output from generateConstants: \n %s\n", printDiff(expected, trimmed))
	}
}

func Test_generate(t *testing.T) {
	bb, err := generate("main", "Something", "Aaa", "Bbb", "Cc_xxxx_zzz")
	if err != nil {
		t.Errorf(err.Error())
	}

	generated := string(bb)

	println(string(generated))

	bb, err = format.Source([]byte(`
package mainenum

// Code generated by \"strenum\"; DO NOT EDIT.
// strenum <dir> main Something,Aaa,Bbb,Cc_xxxx_zzz

import "errors"

const (
	strMainSomething = "Something"
	strMainAaa       = "Aaa"
	strMainBbb       = "Bbb"
	strMainCcXxxxZzz = "Cc_xxxx_zzz"

	EnumMainSomething = privateMainEnumType(strMainSomething)
	EnumMainAaa       = privateMainEnumType(strMainAaa)
	EnumMainBbb       = privateMainEnumType(strMainBbb)
	EnumMainCcXxxxZzz = privateMainEnumType(strMainCcXxxxZzz)
)

type EnumMain interface {
	main()
	String() string
}

type privateMainEnumType string

func (s privateMainEnumType) main() {
	// unreachable method
}

func (s privateMainEnumType) String() string {
	return string(s)
}

func FromString(s string) (EnumMain, error) {
	switch s {
	case strMainSomething:
		return EnumMainSomething, nil
	case strMainAaa:
		return EnumMainAaa, nil
	case strMainBbb:
		return EnumMainBbb, nil
	case strMainCcXxxxZzz:
		return EnumMainCcXxxxZzz, nil
	}
	return nil, errors.New("unknown EnumMain " + s)
}

func ToStrings(ss []EnumMain) ([]string, error) {
	if slice == nil {
		return nil
	}
	result := make([]string, len(slice))
	for i, val := range slice {
		if val == nil {
			return nil, errors.New("unexpected enum EnumMain value: nil")
		}
		result[i] = val.String()
	}
	return result, nil
}

func MustToStrings(ss []EnumMain) []string {
	result, err := ToStrings(ss)
	if err != nil {
		panic(err.Error())
	}
	return result
}

func FromStrings(ss []strings) ([]EnumMain, error) {
	if slice == nil {
		return nil
	}
	result := make([]string, len(slice))
	for i, val := range slice {
		result[i] = val.String()
	}
	return result
}
`))

	if err != nil {
		t.Errorf(err.Error())
	}

	expected := string(bb)

	if expected != generated {
		t.Errorf("incorrect generated output from generate: \n %s\n", printDiff(expected, generated))
	}
}

func Test_1(t *testing.T) {
	bb, err := generate("Status", strings.Split("backlog in_review done", " ")...)
	if err != nil {
		t.Errorf(err.Error())
	}
	println(string(bb))
}
