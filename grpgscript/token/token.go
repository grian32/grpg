package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

// TODO: convert to byte, but i'm just following along with the book now
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENT"
	INT   = "INT"

	ASSIGN = "="
	PLUS   = "+"

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "FUNCTION"
	VAR      = "VAR"
)

var keywords = map[string]TokenType{
	"fnc": FUNCTION,
	"var": VAR,
}

func LookupIdent(ident string) TokenType {
	tok, exists := keywords[ident]
	if exists {
		return tok
	}

	return IDENT
}
