#!/bin/bash
go build -o grpgpack
./grpgpack tex -m testdata/test_tex_manifest.toml -o testdata/textures.grpgtex
./grpgpack tile -m testdata/test_tile_manifest.toml -o testdata/tiles.grpgtile -t testdata/textures.grpgtex
./grpgpack obj -m testdata/test_obj_manifest.toml -o testdata/objs.grpgobj -t testdata/textures.grpgtex
./grpgpack item -m testdata/test_item_manifest.toml -o testdata/items.grpgitem -t testdata/textures.grpgtex
rm -f grpgpack
