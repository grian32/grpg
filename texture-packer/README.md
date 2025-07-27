# texture-packer

Packs PNG textures into a texture.pak, this basically contains a "map" of texture ids to their PNG files, along with some extra metadata.

## Manifest File Format

```toml
texture = [
    { name = "texture_name", id = 0, path = "texture_file_path.png", type = "TILE" },
    { name = "texture_name_2", id = 1, path = "texture_file_path.png_2", type = "OBJ" },
]
```

`texture` is an array of textures metadata definitions:

- `name` is the internal name for the texture, this is used for editors etc.
- `id` is the internal
- `path` is the path to the texture in relation to current working dir
- `type` is the type used by the texture format and translated into other formats.

If you change the `id` for a texture, this may break compatibility with other software such as map files etc that were created using a previous textures.pak file, as that's the id it uses to reference each textures.

# Endianness

The textures.pak output by this texture packer uses big endian ordering as that is what's generally expected by other libraries this file will be read by.

# Testing

This project contains some testdata in `testdata/`, you can pack this by running the `test.sh` script found in the root of the project, which will build the project and run it on the testdata. This will output a textures.pak in `testdata/`
