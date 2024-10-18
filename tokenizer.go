package main

import (
	"fmt"
	"log"
	"unicode"
)

type Token struct {
	value string
	kind  int
}

const (
	EOF     = -1
	NEWLINE = 0
	NUMBER  = 1
	IDENT   = 2
	STRING  = 3

	// KEYWORDS
	LABEL    = 101
	GOTO     = 102
	PRINT    = 103
	INPUT    = 104
	LET      = 105
	IF       = 106
	THEN     = 107
	ENDIF    = 108
	WHILE    = 109
	REPEAT   = 110
	ENDWHILE = 111

	EQ       = 201
	PLUS     = 202
	MINUS    = 203
	ASTERISK = 204
	SLASH    = 205
	EQEQ     = 206
	NOTEQ    = 207
	LT       = 208
	LTEQ     = 209
	GT       = 210
	GTEQ     = 211
	COMMA    = 212
)

var keywords = map[string]int{
	"LABEL":    LABEL,
	"GOTO":     GOTO,
	"PRINT":    PRINT,
	"INPUT":    INPUT,
	"LET":      LET,
	"IF":       IF,
	"THEN":     THEN,
	"ENDIF":    ENDIF,
	"WHILE":    WHILE,
	"REPEAT":   REPEAT,
	"ENDWHILE": ENDWHILE,
}

type Tokenizer struct {
	src      string
	caretPos int
}

func NewTokenizer(src string) *Tokenizer {
	return &Tokenizer{
		src:      src,
		caretPos: 0,
	}
}

func (t *Tokenizer) skipWhitespaces() {
	for {
		if !(t.GetChar() == " " || t.GetChar() == "\\t" || t.GetChar() == "\\r") {
			break
		}
		t.nextChar()
	}
}

func (t *Tokenizer) nextChar() {
	t.caretPos++
}

func (t *Tokenizer) nextCharBy(by int) {
	t.caretPos += by
}

func (t *Tokenizer) GetChar() string {
	if t.caretPos > len(t.src)-1 {
		return "\\0"
	}

	return string(t.src[t.caretPos])
}

func (t Tokenizer) GetNextChar() string {
	if t.caretPos > len(t.src)-1 {
		return "\\0"
	}

	return string(t.src[t.caretPos+1])
}

func MatchKeyword(keyword string) (matches bool, kind int) {
	kind, exists := keywords[keyword]
	if exists {
		return true, kind
	}
	return false, -999
}

func (t *Tokenizer) GetToken() (Token, error) {
	t.skipWhitespaces()
	curChar := t.GetChar()

	switch curChar {
	case "+":
		t.nextChar()
		return Token{value: curChar, kind: PLUS}, nil
	case "-":
		if unicode.IsDigit(rune(t.GetNextChar()[0])) {
			break
		} else {
			t.nextChar()
			return Token{value: curChar, kind: MINUS}, nil
		}
	case "*":
		t.nextChar()
		return Token{value: curChar, kind: ASTERISK}, nil
	case "/":
		t.nextChar()
		return Token{value: curChar, kind: SLASH}, nil
	case "=":
		if t.GetNextChar() == "=" {
			t.nextCharBy(2)
			return Token{value: "==", kind: EQEQ}, nil
		}
		t.nextChar()
		return Token{value: curChar, kind: EQ}, nil
	case "!":
		if t.GetNextChar() == "=" {
			t.nextCharBy(2)
			return Token{value: "!=", kind: NOTEQ}, nil
		}
	case ">":
		if t.GetNextChar() == "=" {
			t.nextCharBy(2)
			return Token{value: ">=", kind: GTEQ}, nil
		}
		t.nextChar()
		return Token{value: curChar, kind: GT}, nil
	case "<":
		if t.GetNextChar() == "=" {
			t.nextCharBy(2)
			return Token{value: "<=", kind: LTEQ}, nil
		}
		t.nextChar()
		return Token{value: curChar, kind: LT}, nil
	case "\"":
		pos := t.caretPos + 1
		str := ""

		for {
			if pos > len(t.src)-1 {
				return Token{value: str, kind: STRING}, fmt.Errorf("Never ending string")
			}

			char := string(t.src[pos])

			if char == "\"" {
				break
			}
			str += char
			pos++
		}

		t.nextCharBy(pos - t.caretPos + 1)
		return Token{value: str, kind: STRING}, nil
	case ",":
		t.nextChar()
		return Token{value: curChar, kind: COMMA}, nil
	case "\n":
		t.nextChar()
		return Token{value: curChar, kind: NEWLINE}, nil
	case "\\0":
		t.nextChar()
		return Token{value: curChar, kind: EOF}, nil

	}

	if unicode.IsDigit(rune(curChar[0])) || curChar == "-" {
		pos := t.caretPos
		str := ""
		if curChar == "-" {
			pos++
			str = "-"
		}
		for {
			if pos > len(t.src)-1 {
				return Token{value: str, kind: NUMBER}, fmt.Errorf("Number is invalid")
			}
			byte := t.src[pos]
			char := string(byte)
			if !unicode.IsDigit(rune(byte)) && char != "." {
				break
			}
			str += char
			pos++
		}

		t.nextCharBy(pos - t.caretPos)
		return Token{value: str, kind: NUMBER}, nil
	}

	if unicode.IsLetter(rune(curChar[0])) {
		pos := t.caretPos
		str := ""
		for {
			if pos > len(t.src)-1 {
				return Token{value: str, kind: NUMBER}, fmt.Errorf("Text is invalid")
			}
			byte := t.src[pos]
			char := string(byte)
			if !unicode.IsLetter(rune(byte)) {
				break
			}
			str += char
			pos++
		}

		t.nextCharBy(pos - t.caretPos)
		isKeyword, kind := MatchKeyword(str)
		if isKeyword {
			return Token{value: str, kind: kind}, nil
		}
		return Token{value: str, kind: IDENT}, nil
	}

	return Token{}, fmt.Errorf("Unknown token: %v", curChar)
}

func (t *Tokenizer) Lex() []Token {
	tokens := []Token{}
	for {
		token, err := t.GetToken()
		if err != nil {
			log.Fatalf("Tokenizer error: %v", err)
		}
		tokens = append(tokens, token)
		if token.kind == EOF {
			fmt.Printf("Tokenization has finished, tokens: %v\n", len(tokens))
			break
		}
	}
	return tokens
}
