# Thingsboard-stub

The name is quite self-explanatory : this program is meant to simulate a device (an arduino/industruino) interacting
with a Thingsboard server, for testing purposes.

For now you can simulate two types of devices :
- A boolean type : like the state of the switch
- A float type : for temperatures for example.

### Usage

```
J:\dev\tests-econtrols-supervisor>main.exe --help
NAME:
   ecstub - Stub for the eControls Supervisor run by Thingsboard

USAGE:
   main.exe [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --address value, -a value       The address of the broker
   --port value, -p value          Port of the MQTT broker (default: 1883)
   --token value                   The token (identifier) of the device on thingsboard
   --switches value, -s value      Path to the switches definition file (default: "resources/default_switches.json")
   --temperatures value, -t value  Path to the temperatures definition file (default: "resources/default_temperatures.json")
   --help, -h                      show help
   --version, -v                   print the version
```