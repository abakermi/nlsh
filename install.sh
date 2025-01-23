#!/bin/bash

# Compile the Go program
go build -o nlsh

# Create bin directory if it doesn't exist
mkdir -p ~/bin

# Move the binary to bin directory
mv nlsh ~/bin/

# Add to PATH if not already added
if ! grep -q "export PATH=\$PATH:\$HOME/bin" ~/.zshrc 2>/dev/null && ! grep -q "export PATH=\$PATH:\$HOME/bin" ~/.bashrc 2>/dev/null; then
    if [ -f "$HOME/.zshrc" ]; then
        echo 'export PATH=$PATH:$HOME/bin' >>~/.zshrc
        echo "Added ~/bin to PATH in .zshrc"
    elif [ -f "$HOME/.bashrc" ]; then
        echo 'export PATH=$PATH:$HOME/bin' >>~/.bashrc
        echo "Added ~/bin to PATH in .bashrc"
    fi
fi

echo "Installation complete! Please restart your terminal or run: source ~/.zshrc (or source ~/.bashrc)"
echo "You can now use 'nlsh' command from anywhere in your terminal"
