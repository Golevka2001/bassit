# 𝄢 bassit - bass in terminal

Tired of typing commands? Time to play some bass 🎸

`bassit` is a terminal-based bass guitar simulator written in Go. It allows you to play bass lines using your keyboard.

## Quick Start

```sh
# Clone the repository
git clone git@github.com:Golevka2001/bassit.git
cd bassit

# Install dependencies
# See prerequisites section below

# Build
go build

# Run
./bassit
```

## Prerequisites

### Linux

**Install `alsa` (required by `Oto`):**

See [Oto's README](https://github.com/ebitengine/oto?tab=readme-ov-file#linux).

**Install `rubberband-cli` or `rubberband`:**

```sh
# Debian, Ubuntu, openSUSE users should install `rubberband-cli`
apt install rubberband-cli
# Other Linux distributions should install `rubberband`
apk add rubberband
# OR dnf install rubberband
# OR pacman -S rubberband
# OR yum install rubberband
# etc.
```

If [pkgs.org](https://pkgs.org/download/rubberband) doesn't list your Linux distribution, you will need to install from source, see [their repository](https://github.com/breakfastquay/rubberband).

## 3rd Party Libraries & Assets

- [Rubber Band Library](https://breakfastquay.com/rubberband/)
  - License: [GPLv2](https://github.com/breakfastquay/rubberband/blob/default/COPYING)
  - It is used for generating pitch-shifted audio samples.
- [A bass sound sample from Pixabay](https://pixabay.com/sound-effects/bass-guitar-c2-raw-101274/)
  - License: [Pixabay Content License](https://pixabay.com/service/license-summary/)
  - It is trimmed and pitch-shifted to create the bass sound.
