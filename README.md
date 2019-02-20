# orthus

Simple utility to recognise touchpad gestures, and log them to
stdout. Install + Usage:

```sh
$ pip install evdev

# list devices
$ sudo python orthus.py list
...

# by path
$ sudo python orthus.py /dev/input/event12

# by name
$ sudo python orthus.py 'DELL07E6:00 06CB:76AF Touchpad'

# to dispatch
$ sudo python orthus.py ... | ./dispatch.sh
```
