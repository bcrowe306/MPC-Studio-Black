import mido
from mido import Message
import numpy as np
import cairo


midi_port_name = "MPC Studio Black MPC Public"

def choose_output():
    outPorts = mido.get_output_names()
    found_port = False
    outPortIndex: int
    for i, port in enumerate(outPorts):
        if port == midi_port_name:
            found_port = True
            outPortIndex = i
            break
    if not found_port:
        for i, port in enumerate(outPorts):
            print(f"{i}:, {port}")
        portNumber = input("Select output index ")
        outPortIndex = int(portNumber)
    return mido.open_output(outPorts[outPortIndex])

def choose_input():
    in_ports = mido.get_input_names()
    found_port = False
    in_port_index: int
    for i, port in enumerate(in_ports):
        if port == midi_port_name:
            found_port = True
            in_port_index = i
            break
    if not found_port:
        for i, port in enumerate(in_ports):
            print(f"{i}:, {port}")
        in_port_selection = input("Select output index ")
        in_port_index = int(in_port_selection)
    return mido.open_input(in_ports[in_port_index])

def gen_mpc_studio_black_line(line: int, pixel_array : list[int]):
    msg = Message("sysex", data = [0x47 ,0x7f, 0x3d, 0x04, 0x00, 0x7e, 0x02, 0x68, 0x00, 0x00, 0x00])
    msg.data += [line]
    msg.data += pixel_array
    return msg

in_port = choose_input()
out_port = choose_output()

def square_array(rgb24_data: list):
    return 0 if rgb24_data[0] < 127 else 1

def generate_pixel_encoding(bit_arry: list[int]):
    """Expects a list of int with a length of 3. The int should be binary, 0 || 1"""
    bitmap = [0x30,0x0c,0x03]
    bit = 0x00

    # Need to use numpy to do this operation. This will greatly reduce screen latency. Python is slow at looping 
    for i, pixel in enumerate(bit_arry):
        if pixel > 0:
            bit = bit | bitmap[i]
        else:
            bit = bit | 0x00
    return bit

def encode_image_buffer(buffer):
    pixels_array = np.split(buffer, 34560)
    dot_matrix = np.apply_along_axis(square_array,1, pixels_array)
    dot_matrix_triplets = np.split(dot_matrix, 11520)
    dot_matrix_encoded = np.apply_along_axis(generate_pixel_encoding, 1, dot_matrix_triplets)
    return np.split(dot_matrix_encoded, 96)

def create_surface(x: int = 360, y: int = 96):
    return cairo.ImageSurface(cairo.FORMAT_RGB24, x, y)

def init_surface(ctx: cairo.Context):
    """Set the background of the surface to all black. Needed for accurate dot_matrix style display"""
    ctx.set_source_rgb(0,0,0)
    ctx.rectangle(0,0, 360, 96)
    ctx.fill()

def draw_text(ctx: cairo.Context, text:str, x:int, y:int, font_size: int = 24, font: str="OpenSans"):
    ctx.set_source_rgb(1,1,1)
    ctx.set_font_size(font_size)
    ctx.select_font_face(font, cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_NORMAL)
    ctx.save()
    ctx.move_to(x,y)
    ctx.show_text(text)
    ctx.restore()

def send_sysex_lines_to_controller(lines_array, prev_lines_array):
    for line_number, byte_array in enumerate(lines_array):
        # Only send the lines that need to be updated. This should reduce the display latency
        if prev_lines_array is None:
            msg = gen_mpc_studio_black_line(int(line_number), byte_array)
            out_port.send(msg)
        else:
            comparison = byte_array == prev_lines_array[line_number]
            equal_arrays = comparison.all()
            if not equal_arrays:
                msg = gen_mpc_studio_black_line(int(line_number), byte_array)
                out_port.send(msg)

def render_surface_text(ctx, text):
    init_surface(ctx)
    draw_text(ctx,text, 10, 60, font_size=12)

surface = create_surface()
ctx = cairo.Context(surface)


prev_lines_array = None
for msg in in_port:
    render_surface_text(ctx, str(msg))
    lines = encode_image_buffer(surface.get_data())
    send_sysex_lines_to_controller(lines, prev_lines_array)
    prev_lines_array = lines
