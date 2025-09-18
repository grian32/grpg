package grpgscript_lsp

func UppercaseAll(str string) string {
	chars := []int32(str)

	for i, b := range str {
		if b >= 'a' && b <= 'z' {
			chars[i] = b - 32
		}
	}
	return string(chars)
}

func IsAlpha(char byte) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		char == '_'
}
