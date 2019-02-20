# midas

Simple utility to recognise touchpad gestures, and log them to stdout.

### Install

```sh
$ pip install evdev
```

### Usage

If you use pyenv, midas automatically detects that and runs `sudo ...`
with the correct python path.

```
# list devices
$ ./midas list

# by path
$ ./midas /dev/input/event12

# by name
$ ./midas 'DELL07E6:00 06CB:76AF Touchpad'

# to dispatch
$ ./midas ... | ./dispatch.sh
```

See `dispatch.sh` for an example of what could be done.

### Algorithm

The meat of the "gesture guessing" algorithm is deciding, given a series of X and Y
coordinates of the user's finger on the touchpad, the overall direction of the gesture.
To do this we first need to gather a buffer of the events:

```
buffer = []
while True:
    try:
        # wait 0.1s for next event
        buffer.append(get_next_event(wait=0.1))
    except Timeout:
        process_buffer(buffer)
        clear buffer
```

Once we have the events we separate them into X and Y components:

```
events = [(x, 1), (y, 2), (x, 2), (y, 3), ...]
 => X = [1, 2, ...]
 => Y = [2, 3, ...]
```

Now the user's gestures are going to naturally have some jitter due to the natural hand
movements. To decide whether the _overall_ direction is left/right (or, up/down) we do
the following.

1. Split up the data into overlapping chunks of 4.

   ```
   [1,2,1,2,3,4,3,4,..] => [1,2,1,2], [2,1,2,3], [1,2,3,4], ...
   ```

2. We calculate the agreed direction of each chunk.
   A chunk agrees on one direction if it is in ascending/descending
   order. For instance, `[1,2,3,4]` agrees on the +1 direction, and `[4,3,2,1]`
   agrees on -1, but `[1,2,1,0]` is neither so we treat it as if it
   agrees on 0.

3. Sum all of the directions. If this sum is < 0 then we treat it as
   if the user does direction -1 (left/up), otherwise we guess direction
   +1 (right/down).

That's just for one component. To decide whether the overall movement is left/right/up/down,
we calculate the sum of directions for X and Y, and pick whichever one which has the larger
absolute value. E.g. if sum for X = -10, and sum for Y = 35, we pick `down`.
