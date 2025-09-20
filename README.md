# grpg
[![wakatime](https://wakatime.com/badge/user/6c690f0c-0d2d-4a6f-90f6-525168138404/project/9abeef4d-0b37-498e-8c50-64de0da52d4d.svg)](https://wakatime.com/badge/user/6c690f0c-0d2d-4a6f-90f6-525168138404/project/9abeef4d-0b37-498e-8c50-64de0da52d4d)

GRPG is a hobby MMO Game project that started out as a grid networking prototype and has since been growing into a full MMO Engine.

The philosophy for this is that while it's not the fastest way to develop a game, making everything custom and tailor-made to the game's needs rather than using general purpose engines and tools will provide much better performance and a much better developer experience.

Currently the project is still in its early stages. There isn't much content but there have been a good chunk of core systems implemented: NPCs, player inventories/saves, stateful/interactable objects, etc.

## Project Structure

### Client

The GRPG Client is located in `client-go/`, it uses Raylib for rendering

### Server

The GRPG server is located in `server-go/`. 

### Data GO

The `data-go` package contains reading/writing for various formats used in other parts of GRPG.

Most of the formats are self explanatory but part of the GRPGTEX format is that it encodes images as JPEGXL.

### Data Packer

A CLI tool located in `data-packer` using `cobra` & `charmbracelet/fang`, this mainly handles reading in from [GCFG](https://github.com/grian32/gcfg) manifest files and encoding them to the binary data formats used in GRPG. The manifests are mainly 1:1 with the actual binary data formats, the only transformations they do is converting PNG images to JPEGXL for the GRPGTEX format

### Map Editor

The Map Editor is GUI tool for editing chunk maps and exporting them to GRPG's binary data format. It uses `github.com/AllenDang/giu` for rendering the GUI.

### GRPGScript

GRPGScript is GRPG's content scripting language. It can technically be used independently of GRPG although you won't be able to do much with it. The language is meant to be used as a library.

There is also a [tree-sitter](https://github.com/grian32/tree-sitter-grpgscript) for GRPGScript, along with an LSP located in `grpgscript-lsp/`, a Zed Plugin is also located in `grpgscript-lsp/editors/zed` and an Intellij Plugin is planned.
