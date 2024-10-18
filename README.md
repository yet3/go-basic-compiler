# Basic tokenizer and parser in Go

## Usage
```sh
./tiny-basic-compiler [input_file] [output_file]

./tiny-basic-compiler ./examples/hello.tb ./dist/output.go
```

program ::= {statement}
statement ::= "PRINT" (expression | string)+ nl
  | "IF" comparison "THEN" nl {statement} "ENDIF" nl
  | "WHILE" comparison "REPEAT" nl {statement} "ENDWHILE" nl
  | "LABEL" ident nl
  | "GOTO" ident nl
  | "LET" ident "=" expression nl
  | "INPUT" ident nl
  | ident EQ (expression | string)+ nl
comparison ::= expression (("==" | "!=" | ">" | ">=" | "<" | "<=") expression)+
expression ::= term {( "-" | "+" ) term}
term ::= unary {( "/" | "*" ) unary}
unary ::= ["+" | "-"] primary
primary ::= number | ident
nl ::= '\n'+

