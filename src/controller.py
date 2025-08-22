# artnet_system.py
from __future__ import annotations
import socket
from typing import Dict


class ArtnetController:
    """
    Art-Net sender that PREBUILDS one packet per universe.
    You mutate the 512-byte DMX payload in-place; send() only updates the sequence byte.

    address:         destination node IP (e.g., "192.168.0.50")
    universes:       number of consecutive universes to control
    start_universe:  absolute starting universe index (0-based)
    port:            UDP port (default 6454)
    """

    # ArtDmx fixed layout offsets
    _OFF_ID      = 0          # 8 bytes: "Art-Net\0"
    _OFF_OPCODE  = 8          # 2 bytes: little-endian 0x5000
    _OFF_PVER    = 10         # 2 bytes: big-endian protocol version (0x000E)
    _OFF_SEQ     = 12         # 1 byte : sequence (0 disables)
    _OFF_PHYS    = 13         # 1 byte : physical port
    _OFF_SUBUNI  = 14         # 1 byte : low 8 bits of universe
    _OFF_NET     = 15         # 1 byte : high bits (7-bit)
    _OFF_LEN     = 16         # 2 bytes: big-endian DMX length
    _OFF_DATA    = 18         # start of DMX payload

    def __init__(self, address: str, universes: int, *, start_universe: int = 0, port: int = 6454) -> None:
        if universes <= 0:
            raise ValueError("universes must be >= 1")

        self.address = address
        self.port = port
        self.start_universe = int(start_universe)
        self.universe_count = int(universes)

        # One prebuilt packet per universe: header (18 bytes) + 512-byte payload
        self._packets: Dict[int, bytearray] = {}
        self._dmx_views: Dict[int, memoryview] = {}

        for i in range(self.universe_count):
            abs_uni = self.start_universe + i
            pkt = bytearray(self._OFF_DATA + 512)

            # Header
            pkt[self._OFF_ID:self._OFF_ID+8] = b"Art-Net\x00"
            pkt[self._OFF_OPCODE:self._OFF_OPCODE+2] = bytes((0x00, 0x50))   # ArtDmx
            pkt[self._OFF_PVER:self._OFF_PVER+2] = bytes((0x00, 0x0E))       # ProtVer 14
            pkt[self._OFF_SEQ]  = 0                                           # sequence (0=disable)
            pkt[self._OFF_PHYS] = 0
            pkt[self._OFF_SUBUNI] = abs_uni & 0xFF
            pkt[self._OFF_NET]    = (abs_uni >> 8) & 0x7F
            pkt[self._OFF_LEN:self._OFF_LEN+2] = bytes((0x02, 0x00))         # 512 (big-endian)

            # Payload is initially zero; keep a view for fast writes
            self._packets[abs_uni] = pkt
            self._dmx_views[abs_uni] = memoryview(pkt)[self._OFF_DATA:self._OFF_DATA+512]

        self._sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        self._sock.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)
        self._sequence = 0  # set to 0 to disable, or increment each send if you want ordering

    # --- Public API ------------------------------------------------------------

    def get_buffer(self, universe: int) -> memoryview:
        """
        Returns a mutable 512-byte memoryview for the given ABSOLUTE universe’s DMX data.
        Example: buf = art.get_buffer(0); buf[0:3] = (255, 0, 0)
        """
        return self._dmx_views[universe]

    def set_buffer(self, universe: int, data: bytes | bytearray) -> None:
        """
        Replace the 512-byte DMX payload for the given ABSOLUTE universe.
        """
        if len(data) != 512:
            raise ValueError("DMX data must be exactly 512 bytes.")
        self._dmx_views[universe][:] = data

    def write_pixel(self, universe: int, addr: int, rgb: tuple[int, int, int]) -> None:
        if not (1 <= addr <= 512):
            raise ValueError(f"addr must be 1..512, got {addr}")
        if addr + 2 > 512:
            raise RuntimeError(
                f"DMX rollover forbidden: universe={universe}, addr={addr} (needs {addr+2} > 512)"
            )

        r, g, b = (int(rgb[0]) & 0xFF, int(rgb[1]) & 0xFF, int(rgb[2]) & 0xFF)
        buf = self._dmx_views[universe]
        i = addr - 1  # 0-based

        # assume well behaved universes
        buf[i] = r
        buf[i + 1] = g
        buf[i + 2] = b

    def blank_pixel(self, universe: int, addr: int) -> None:
        """Set the pixel at (universe, addr) to black (handles spill)."""
        self.write_pixel(universe, addr, (0, 0, 0))


    def send(self, *, use_sequence: bool = True) -> None:
        """
        Send all prebuilt packets. If use_sequence=True, increments and writes sequence per send.
        """
        if use_sequence:
            self._sequence = (self._sequence + 1) & 0xFF

        seq = self._sequence if use_sequence else 0

        for _, pkt in self._packets.items():
            pkt[self._OFF_SEQ] = seq
            self._sock.sendto(pkt, (self.address, self.port))

    def close(self) -> None:
        try:
            self._sock.close()
        except Exception:
            pass
