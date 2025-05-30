# 𝄢 bassit - bass in terminal

[![Go](https://github.com/Golevka2001/bassit/actions/workflows/go.yml/badge.svg)](https://github.com/Golevka2001/bassit/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Golevka2001/bassit)](https://goreportcard.com/report/github.com/Golevka2001/bassit)
[![go.dev reference](https://godoc.org/github.com/Golevka2001/bassit?status.svg)](https://pkg.go.dev/github.com/Golevka2001/bassit)

Tired of typing commands? Time to play some bass 🎸

`bassit` is a terminal-based bass guitar simulator written in Go. It allows you to play bass lines using your keyboard.

Supported platforms include **Linux**, **macOS**, and **Windows**.

## :rocket: Quick Start

### :building_construction: Install Dependencies

See the [prerequisites](#building_construction-prerequisites) section below for instructions on installing dependencies.

### :computer: Install and Run

```sh
# Ensure Go>=1.23 is installed
go version

# Install bassit
go install github.com/Golevka2001/bassit@latest

# Run
./bassit
```

### :keyboard: Keybindings

![Keybindings](./README.assets/keybindings.png)

#### Basic Controls

- `Esc`: Exit the program

#### Fret Controls

Press one of the following keys to **press a fret** on a specific string:

|    Key    | String |    Range     |
| :-------: | :----: | :----------: |
| `1` ~ `-` |   G    | fret 1 to 11 |
| `q` ~ `p` |   D    | fret 1 to 10 |
| `a` ~ `l` |   A    | fret 1 to 9  |
| `z` ~ `,` |   E    | fret 1 to 8  |

#### Pluck Controls

Press one of the following keys to **pluck** the corresponding string (make the note sound):

|       Key        | String |
| :--------------: | :----: |
| `=`, `Backspace` |   G    |
|     `[`, `]`     |   D    |
|     `;`, `'`     |   A    |
|     `.`, `/`     |   E    |

> 💡 **Tip:** To play a note, first press a fret key, then pluck the corresponding string.

## :building_construction: Prerequisites

### :penguin: Linux

<details>
<summary><strong>TL;DR</strong></summary>

- On **Debian-based** distributions, the following packages are required:
  - `libasound2-dev` (required by [Oto](https://github.com/ebitengine/oto?tab=readme-ov-file#linux))
  - `libx11-dev xorg-dev libxtst-dev xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev libxkbcommon-dev` (required by [GoHook](https://github.com/robotn/gohook?tab=readme-ov-file#requirements-linux))
  - `rubberband-cli`

- On **RedHat-based** distributions, the following packages are required:
  - `alsa-lib-devel` (required by [Oto](https://github.com/ebitengine/oto?tab=readme-ov-file#linux))
  - `libXtst-devel libxkbcommon-devel libxkbcommon-x11-devel xorg-x11-xkb-utils-devel` (required by [GoHook](https://github.com/robotn/gohook?tab=readme-ov-file#requirements-linux))
  - `rubberband` *(except for OpenSUSE, which uses `rubberband-cli`)*

</details>
</br>

**Install `alsa` (required by Oto):**

See [Oto's README](https://github.com/ebitengine/oto?tab=readme-ov-file#linux).

**Install `x11` related packages (required by GoHook):**

See [RobotGo's README](https://github.com/go-vgo/robotgo?tab=readme-ov-file#requirements).

**Install `rubberband-cli` or `rubberband`:**

Debian, Ubuntu, and openSUSE users can install `rubberband-cli`:

```sh
apt install rubberband-cli
```

Other distributions can install `rubberband`:

```sh
apk add rubberband
# OR dnf install rubberband
# OR pacman -S rubberband
# etc.
```

If [pkgs.org](https://pkgs.org/download/rubberband) doesn't list your Linux distribution, you will need to install from source, see [their repository](https://github.com/breakfastquay/rubberband).

## :toolbox: 3rd Party Libraries & Assets

- [Rubber Band Library](https://breakfastquay.com/rubberband/)
  - License: [GPLv2](https://github.com/breakfastquay/rubberband/blob/default/COPYING)
  - It is used for generating pitch-shifted audio samples.
- [A bass sound sample from Pixabay](https://pixabay.com/sound-effects/bass-guitar-c2-raw-101274/)
  - License: [Pixabay Content License](https://pixabay.com/service/license-summary/)
  - It is trimmed and pitch-shifted to create the bass sound.
