#!/bin/bash

# Create new tmux session
tmux new-session -d -s blockchain -x 80 -y 24  # Adjusted for M1 screen

# Main window layout (75% top, 25% bottom)
tmux send-keys -t blockchain "go run mainNode/main.go" C-m
tmux split-window -v -p 25  # Smaller bottom pane

# Bottom pane split (50/50 left/right)
tmux send-keys "sleep 2 && go run validator/validator.go" C-m
tmux split-window -h -p 50

# Right side split (top/bottom)
tmux send-keys "go run mainNode/tcp/main_tcp.go" C-m
tmux select-pane -U
tmux split-window -v -p 50
tmux send-keys "go run producer/main.go" C-m

# Final layout adjustment
tmux select-layout -t blockchain main-vertical
tmux set-window-option -t blockchain synchronize-panes off

# Attach to session with proper sizing
tmux attach-session -t blockchain