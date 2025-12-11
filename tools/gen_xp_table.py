#!/usr/bin/python3
from math import floor

xps = []

for i in range(74):
    xps.append(floor(211 * (i ** 2.5)))

build_str = "var LEVEL_XP [75]uint32 = [75]uint32 {\n"

for idx, xp in enumerate(xps):
    build_str += "    " + str(xp) + ",\n"

build_str += "}"

print(build_str)
