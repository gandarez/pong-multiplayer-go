#!/bin/bash

install_ubuntu() {
  echo "Detected Ubuntu. Installing dependencies..."
  sudo apt-get update
  sudo apt-get install -y build-essential gcc-multilib libc6-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config
}

install_fedora() {
  echo "Detected Fedora. Installing dependencies..."
  sudo dnf install -y @development-tools mesa-libGL-devel mesa-libGLES-devel libXrandr-devel libXcursor-devel libXinerama-devel libXi-devel libXxf86vm-devel alsa-lib-devel pkg-config
}

install_arch() {
  echo "Detected Arch Linux. Installing dependencies..."
  sudo pacman -Sy --needed base-devel mesa libxrandr libxcursor libxinerama libxi pkg-config
}

install_opensuse() {
  echo "Detected openSUSE. Installing dependencies..."
  sudo zypper install -t pattern devel_basis
  sudo zypper install libGL-devel libXrandr-devel libXinerama-devel libXcursor-devel libXi-devel alsa-devel libXxf86vm1 libXxf86vm-devel
}

install_xcode_select() {
  echo "Detected macOS. Installing xcode-select..."
  if ! xcode-select --print-path &> /dev/null; then
    echo "xcode-select not found. Installing..."
    xcode-select --install
    if [ $? -ne 0 ]; then
      echo "Failed to install xcode-select. Please install it manually."
      exit 1
    fi
  else
    echo "xcode-select is already installed."
  fi
}

if [ "$(uname)" == "Darwin" ]; then
  install_xcode_select
elif [ -f /etc/os-release ]; then
  . /etc/os-release
  case "$ID" in
    ubuntu)
      install_ubuntu
      ;;
    fedora)
      install_fedora
      ;;
    arch)
      install_arch
      ;;
    opensuse*)
      install_opensuse
      ;;
  esac
else
  echo "Couldn't detect your operating system or install Ebiten dependencies automatically"
  echo "For Windows and other Linux distributions, please refer to: https://ebitengine.org/en/documents/install.html to download the dependencies."
  exit 1
fi
