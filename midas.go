package main

import "os"
import "fmt"
import "time"
import "github.com/desertbit/timer"
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

func update_minmax(value int32, min, max int32) (int32, int32) {
	if min == 0 || value < min {
		min = value
	}
	if value > max {
		max = value
	}
	return min, max
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
	asc := buff[0] <= buff[1] && buff[1] <= buff[2] && buff[2] <= buff[3]
	dsc := buff[0] >= buff[1] && buff[1] >= buff[2] && buff[2] >= buff[3]
	if asc && !dsc {
		return +1, i
	} else if !asc && dsc {
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

func guess_event(min_x, min_y, max_x, max_y int32, dx, dy int, c_max uint16) (event string, ok bool) {
	switch c_max {
	case evdev.BTN_TOOL_QUADTAP:
		event = "4."
	case evdev.BTN_TOOL_TRIPLETAP:
		event = "3."
	default:
		ok = false
		return
	}
	is_x := abs(dx) > abs(dy)
	if is_x && max_x-min_x > 100 {
		ok = true
		if dx < 0 {
			event += "left"
		} else {
			event += "right"
		}
	} else if !is_x && max_y-min_y > 100 {
		ok = true
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
	diff := 0
	i, j := 0, 0
	dy, dx := 0, 0
	min_y, max_y := int32(0), int32(0)
	min_x, max_x := int32(0), int32(0)
	x_buf := [4]int32{}
	y_buf := [4]int32{}
	c_max := uint16(0)
	duration := 100 * time.Millisecond
	t := timer.NewTimer(duration)
	for {
		select {
		case ev := <-events:
			if ev == nil {
				break
			}
			t.Reset(duration)
			switch ev.Type {
			case evdev.EV_ABS:
				switch ev.Code {
				case evdev.ABS_X:
					diff, i = update_buff(&x_buf, ev.Value, i)
					min_x, max_x = update_minmax(ev.Value, min_x, max_x)
					dx += diff
				case evdev.ABS_Y:
					diff, j = update_buff(&y_buf, ev.Value, j)
					min_y, max_y = update_minmax(ev.Value, min_y, max_y)
					dy += diff
				}
			case evdev.EV_KEY:
				if ev.Code == evdev.BTN_TOOL_QUADTAP || ev.Code == evdev.BTN_TOOL_TRIPLETAP {
					if ev.Code > c_max {
						c_max = ev.Code
					}
				}
			}
		case <-t.C:
			event, ok := guess_event(min_x, min_y, max_x, max_y, dx, dy, c_max)
			if ok {
				fmt.Println(event)
			}
			i, j = 0, 0
			dx, dy = 0, 0
			min_y, max_y = 0, 0
			min_x, max_x = 0, 0
			c_max = 0
		}
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println(`usage:
	midas list
	midas <name>
	midas <path>`)
		os.Exit(1)
	}
	if os.Args[1] == "list" {
		devices, _ := evdev.ListInputDevices()
		for _, device := range devices {
			fmt.Println(device.File.Name(), device.Name)
		}
	} else if len(os.Args[1]) > 0 && os.Args[1][0] == '/' {
		device, err := evdev.Open(os.Args[1])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "midas: %e\n", err)
		}
		watch(device)
	} else {
		devices, _ := evdev.ListInputDevices()
		for _, device := range devices {
			if device.Name == os.Args[1] {
				watch(device)
				break
			}
		}
	}
}
