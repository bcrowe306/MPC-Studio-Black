package main

import (
	"bufio"
	"fmt"
	"image"
	"os"

	"github.com/fogleman/gg"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

var screen_width int = 360
var screen_height int = 96
var prev_image_buf [96][121]byte

func getPorts(inPortName string, outPortName string) (drivers.In, drivers.Out, error) {
	in_port, err := midi.FindInPort(inPortName)
	if err != nil {
		return nil, nil, err
	}
	out_port, err := midi.FindOutPort(outPortName)
	if err != nil {
		return nil, nil, err
	}
	return in_port, out_port, nil
}

func draw_image(ctx *gg.Context, text string) image.Image {
	ctx.DrawRectangle(0.0, 0.0, float64(screen_width), float64(screen_height))
	ctx.SetRGB(0.0, 0.0, 0.0)
	ctx.Fill()
	ctx.SetRGB(1.0, 1.0, 1.0)
	err := ctx.LoadFontFace("OpenSans.ttf", 11)
	if err != nil {
		fmt.Printf("ERROR: %s", err)
	}
	ctx.DrawString(text, 20.0, 48.0)
	return ctx.Image()
}

// send_sysex_msg: Sends a write screen sysex message to the MPC Studio Black.
func send_sysex_msg(out_port *drivers.Out, payload *[121]byte) {
	sysex_array := [132]byte{0x47, 0x7f, 0x3d, 0x04, 0x00, 0x7e, 0x02, 0x68, 0x00, 0x00, 0x00}
	for i := 0; i < len(*payload); i++ {
		sysex_array[i+11] = (*payload)[i]
	}
	send, err := midi.SendTo(*out_port)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}
	// fmt.Println(sysex_array)
	mMsg := midi.SysEx(sysex_array[:132])
	err = send(mMsg)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}

}

func write_to_screen(out_port *drivers.Out, image_data *image.Image, pre_alloc_image_buf *[121]byte) int {
	lines_drawn := 0
	for y := range screen_height {

		bit := 0
		pixelCount := 1
		pre_alloc_image_buf[0] = byte(y)
		for x := range screen_width {
			stride_position := x % 3
			r, _, _, _ := (*image_data).At(x, y).RGBA()
			var val = 0
			if r > 26214 {
				val = 3
			}
			bitShift := (2 - stride_position) * 2
			bit = bit | val<<bitShift
			if x > 0 && stride_position == 2 {
				(*pre_alloc_image_buf)[pixelCount] = byte(bit)
				// fmt.Printf("%#x ", bit)
				pixelCount++
				bit = 0
			}
		}
		if (*pre_alloc_image_buf) != (prev_image_buf[y]) {
			send_sysex_msg(out_port, pre_alloc_image_buf)
			prev_image_buf[y] = (*pre_alloc_image_buf)
			lines_drawn++
		}
	}
	return lines_drawn
}

func main() {
	var lines [121]byte
	ctx := gg.NewContext(screen_width, screen_height)

	defer midi.CloseDriver()
	inPort, outPort, err := getPorts("MPC Studio Black MPC Public", "MPC Studio Black MPC Public")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("Using Input: %s\n", inPort.String())
	fmt.Printf("Using Output: %s\n", outPort.String())

	stop, err := midi.ListenTo(inPort, func(msg midi.Message, timestampms int32) {

		// Generate image: Prints a blank screen with the string of midi msg. Outputs Image data in RGBA in uint32 (r,g,b,a) format
		img := draw_image(ctx, msg.String())

		// Write the lines to the screen
		lines_drawn := write_to_screen(&outPort, &img, &lines)
		fmt.Printf("Lines Drawn: %d\n", lines_drawn)
	}, midi.UseSysEx())

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(input)

	stop()
}
