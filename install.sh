#!/bin/bash

# Function to detect OS and architecture
detect_platform() {
    OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
    ARCH="$(uname -m)"
    case "${ARCH}" in
    x86_64) ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    armv7l) ARCH="arm" ;;
    esac
    echo "${OS}-${ARCH}"
}

# Function to download and verify the binary
download_binary() {
    local platform=$1
    local temp_dir=$(mktemp -d)

    # Fetch latest release tag from GitHub API
    echo "Fetching latest release version..."
    local version=$(curl -s https://api.github.com/repos/abakermi/nlsh/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$version" ]; then
        echo "Error: Failed to fetch latest version"
        rm -rf "${temp_dir}"
        return 1
    fi

    local binary_url="https://github.com/abakermi/nlsh/releases/download/${version}/nlsh-${platform}"
    echo $binary_url
    echo "Downloading nlsh for ${platform}..."
    if ! curl -fsSL "${binary_url}" -o "${temp_dir}/nlsh"; then
        echo "Error: Failed to download binary for ${platform}"
        rm -rf "${temp_dir}"
        return 1
    fi

    chmod +x "${temp_dir}/nlsh"
    mkdir -p ~/bin
    mv "${temp_dir}/nlsh" ~/bin/
    rm -rf "${temp_dir}"
    return 0
}

# Detect platform
PLATFORM=$(detect_platform)

# Download and install binary
if ! download_binary "${PLATFORM}"; then
    echo "Error: Installation failed"
    exit 1
fi

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
