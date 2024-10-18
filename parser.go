package main

import (
	"fmt"
	"log"
	"slices"
)

type Parser struct {
	curToken  Token
	peekToken Token
	t         *Tokenizer
	e         *Emitter

	symbols []string
	labels  []string
	gotos   []string
}

func NewParser(t *Tokenizer, e *Emitter) *Parser {
	return &Parser{
		symbols: []string{},
		labels:  []string{},
		gotos:   []string{},
		t:       t,
		e:       e,
	}
}

func (p *Parser) checkKind(kind int) bool {
	return p.curToken.kind == kind
}

func (p *Parser) matchKind(kind int) (Token, bool) {
	curToken := p.curToken
	if curToken.kind != kind {
		log.Fatalf("Kind not matching, expected: %v, got: %v\n", kind, p.curToken.kind)
	}
	p.nextToken()
	return curToken, true
}

func (p *Parser) nextToken() {
	token, err := p.t.GetToken()
	if err != nil {
		log.Fatalf("Couldn't get token: %v\n", err)
	}
	p.curToken = p.peekToken
	p.peekToken = token
}

func (p *Parser) newLine() {
	p.e.EmitCode("\n")
	p.matchKind(NEWLINE)

	for {
		if !p.checkKind(NEWLINE) {
			break
		}
		p.nextToken()
	}
}

func (p *Parser) primary() {
	p.e.EmitCode(p.curToken.value)
	if p.checkKind(NUMBER) {
		p.nextToken()
	} else if p.checkKind(IDENT) {
		if !slices.Contains(p.symbols, p.curToken.value) {
			log.Fatalf("Accessing undeclared variable: %v\n", p.curToken.value)
		}
		p.nextToken()
	} else {
		log.Fatal("Not a primary")
	}
}

func (p *Parser) unary() {
	if p.checkKind(PLUS) || p.checkKind(MINUS) {
		p.e.EmitCode(p.curToken.value)
		p.nextToken()
	}
	p.primary()
}

func (p *Parser) term() {
	p.unary()

	for {
		if !(p.checkKind(ASTERISK) || p.checkKind(SLASH)) {
			break
		}
		p.e.EmitCode(p.curToken.value)
		p.nextToken()
		p.unary()
	}
}

func (p *Parser) expression() {
	p.term()

	for {
		if !(p.checkKind(PLUS) || p.checkKind(MINUS)) {
			break
		}
		p.e.EmitCode(p.curToken.value)
		p.nextToken()
		p.term()
	}
}

func isComparisonOp(token Token) bool {
	if token.kind == EQEQ || token.kind == NOTEQ || token.kind == GT || token.kind == GTEQ || token.kind == LT || token.kind == LTEQ {
		return true
	}
	return false
}

func (p *Parser) comparison() {
	p.expression()

	if isComparisonOp(p.curToken) {
		p.e.EmitCode(p.curToken.value)
		p.nextToken()
		p.expression()
	} else {
		log.Fatalf("Is not comparison operator: %v\n", p.curToken.value)
	}

	for {
		if !isComparisonOp(p.curToken) {
			break
		}
		p.e.EmitCode(p.curToken.value)
		p.nextToken()
		p.expression()
	}
}

func (p *Parser) statement() {
	if p.checkKind(PRINT) {
		p.e.AddImport("fmt")

		p.e.EmitCode("fmt.Println(")

		p.nextToken()

		for {
			if p.checkKind(NEWLINE) {
				break
			}

			if p.checkKind(STRING) {
				p.e.EmitCode(fmt.Sprintf("\"%v\"", p.curToken.value))
				p.nextToken()
			} else if p.checkKind(COMMA) {
				p.e.EmitCode(",")
				p.nextToken()
			} else {
				p.expression()
			}
		}

		p.e.EmitCode(")")
	} else if p.checkKind(IF) {
		p.e.EmitCode("if ")
		p.nextToken()
		p.comparison()

		p.e.EmitCode(" {")
		p.matchKind(THEN)
		p.newLine()

		for {
			if p.checkKind(ENDIF) {
				break
			}
			p.e.EmitCode("")
			p.statement()
		}
		p.e.EmitCode("}")
		p.matchKind(ENDIF)

	} else if p.checkKind(WHILE) {
		p.e.EmitCode("for ")
		p.nextToken()
		p.comparison()
		p.e.EmitCode(" {\n")
		p.matchKind(REPEAT)
		p.newLine()
		for {
			if p.checkKind(ENDWHILE) {
				break
			}
			p.statement()
		}
		p.e.EmitCode("}")
		p.matchKind(ENDWHILE)
	} else if p.checkKind(LABEL) {
		p.nextToken()

		p.e.EmitCode(p.curToken.value + ":")
		if slices.Contains(p.labels, p.curToken.value) {
			log.Fatalf("Label is already declared: %v\n", p.curToken.value)
		} else {
			p.labels = append(p.labels, p.curToken.value)
		}

		p.matchKind(IDENT)
	} else if p.checkKind(GOTO) {
		p.nextToken()

		p.e.EmitCode("goto " + p.curToken.value)
		if !slices.Contains(p.gotos, p.curToken.value) {
			p.gotos = append(p.gotos, p.curToken.value)
		}

		p.matchKind(IDENT)
	} else if p.checkKind(INPUT) {
		p.e.AddImport("fmt")
		p.nextToken()

		if !slices.Contains(p.symbols, p.curToken.value) {
			p.symbols = append(p.symbols, p.curToken.value)
			p.e.EmitCode(fmt.Sprintf("var %v = \"\"\n", p.curToken.value))
		}
		p.e.EmitCode(fmt.Sprintf("fmt.Scanln(&%v)", p.curToken.value))

		p.matchKind(IDENT)
	} else if p.checkKind(LET) {
		p.nextToken()

		if !slices.Contains(p.symbols, p.curToken.value) {
			p.symbols = append(p.symbols, p.curToken.value)
		}

		p.e.EmitCode(fmt.Sprintf("var %v", p.curToken.value))
		p.matchKind(IDENT)
		p.e.EmitCode(" = ")
		p.matchKind(EQ)
		p.expression()
	} else if p.checkKind(IDENT) {
		p.e.EmitCode(p.curToken.value)
		p.nextToken()
		p.e.EmitCode(" = ")
		p.matchKind(EQ)

		for {
			if p.checkKind(NEWLINE) {
				break
			}

			if p.checkKind(STRING) {
				p.e.EmitCode(fmt.Sprintf("\"%v\"", p.curToken.value))
				p.nextToken()
			} else {
				p.expression()
			}
		}
	} else {
		log.Fatalf("WRONG STATMENT: %v\n", p.curToken.value)

	}

	p.newLine()
}

func (p *Parser) Parse() {
	p.nextToken()
	p.nextToken()
	fmt.Println("Parsing...")

	p.e.pkg = "main"
	p.e.EmitCode("func main() {\n")

	for {
		if p.checkKind(EOF) {
			fmt.Println("Parsing has finished")
			break
		}
		p.statement()
	}

	for _, label := range p.gotos {
		if !slices.Contains(p.labels, label) {
			log.Fatalf("Label doesn't exist: %v\n", label)
		}
	}

	p.e.EmitCode("}")
}

