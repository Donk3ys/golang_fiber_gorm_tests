package util

import (
	"testing"

	"github.com/jaswdr/faker"
)

func TestFormatPhoneNumberLocal(t *testing.T) {
	numbersBad := []string{
		"12345",
		" 221097-330-9712 ",
	}
	for _, no := range numbersBad {
		_, err := FormatPhoneNumberLocal(no)
		t.Log(no)
		if err == nil {
			t.Fatalf("Phone number should error on format %v", no)
		}
	}

	numbersGood := []string{
		"079 330 9715",
		" +27 973309712 ",
		" 27 97 330 9712 ",
		"27-97-330-9711 ",
	}
	numbersGoodFmt := []string{
		"0793309715",
		"0973309712",
		"0973309712",
		"0973309711",
	}
	for i, no := range numbersGood {
		formatted, err := FormatPhoneNumberLocal(no)
		t.Log(no, " \t", formatted)
		if err != nil {
			t.Fatalf("Phone number should not error on format %v", no)
		}
		if formatted != numbersGoodFmt[i] {
			t.Fatalf("%v != %v", formatted, numbersGoodFmt[i])
		}
	}
}

func TestFormatPhoneNumberToInternational(t *testing.T) {
	numbersBad := []string{
		"12345",
		" 221097-330-9712 ",
	}
	for _, no := range numbersBad {
		_, err := FormatPhoneNumberToInternational(no)
		t.Log(no)
		if err == nil {
			t.Fatalf("Int phone number should error on format %v", no)
		}
	}

	numbersGood := []string{
		"079 330 9715",
		" +27 973309712 ",
		" 27 97 330 9712 ",
		"27-97-330-9711 ",
	}
	numbersGoodFmt := []string{
		"+27793309715",
		"+27973309712",
		"+27973309712",
		"+27973309711",
	}
	for i, no := range numbersGood {
		formatted, err := FormatPhoneNumberToInternational(no)
		t.Log(no, " \t", formatted)
		if err != nil {
			t.Fatalf("Int phone number should not error on format %v", no)
		}
		if formatted != numbersGoodFmt[i] {
			t.Fatalf("%v != %v", formatted, numbersGoodFmt[i])
		}
	}
}

func TestGenerateCode(t *testing.T) {
	n := 0
	for n < 10000 {
		code := GenerateCode()
		if code >= 10000 {
			t.Fatalf("Should not generate number larger than 10000 %v", code)
		}
		if code < 1000 {
			t.Fatalf("Should not generate number smaller than 1000 %v", code)
		}

		n++
	}
}

func TestCheckPasswordMatch(t *testing.T) {
	fake := faker.New()

	n := 0
	for n < 5 {
		pw := fake.Internet().Password()
		t.Log(pw)
		hash, err := GenPasswordHash(pw)
		if err != nil {
			t.Fatalf("Should generate password from %v\n%v", pw, err)
		}

		// Matches
		err = CheckPasswordMatch(hash, pw)
		if err != nil {
			t.Fatalf("Should match password from %v %v\n%v", hash, pw, err)
		}

		// Does not match
		diffPw := fake.Internet().Password()
		err = CheckPasswordMatch(hash, diffPw)
		if err == nil {
			t.Fatalf("Should not match password from %v %v\n%v", hash, diffPw, err)
		}

		n++
	}
}
