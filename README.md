# orthus

Simple utility to recognise touchpad gestures, and log them to stdout.

### Install

```sh
$ pip install evdev
```

### Usage

If you use pyenv, orthus automatically detects that and runs `sudo ...`
with the correct python path.

```
# list devices
$ ./orthus list

# by path
$ ./orthus /dev/input/event12

# by name
$ ./orthus 'DELL07E6:00 06CB:76AF Touchpad'

# to dispatch
$ ./orthus ... | ./dispatch.sh
```

See `dispatch.sh` for an example of what could be done.
