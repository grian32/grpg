# texture-packer

Packs PNG textures into a texture.pak

Manifest File Format:
```
texture_name=texture_file_path.png
```

`texture_name` is the internal name for the texture, while `texture_file_path.png` is the path to the texture in PWD, outputs a textures.pak in PWD

If another program is using the textures.pak, and you change the `texture_name`, aforementioned programs may break as that is what it uses as a reference.