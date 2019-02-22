package main

import "os"
import "fmt"
import "time"
import "github.com/gvalkov/golang-evdev"

func gatherEvents(device *evdev.InputDevice) chan *evdev.InputEvent {
	sink := make(chan *evdev.InputEvent)
	go func() {
		for {
			ev, err := device.ReadOne()
			if err != nil {
				sink <- nil
				close(sink)
				break
			}
			sink <- ev
		}
	}()
	return sink
}

func update_buff(buff *[4]int32, value int32, i int) (int, int) {
	if i < 4 {
		buff[i] = value
		return 0, i + 1
	}
	// shift values and recalculate
	buff[0] = buff[1]
	buff[1] = buff[2]
	buff[2] = buff[3]
	buff[3] = value
	lt := buff[1] <= buff[0] && buff[2] <= buff[1] && buff[3] <= buff[2]
	gt := buff[1] >= buff[0] && buff[2] >= buff[1] && buff[3] >= buff[2]
	if gt && !lt {
		return +1, i
	} else if !gt && lt {
		return -1, i
	} else {
		return 0, i
	}
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func guess_event(dx, dy int, c_max uint16) (event string, ok bool) {
	switch c_max {
	case evdev.BTN_TOOL_QUADTAP:
		event = "4."
	case evdev.BTN_TOOL_TRIPLETAP:
		event = "3."
	default:
		ok = false
		return
	}
	ok = true
	if abs(dx) > abs(dy) {
		if dx < 0 {
			event += "left"
		} else {
			event += "right"
		}
	} else {
		if dy < 0 {
			event += "up"
		} else {
			event += "down"
		}
	}
	return
}

func watch(device *evdev.InputDevice) {
	events := gatherEvents(device)
	i := 0
	j := 0
	ddx := 0
	ddy := 0
	dx := 0
	dy := 0
	x_buf := [4]int32{}
	y_buf := [4]int32{}
	c_max := uint16(0)
	for {
		select {
		case ev := <-events:
			if ev == nil {
				break
			}
			switch ev.Type {
			case evdev.EV_ABS:
				switch ev.Code {
				case evdev.ABS_X:
					ddx, i = update_buff(&x_buf, ev.Value, i)
					dx += ddx
				case evdev.ABS_Y:
					ddy, j = update_buff(&y_buf, ev.Value, j)
					dy += ddy
				}
			case evdev.EV_KEY:
				if ev.Code == evdev.BTN_TOOL_QUADTAP || ev.Code == evdev.BTN_TOOL_TRIPLETAP {
					if ev.Code > c_max {
						c_max = ev.Code
					}
				}
			}
		case <-time.After(100 * time.Millisecond):
			// do some processing
			event, ok := guess_event(dx, dy, c_max)
			if ok {
				fmt.Println(event)
			}
			dx = 0
			dy = 0
			i = 0
			j = 0
			c_max = 0
		}
	}
}

func main() {
	devices, _ := evdev.ListInputDevices()
	if len(os.Args) == 1 {
		fmt.Println(`usage:
	midas list
	midas <name>
	midas <path>`)
		os.Exit(1)
	}
	if os.Args[1] == "list" {
		for _, device := range devices {
			fmt.Println(device.File.Name(), device.Name)
		}
	} else if len(os.Args[1]) > 0 && os.Args[1][0] == '/' {
		device, err := evdev.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		watch(device)
	} else {
		for _, device := range devices {
			if device.Name == os.Args[1] {
				watch(device)
				break
			}
		}
	}
}
