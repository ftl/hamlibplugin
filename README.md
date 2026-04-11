# Hamlib OpenAction Plugin

An OpenAction ([OpenDeck](https://github.com/nekename/OpenDeck)) plugin to control ham radio transceivers through a hamlib rigctld daemon.

## Configuration

The plugin currently supports only one rigctld connection to `localhost:4532`.

## Actions

The following actions are currently available:

| Action | rigctld Command | Description |
|--------|----------------|-------------|
| Frequency Dial | `\set_freq`, `\get_freq` | Encoder to tune the VFO frequency in configurable steps |
| Level Encoder | `\set_level`, `\get_level` | Encoder to adjust a level value in configurable steps |
| On/Off | `\get_powerstat`, `\set_powerstat` | Toggle the rig power between on and off |
| Power State | `\set_powerstat` | Set the rig to a specific power state (off, on, standby, operate) |
| RIT | `\set_rit`, `\get_rit`, `\set_func` | Encoder to adjust the RIT offset, press/keypad to toggle RIT on/off |
| Select Mode | `\set_mode` | Set the operating mode for a VFO |
| Select VFO | `\set_vfo` | Select the active VFO |
| Send Morse | `\send_morse` | Send a preconfigured text as Morse code |
| Set Antenna | `\set_ant` | Select an antenna for a VFO |
| Set Frequency | `\set_freq` | Set the frequency of a VFO to a fixed value |
| Set Frequency (rel) | `\set_freq`, `\get_freq` | Set a VFO frequency relative to another VFO's frequency |
| Set Function | `\set_func` | Enable or disable a rig function |
| Set Level | `\set_level` | Set a level to a fixed value |
| Set Parameter | `\set_parm` | Set a rig-wide parameter value |
| Set Split VFO | `\set_split_vfo` | Enable or disable split operation and set the TX VFO |
| Stop Morse | `\stop_morse` | Stop the current Morse code transmission |
| Toggle Function | `\get_func`, `\set_func` | Toggle a rig function on/off |
| VFO Operation | `\vfo_op` | Perform a VFO operation (e.g. copy, exchange, band up/down) |
| VFO Operation Encoder | `\vfo_op` | Encoder with configurable VFO operations for clockwise, counter-clockwise, and press |
| XIT | `\set_xit`, `\get_xit`, `\set_func` | Encoder to adjust the XIT offset, press/keypad to toggle XIT on/off |

## License
This software is published under the [MIT License](https://www.tldrlegal.com/l/mit).

Copyright [Florian Thienel](http://thecodingflow.com/)
