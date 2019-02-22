#!/bin/bash
while IFS= read -r line; do
    case "$line" in
        #3.up)    herbstclient focus up    ;;
        #3.down)  herbstclient focus down  ;;
        #3.left)  herbstclient focus left  ;;
        #3.right) herbstclient focus right ;;
        #4.up)    st -e "htop" & ;;
        #4.up)    herbstclient close ;;
        #4.down)  firefox & ;;
        4.left)  herbstclient cycle_monitor -1 ;;
        4.right) herbstclient cycle_monitor +1 ;;
        #4.left)  herbstclient use_index -1 ;;
        #4.right) herbstclient use_index +1 ;;
    esac
done
