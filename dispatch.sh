#!/bin/bash
while IFS= read -r line; do
    case "$line" in
        3.up)    herbstclient focus up    ;;
        3.down)  herbstclient focus down  ;;
        3.left)  herbstclient focus left  ;;
        3.right) herbstclient focus right ;;
        # 4.up
        # 4.down
        # 4.left
        # 4.right
    esac
done
