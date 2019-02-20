import sys
from functools import partial
from threading import Thread
from queue import Queue, Empty
from evdev import InputDevice, ecodes, list_devices


def guess_direction(X, chunk_size=4):
    dir = 0
    for i in range(len(X) - chunk_size):
        seen_lt = False
        seen_gt = False
        for j in range(chunk_size - 1):
            a = X[i+j]
            b = X[i+j+1]
            seen_lt |= b < a
            seen_gt |= b > a
        dir += seen_gt - seen_lt
    return dir


def guess_event(data):
    if not data:
        return
    y_data = []
    x_data = []
    c_data = []
    for t, y in data:
        if   t == 0: y_data.append(y)
        elif t == 1: x_data.append(y)
        else:        c_data.append(y)
    if not c_data:
        return
    dy = guess_direction(y_data)
    dx = guess_direction(x_data)
    c = max(c_data)
    # find out main direction change
    if abs(dx) > abs(dy):
        dir_name = "left" if dx < 0 else "right"
    else:
        dir_name = "up" if dy < 0 else "down"
    key_name = "4"  if c == ecodes.BTN_TOOL_QUADTAP else "3"
    return "%s.%s" % (key_name, dir_name)


def timed(q):
    buffer = []
    while True:
        try:
            buffer.append(q.get(timeout=0.1))
        except Empty:
            evt = guess_event(buffer)
            if evt is not None:
                print(evt, flush=True)
            buffer.clear()


def main(device_path):
    dev = InputDevice(device_path)
    q = Queue()

    Thread(target=partial(timed, q)).start()
    taps = (ecodes.BTN_TOOL_TRIPLETAP,
            ecodes.BTN_TOOL_QUADTAP)

    for event in dev.read_loop():
        type  = event.type
        code  = event.code
        value = event.value
        if type == ecodes.EV_ABS:
            if code == ecodes.ABS_Y:   q.put((0, value))
            elif code == ecodes.ABS_X: q.put((1, value))
        elif type == ecodes.EV_KEY and code in taps:
            q.put((2, code))


if __name__ == '__main__':
    devices = [InputDevice(path) for path in list_devices()]
    if sys.argv[1] == "list":
        for device in devices:
            print(device.path, device.name)
    elif not sys.argv[1].startswith('/'):
        for device in devices:
            if device.name == sys.argv[1]:
                main(device.path)
                break
    else:
        main(sys.argv[1])
