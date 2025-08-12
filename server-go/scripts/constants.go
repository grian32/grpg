package scripts

import (
	"errors"
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgobj"
	"grpgscript/object"
	"io"
	"os"
)

func LoadObjConstants(env *object.Environment, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	buf := gbuf.NewGBuf(bytes)
	header, err := grpgobj.ReadHeader(buf)
	if err != nil {
		return err
	}

	if header.Magic != [8]byte{'G', 'R', 'P', 'G', 'O', 'B', 'J', 0x00} {
		return errors.New("GRPGOBJ file does not have correct magic")
	}

	objs, err := grpgobj.ReadObjs(buf)
	if err != nil {
		return err
	}

	for _, obj := range objs {
		env.Set(uppercaseAll(obj.Name), &object.Integer{Value: int64(obj.ObjId)})
	}

	return nil
}

func uppercaseAll(str string) string {
	chars := []int32(str)

	for i, b := range str {
		if b >= 'a' && b <= 'z' {
			chars[i] = b - 32
		}
	}

	return string(chars)
}
