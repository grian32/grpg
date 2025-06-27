# texture-packer

Packs PNG textures into a texture.pak, this basically contains a "map" of texture ids to their PNG files.

Manifest File Format:
```
texture_name=texture_file_path.png
```

`texture_name` is the internal name for the texture, while `texture_file_path.png` is the path to the texture in PWD, outputs a textures.pak in PWD

If another program is using the textures.pak, and you change the `texture_name`, aforementioned programs may break as that is what it uses as a reference.

# Endianness

The textures.pak output by this texture packer uses big endian ordering as that is what's generally expected by other libraries this file will be read by.

# Testing

This project contains some testdata in `testdata/`, you can pack this by running the `test.sh` script found in the root of the project, which will build the project and run it on the testdata. This will output a textures.pak in `testdata/`