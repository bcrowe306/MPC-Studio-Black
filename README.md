![alt](mpc_studio_black.jpeg)

# MPC Studio Black - Midi Speceification
This repository holds the midi specification for the MPC Studio Black. This information was not provided by AKAI Professional, and I am not affiliated with them in any manner. I am simply an owner of the product and have documented it's MIDI protocol for my personal use. The repo is a product of that my findings.

## PAD Note Mapping
The pads on the MPC Studio Black communicate on channel 10(channel 9 in code) and are not mapped in contiguous numbers. They send NoteOn, NoteOff, and Aftertouch messages for each pad. Listed below is the mapping from Pad number to note. 

|Channel|Note|Pad Number|
|-------|----|----------|
|9      |37  |1         |
|9      |36  |2         |
|9      |42  |3         |
|9      |82  |4         |
|9      |40  |5         |
|9      |38  |6         |
|9      |46  |7         |
|9      |44  |8         |
|9      |48  |9         |
|9      |47  |10        |
|9      |45  |11        |
|9      |43  |12        |
|9      |49  |13        |
|9      |55  |14        |
|9      |51  |15        |
|9      |53  |16        |

## Function Buttons Mapping
The Function buttons to immediately North of the top row of Pads communicate on channel 1(0 in code), and send only NoteOn Message. On Press, the Function buttons send a velocity value of 127. On release, the buttons send a velocity of 0, which your daw will recognize as a NoteOff value, but it's important for you to know if you decide to develope custome software for this device. Listed below is a table for the Function Button Mapping.

|Channel|Note|Function Button|
|-------|----|------------|
|0      |12  |f1          |
|0      |13  |f2          |
|0      |14  |f3          |
|0      |15  |f4          |
|0      |33  |f5          |
|0      |34  |f6          |

## QLink Knobs
The QLinkj knobs to the left side of the controller send two types of messages:
* **NoteOn**: The knobs are touch sensitive. On touch, it will send NoteOn with velocity 127. On Release, it will send NoteOn with velocity 0.
* **ControlChange**: When turning the knob, they will send CC values from 0-127

The Scroll knob is not touch sensitive.

Listed below is a table mapping the knobs to there corresponding ControlChange values.
|Channel|CC(Control Change) |QLINK|
|-------|---|-----|
|0      |101 |Scroll|
|0      |16 |QLink1|
|0      |17 |QLink2|
|0      |18 |QLink3|
|0      |19 |QLink4|


## Jog Wheel
The Jog wheel is the only endless encoder on the device. It sends a signed bit on CC 100. This means it sends the follwing values when turned:
* **Clockwise turn**: MSG= channel:0, cc: 100, value: 1
* **Counter turn**: MSG= channel:0, cc: 100, value: 127


## Button Controls
The buttons on the MPC Studio Black primarily send NoteOn values of 127 for press and 0 for release. There are however a few buttons that send CC message in the same range instead of NoteOn message. The table below lists all the button mappings:

|Channel|Note|Name        |Control Change|
|-------|----|------------|--------------|
|0      |35  |PadBankA    |n/a           |
|0      |36  |PadBankB    |n/a           |
|0      |37  |PadBankC    |n/a           |
|0      |38  |PadBankD    |n/a           |
|0      |39  |FullLevel   |n/a           |
|0      |40  |16Level     |n/a           |
|0      |41  |StepSeq     |n/a           |
|0      |42  |NexSeq      |n/a           |
|0      |43  |TrackMute   |n/a           |
|0      |2   |ProgEdit    |n/a           |
|0      |3   |ProgMix     |n/a           |
|0      |4   |SeqEdit     |n/a           |
|0      |5   |SampleEdit  |n/a           |
|0      |8   |Song        |n/a           |
|0      |44  |Project     |n/a           |
|0      |45  |SeqEdit     |n/a           |
|0      |46  |Prog        |n/a           |
|0      |47  |Sample      |n/a           |
|0      |48  |NoFilter    |n/a           |
|0      |52  |Main        |n/a           |
|0      |111 |Numeric     |n/a           |
|0      |49  |Shift       |n/a           |
|0      |50  |Browser     |n/a           |
|0      |55  |DataMinus   |n/a           |
|0      |54  |DataPlus    |n/a           |
|0      |51  |Window      |n/a           |
|0      |53  |TapTempo    |n/a           |
|0      |56  |Up          |n/a           |
|0      |57  |Down        |n/a           |
|0      |65  |Left        |n/a           |
|0      |66  |Right       |n/a           |
|0      |67  |Undo        |n/a           |
|0      |n/a |**StepLeft**    |115           |
|0      |n/a |**StepRight**   |116           |
|0      |70  |GoTo        |n/a           |
|0      |71  |BarLeft     |n/a           |
|0      |72  |BarRight    |n/a           |
|0      |n/a |**Record**      |119           |
|0      |n/a |**Overdub**     |114           |
|0      |n/a |**Stop**        |117           |
|0      |n/a |**Play**        |118           |
|0      |n/a |**PlayStart**   |113           |
|0      |113 |QLinkTrigger|n/a           |
|0      |9   |Erase       |n/a           |
|0      |11  |NoteRepeat  |n/a           |



## Button LED Feedback
The MPC Studio Black is able to light the LEDs behind certain buttons using CC *(Control Change)* messages sent on channel 0 *(zero is channel one when observing midi and developing)*. 

#### Button Types
There are two types of LED buttons on this controller. **One Color** and **Two Color**. 
* **One Color**: These buttons light up with one color by sending a CC value of 1 to the controller on channel 0, and the corresponding note(0-127)
    * CC#: 0 = OFF
    * CC#: 1 = ON
* **Two Color**: These buttons light up with two colors by sending a CC value of (0-2) to the controller on channel 0, and the corresponding note(0-127)
    * CC#: 0 = OFF
    * CC#: 1 = ON - 1st Color
    * CC#: >=2 = ON - 2nd Color

Listed below is a table for all LED buttons on the controller.
You may use this table to send messages and light the buttons. Here is an example Midi Message in different formats...

Suppose I wanted to light the ProgEdit button with the second color. Using the table below, here is sudo python code to accomplish this.

```python
# Import mido
import mido

# List all output ports
ports = mido.get_output_names()

# Select and instantiate the correct port
portIndex: int = 6
outport = mido.open_output(ports[portIndex])

# Create the msg to send (disreagard time for this example, n/a)
msg = Message('control_change', channel=0, control=2, value=2, time=0)

# Send the message
outport.send(msg)

```

|Channel | Note/Data1 |Name |Button Type |
|-------|----|------------|---------|
|0      |2   |ProgEdit    |two_color|
|0      |3   |ProgMix     |two_color|
|0      |5   |SeqEdit     |two_color|
|0      |6   |SampleEdit  |two_color|
|0      |8   |Song        |two_color|
|0      |35  |PadBankA    |two_color|
|0      |36  |PadBankB    |two_color|
|0      |37  |PadBankC    |two_color|
|0      |38  |PadBankD    |two_color|
|0      |39  |FullLevel   |one_color|
|0      |40  |16Level     |one_color|
|0      |41  |StepSeq     |one_color|
|0      |42  |NextSeq     |one_color|
|0      |43  |TrackMute   |two_color|
|0      |44  |Project     |one_color|
|0      |45  |Seq         |one_color|
|0      |46  |Prog        |one_color|
|0      |47  |Sample      |one_color|
|0      |48  |NoFilter    |one_color|
|0      |50  |Browser     |one_color|
|0      |51  |Window      |one_color|
|0      |52  |Main        |two_color|
|0      |53  |TapTempo    |one_color|
|0      |67  |Undo        |one_color|
|0      |73  |Rec         |one_color|
|0      |80  |OverDub     |one_color|
|0      |82  |Play        |one_color|
|0      |112 |PadAssign   |one_color|
|0      |113 |QlinkTrigger|one_color|
