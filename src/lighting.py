from __future__ import annotations
from dataclasses import dataclass, field
from enum import Enum
from typing import List, Tuple


# --- Types / Enums ------------------------------------------------------------

class BoardType(Enum):
    Board7x7 = 0
    Board1x7 = 1


class BoardOrientation(Enum):
    Orient0   = 0
    Orient90  = 1
    Orient180 = 2
    Orient270 = 3


@dataclass
class Point:
    x: int
    y: int

    def alter(self, new_x: int, new_y: int) -> None:
        self.x = new_x
        self.y = new_y

    def is_equal(self, other: "Point") -> bool:
        return self.x == other.x and self.y == other.y

    def as_tuple(self) -> Tuple[int, int]:
        return (self.x, self.y)


@dataclass
class Light:
    position: Point
    pixel:    int
    show:     bool = True
    universe: int  = 0 
    addr:     int  = 1  


@dataclass
class BoardConfiguration:
    board_type: BoardType
    orientation: BoardOrientation
    starting_point: Point

    @staticmethod
    def new(t: BoardType, o: BoardOrientation, starting_point: Point) -> "BoardConfiguration":
        return BoardConfiguration(board_type=t, orientation=o, starting_point=starting_point)


# --- Board generator ----------------------------------------------------------

def generate_lights_from_boards(
    boards: List[BoardConfiguration],
    disallowed: List[Point],
    *,
    start_universe: int = 0,
    start_address: int = 1,
    dmx_slots_per_universe: int = 512,
) -> List[Light]:
    if not (1 <= start_address <= dmx_slots_per_universe):
        raise ValueError("start_address must be within 1..dmx_slots_per_universe (1..512).")

    lights: List[Light] = []
    last_light = 0
    cpp = 3  # RGB
    disallowed_set = {p.as_tuple() for p in disallowed}

    # Running DMX write head
    cur_uni = start_universe
    cur_addr = start_address  # 1-based

    for config in boards:
        if config.board_type == BoardType.Board7x7:
            led_count = 49
            max_nominal_y = 6
        else:
            led_count = 7
            max_nominal_y = 0

        for i in range(led_count):
            y_pos = i // 7
            is_odd_y = (y_pos % 2) != 0
            x_pos = i % 7
            if is_odd_y:
                x_pos = 6 - x_pos

            position = Point(config.starting_point.x, config.starting_point.y)
            o = config.orientation
            if o == BoardOrientation.Orient0:
                position.alter(position.x + y_pos, position.y + (6 - x_pos))
            elif o == BoardOrientation.Orient90:
                position.alter(position.x + x_pos, position.y + y_pos)
            elif o == BoardOrientation.Orient180:
                position.alter(position.x + (max_nominal_y - y_pos), position.y + x_pos)
            elif o == BoardOrientation.Orient270:
                position.alter(position.x + (6 - x_pos), position.y + (max_nominal_y - y_pos))

            is_allowed = (position.as_tuple() not in disallowed_set)

            # --- DMX placement (no rollover allowed) --------------------------
            # If this pixel (addr..addr+2) would exceed the universe, advance to a fresh one.
            if cur_addr + (cpp - 1) > dmx_slots_per_universe:
                cur_uni += 1
                cur_addr = 1  # start at slot 1 in the next universe

            lights.append(
                Light(
                    position=position,
                    pixel=last_light,
                    show=is_allowed,
                    universe=cur_uni,
                    addr=cur_addr,
                )
            )

            # Advance DMX write head by RGB width (always, even if hidden)
            cur_addr += cpp
            last_light += 1

    return lights, cur_uni + 1
