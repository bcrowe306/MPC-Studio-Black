# MPC Studio Black - Midi Speceification
This repository holds the midi specification for the MPC Studio Black. This information was not provided by AKAI Professional, and I am not affiliated with them in any manner. I am simply an owner of the product and have documented it's MIDI protocol for my personal use. The repo is a product of that my findings.


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
