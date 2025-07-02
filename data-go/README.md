# data-go

The data-go package contains binary reading/writing for various formats used in GRPG

## GBuf

A wrapper for buffer writing/reading for convenience, as the default stdlib functions can get lengthy when used many times.

## GRPGTex

GRPGTex is the format used for packing textures & various other data such as the type of the tile into one file, it originally started as packing only texture 
but I've needed extra data since and it seemed like the natural place to add it.

Each texture contains an internal ID as both a string and an uint16, the string is for clarity in tools like map editors etc., and the u16 ID is for representation in other formats.
Both are defined by the user so that upon packing nothing is automatically generated that could make people need to "recompile" their maps or whatever else.

Each texture also contains the PNG of the texture, this is the easiest way I found to achieve a fairly compressed format without much work on compression.