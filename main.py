from typing import List

from pathlib import Path
import time
import moderngl
import numpy as np
from PIL import Image

import moderngl_window
from moderngl_window import resources
from moderngl_window.meta import ProgramDescription
from moderngl_window.utils.scheduler import Scheduler

from src.lighting import BoardConfiguration, BoardType, BoardOrientation, Point, generate_lights_from_boards
from src.controller import ArtnetController
from src.mapper import map_lights_to_artnet

BOARD_CONFIGURATION = [
    BoardConfiguration.new(BoardType.Board7x7, BoardOrientation.Orient0, Point(0, 0)),
    BoardConfiguration.new(BoardType.Board7x7, BoardOrientation.Orient0, Point(0, 7)),
    BoardConfiguration.new(BoardType.Board7x7, BoardOrientation.Orient0, Point(7, 7)),
    BoardConfiguration.new(BoardType.Board7x7, BoardOrientation.Orient0, Point(7, 0)),
]

BOARD_DISALLOWED_POSITIONS: List[Point] = []

LIGHTS, UNIVERSES = generate_lights_from_boards(BOARD_CONFIGURATION, BOARD_DISALLOWED_POSITIONS)

out_w, out_h = 14, 14
pixel_size = 33

screen_w, screen_h = out_w * pixel_size, out_h * pixel_size


ARTNET_CONTROLLER = ArtnetController("192.168.0.102", UNIVERSES)

this_path = Path(__file__).resolve()
this_dir = this_path.parent
shader_dir = this_dir / "data"

resources.register_program_dir(shader_dir)

shaders = [
    "basic", "checkerboard", "stripe-wheel-spread", "bar-hoppin", "tv-no-signal"
]

shader_change_time = 15.0

TARGET_FPS = 30.0
TARGET_DT  = 1.0 / TARGET_FPS


def load_shader(shader_name: str) -> moderngl.Program:
    return resources.programs.load(
        ProgramDescription(
            vertex_shader=f'{shader_name}.vert',
            fragment_shader=f'{shader_name}.frag',
        )
    )

class MainWindow(moderngl_window.WindowConfig):

    window_size = (screen_w, screen_h)
    resizable = False

    def __init__(self, **kwargs):
        super().__init__(**kwargs)

        self.scheduler = Scheduler(self.timer)
        self.shader_change_event = self.scheduler.run_every(
            self.change_shader, delay=shader_change_time, initial_delay=shader_change_time
        )

        self.shaders = [load_shader(shader) for shader in shaders]
        self.shader_index = 0

        vertices = np.array([
            -1.0, -1.0,  # bottom-left  -> 0
             1.0, -1.0,  # bottom-right -> 1
             1.0,  1.0,  # top-right    -> 2
            -1.0,  1.0,  # top-left     -> 3
        ], dtype='f4')

        indices = np.array([0, 1, 2, 0, 2, 3], dtype='i4')

        self.vbo = self.ctx.buffer(vertices.tobytes())
        self.ibo = self.ctx.buffer(indices.tobytes())

        self.build_vao()

        self.last_frame_ts = 0
        self.ctr = 0

        self._last_dump_ts = 0.0
        self._dump_dir = this_dir / "debug"
        self._dump_dir.mkdir(exist_ok=True)


    def build_vao(self):
        if hasattr(self, "vao") and self.vao is not None:
            self.vao.release()
        self.vao = self.ctx.vertex_array(
            self.shaders[self.shader_index],
            [(self.vbo, "2f", "vPosition")],
            index_buffer=self.ibo,
        )

    def on_key_event(self, key, action, modifiers):
        keys = self.wnd.keys
        if action == keys.ACTION_PRESS:
            if key == keys.SPACE:
                self.scheduler.cancel(self.shader_change_event)
                self.change_shader()
                self.shader_change_event = self.scheduler.run_every(
                    self.change_shader, delay=shader_change_time, initial_delay=shader_change_time
                )
                
    def change_shader(self):
        self.shader_index = (self.shader_index + 1) % len(self.shaders)
        self.build_vao()

    def maybe_dump_texture(self, color_data: bytes, color_size: tuple[int, int], every_s: float = 5.0) -> None:
        """Write color_data (H,W,4, uint8, RGB) to PNG every `every_s` seconds."""
        now = time.monotonic()
        if now - self._last_dump_ts < every_s:
            return

        # file names: timestamped and a rolling "latest.png" for quick peek
        ts = time.strftime("%Y%m%d-%H%M%S")
        out_path = self._dump_dir / f"frame-{ts}.png"
        latest   = self._dump_dir / "latest.png"


        image = Image.frombytes("RGBA", color_size, color_data)
        image = image.transpose(Image.FLIP_TOP_BOTTOM)
        image.save(out_path, format="png")
        image.save(latest, format="png")

        self._last_dump_ts = now
        print(f"[dump] wrote {out_path}")

    def on_render(self, render_time, _frame_time):

        elapsed = render_time - self.last_frame_ts

        if elapsed < TARGET_DT:
            time.sleep(TARGET_DT - elapsed)
        
        self.scheduler.execute()

        w, h = self.wnd.fbo.size
        self.ctx.viewport = (0, 0, w, h)
        self.ctx.clear(0.0, 0.0, 0.0, 0.0)


        self.shaders[self.shader_index]['time'].value = float(render_time)
        self.shaders[self.shader_index]['resolution'].value = (float(w), float(h))
        self.vao.render(mode=moderngl.TRIANGLES)

        color_data = self.wnd.fbo.read(components=4) 
        # self.maybe_dump_texture(color_data, (w, h))

        map_lights_to_artnet(LIGHTS, color_data, w, h, ARTNET_CONTROLLER, pixel_size)

        ARTNET_CONTROLLER.send()

        self.last_frame_ts = render_time




if __name__ == "__main__":
    moderngl_window.run_window_config(MainWindow, args=("--window", "glfw"))
