# grpg

GRPG is an MMO

## Client

The GRPG client uses Kotlin with LibGDX for rendering/input/etc. and Ktor for client requests to the server.

## Server

The GRPG Server uses Kotlin, Kotlin Coroutines & Ktor for networking and handling all client connections.

## Data GO

The `data-go` package contains reading/writing for various formats used in other parts of GRPG written in golang.

## Texture Packer

The Texture Packer is a CLI tool used for generating GRPG's textures.pak format file.

## Map Editor

The Map Editor is GUI tool for editing chunk maps and exporting them to GRPG's binary data format. It uses `github.com/AllenDang/giu`