# grpg

GRPG is an MMO Game project that started out as a prototype mmo project and later was further fleshed out.

At this time the project is still in a fairly early stage but both the client/server load map files w/ collision and multiple players can connect & move around etc.

## Client

The GRPG client is located in `client-go/` uses Go with Raylib-Go for rendering/input/etc. along with the stdlib TCP networking utilities for networking. There is also the initial Kotlin client in `client/`, which used LibGDX/Ktor but I later decided to rewrite it in GO.

## Server

The GRPG server is located in `server-go/`, this uses no external dependencies & takes advantage of go's goroutines.

## Data GO

The `data-go` package contains reading/writing for various formats used in other parts of GRPG.

## Texture Packer

The Texture Packer is a CLI tool used for generating GRPG's textures.pak format file.

## Map Editor

The Map Editor is GUI tool for editing chunk maps and exporting them to GRPG's binary data format. It uses `github.com/AllenDang/giu` for rendering the GUI.

## GRPGScript

GRPGScript is GRPG's custom scripting language for content. It's source code is located in `grpgscript/` and currently the plan is to accept "providers" for various content related built-ins, while fleshing out grpgscript as a reasonable to use language.