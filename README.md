# 𝄢 bassit - bass in terminal

Tired of typing commands? Time to play some bass 🎸

`bassit` is a terminal-based bass guitar simulator written in Go. It allows you to play bass lines using your keyboard.

Supported platforms include **Linux**, **macOS**, and **Windows**.

## :rocket: Quick Start

### :computer: Install

```sh
# Ensure Go>=1.23 is installed
go version

# Install bassit
go install github.com/Golevka2001/bassit@latest

# Install dependencies
# See prerequisites section below

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

## :toolbox: 3rd Party Libraries & Assets

- [Rubber Band Library](https://breakfastquay.com/rubberband/)
  - License: [GPLv2](https://github.com/breakfastquay/rubberband/blob/default/COPYING)
  - It is used for generating pitch-shifted audio samples.
- [A bass sound sample from Pixabay](https://pixabay.com/sound-effects/bass-guitar-c2-raw-101274/)
  - License: [Pixabay Content License](https://pixabay.com/service/license-summary/)
  - It is trimmed and pitch-shifted to create the bass sound.
