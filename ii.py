from pprint import pprint
from collections import Counter
from threading import Thread
from queue import Queue, Empty
from evdev import InputDevice, categorize, ecodes

dev = InputDevice("/dev/input/event18")
data = []
q = Queue()


def guess_event(data):
    if data:
        y_data = [y for t,y in data if t == 0]
        c_data = [y for t,y in data if t == 1]
        if not c_data:
            return
        up = 0
        down = 0
        for i in range(len(y_data) - 4):
            chunk = y_data[i:i+4]
            seen_up   = False
            seen_down = False
            for i, a in enumerate(chunk):
                if i < 3:
                    b = chunk[i+1]
                    seen_up   |= b < a
                    seen_down |= b > a
            # if we cannot agree then screw it
            # forget about this chunk
            if seen_up and seen_down: continue
            if seen_up: up += 1
            else:       down += 1
        return max(c_data), 0 if up > down else 1
    return


def timed():
    buffer = []
    while True:
        try:
            buffer.append(q.get(timeout=0.1))
        except Empty:
            evt = guess_event(buffer)
            if evt is not None:
                print(evt)
            buffer.clear()


Thread(target=timed).start()


for event in dev.read_loop():
    # event.value
    # 0 = up
    # 1 = down
    # event.code
    # 335 = quad touch
    # 334 = triple touch
    if event.type == ecodes.EV_ABS and event.code == ecodes.ABS_Y:
        q.put_nowait((0, event.value))
    elif event.type == ecodes.EV_KEY and event.code in (334, 335):
        q.put_nowait((1, event.code))
