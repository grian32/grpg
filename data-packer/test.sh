#!/bin/bash
go build -o grpgpack
./grpgpack tex -m testdata/test_tex_manifest.gcfg -o testdata/textures.grpgtex
./grpgpack tile -m testdata/test_tile_manifest.gcfg -o testdata/tiles.grpgtile -t testdata/textures.grpgtex
./grpgpack obj -m testdata/test_obj_manifest.gcfg -o testdata/objs.grpgobj -t testdata/textures.grpgtex
./grpgpack item -m testdata/test_item_manifest.gcfg -o testdata/items.grpgitem -t testdata/textures.grpgtex
./grpgpack npc -m testdata/test_npc_manifest.gcfg -o testdata/npcs.grpgnpc -t testdata/textures.grpgtex
rm -f grpgpack testdata/items.grpgitem testdata/npcs.grpgnpc testdata/objs.grpgobj testdata/textures.grpgtex testdata/tiles.grpgtile
