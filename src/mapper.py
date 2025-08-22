# mapper.py (or drop into your rendering module)

from typing import Iterable
import numpy as np
from src.lighting import Light
from src.controller import ArtnetController


def map_lights_to_artnet(
    lights: Iterable[Light],
    raw_rgba: bytes,   # length == width * height * 4
    width: int,
    height: int,
    artnet: ArtnetController,
    pixel_size: int = 1,       # “stretch” factor; 1 = no stretch
    flip_y: bool = False,   # <-- set True only if your LED coords are top-left origin
) -> None:
    """
    Pulls RGB from color_data using each light.position (x,y), optionally stretched by pixel_width,
    and writes it into the appropriate Art-Net universe/address.

    - If light.show is False, the pixel is blacked out.
    - Positions are clamped into the image bounds.
    - Stretch is nearest-neighbor: (sx, sy) = (x*pixel_width, y*pixel_width).
    """

    ps = pixel_size if pixel_size > 0 else 1
    stride = width * 4
    buf = memoryview(raw_rgba)

    for l in lights:
        if not l.show:
            artnet.blank_pixel(l.universe, l.addr)
            continue


        sx = l.position.x * ps
        sy = l.position.y * ps

        # clamp
        if sx < 0: sx = 0
        elif sx >= width: sx = width - 1
        if sy < 0: sy = 0
        elif sy >= height: sy = height - 1

        row = (height - 1 - sy) if flip_y else sy
        idx = row * stride + sx * 4

        r = buf[idx + 0]
        g = buf[idx + 1]
        b = buf[idx + 2]

        artnet.write_pixel(l.universe, l.addr, (r, g, b))
