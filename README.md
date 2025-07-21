# 𝄢 bassit - bass in terminal

[![Go](https://github.com/Golevka2001/bassit/actions/workflows/go.yml/badge.svg)](https://github.com/Golevka2001/bassit/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Golevka2001/bassit)](https://goreportcard.com/report/github.com/Golevka2001/bassit)
[![go.dev reference](https://godoc.org/github.com/Golevka2001/bassit?status.svg)](https://pkg.go.dev/github.com/Golevka2001/bassit)

Tired of typing commands? Time to play some bass 🎸

`bassit` is a terminal-based bass guitar simulator written in Go. It allows you to play bass lines on your keyboard.

Supported on **Linux**, **macOS**, and **Windows**.

## :camera: Snapshots

![MacOS Snapshots](./README.assets/snapshots/macos.gif)

<details>
<summary><strong>Windows and Linux Snapshots</strong></summary>
<div style="display: flex; gap: 10px;">
  <img src="./README.assets/snapshots/windows.gif" alt="Windows Snapshot" width="49%" />
  <img src="./README.assets/snapshots/linux.gif" alt="Linux Snapshot" width="49%" />
</div>
</details>


## :rocket: Quick Start

### :vertical_traffic_light: Before You Begin

1. See the [Dependencies](#building_construction-dependencies) section below for installation instructions.

2. Please use a modern terminal that supports [progressive keyboard enhancements](https://sw.kovidgoyal.net/kitty/keyboard-protocol/). See [this page](https://github.com/charmbracelet/bubbletea/discussions/1374#:~:text=Keyboard-,Enhancements,-Progressive%20keyboard%20enhancements) for a list of supported terminals.

3. To avoid character encoding or alignment issues, please make sure:
   - You are using the **UTF-8 character set**
     - On Linux or macOS, you can check your current character set settings by running the `locale` command.
     - On Windows, you can check your current character set settings in PowerShell by running `$OutputEncoding`.
   - You are using a **monospace font** (fixed-width font) in your terminal.

### :computer: Install and Run

```sh
# Ensure Go>=1.23 is installed
go version

# Install bassit
go install github.com/Golevka2001/bassit@latestdocsd

# Run
$GOPATH/bin/bassit
# OR if you have the Go bin directory in your PATH
bassit
```

### :keyboard: Keybindings

![Keybindings](./README.assets/keybindings/keybindings.drawio.png)

#### Basic Controls

- `Esc`: Move focus to the tab selector
- `→` / `Tab`: Switch to the next tab
- `←` / `Shift+Tab`: Switch to the previous tab
- `Enter`: Enter the selected tab
- `Ctrl+C`: Exit the application

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

**Modifier Keys:**

- `\` : Slap mode
- `Space` : Mute mode

If you hold the modifier key while plucking, you can use different techniques to play, such as:

|             Key             |       Default        | Slap Mode |    Mute Mode     |
| :-------------------------: | :------------------: | :-------: | :--------------: |
|     `=`, `[`, `;`, `.`      | moderately soft (𝓂𝓅) | slap (S)  |  dead note (X)   |
| `Backspace`, `]` ,`'` , `/` | moderately loud (𝓂𝒻) |  pop (P)  | palm mute (P.M.) |

The priority is: Mute Mode > Slap Mode > Default.
(e.g. if no modifier key is held -> default mode; if both modifier keys are held -> mute mode)

> 💡 **Tip:** To play a note, first press a fret key, then pluck the corresponding string.

### :art: Themes

You can modify the style of the bass by simply modifying a theme file.

Of course, bassit also provides some themes for you to choose from.

|                    default                     |                   maple                    |                       vintage-maple                        |                     rosewood                     |                         fretless-pauferro                          |
| :--------------------------------------------: | :----------------------------------------: | :--------------------------------------------------------: | :----------------------------------------------: | :----------------------------------------------------------------: |
| ![default](./README.assets/themes/default.png) | ![maple](./README.assets/themes/maple.png) | ![vintage-maple](./README.assets/themes/vintage-maple.png) | ![rosewood](./README.assets/themes/rosewood.png) | ![fretless-pauferro](./README.assets/themes/fretless-pauferro.png) |

Use `bassit --theme <theme-name>` to switch to a theme.

The theme files are located at `~/.config/bassit/themes/`.

## :building_construction: Dependencies

### :penguin: Linux

<details>
<summary><strong>TL;DR</strong></summary>

- On **Debian-based** distributions, the following packages are required:
  - `libasound2-dev` (required by [Oto](https://github.com/ebitengine/oto?tab=readme-ov-file#linux))
  - `rubberband-cli`

- On **RedHat-based** distributions, the following packages are required:
  - `alsa-lib-devel` (required by [Oto](https://github.com/ebitengine/oto?tab=readme-ov-file#linux))
  - `rubberband` *(except for OpenSUSE, which uses `rubberband-cli`)*

</details>
</br>

**Install `alsa` (required by Oto):**

See [Oto's README](https://github.com/ebitengine/oto?tab=readme-ov-file#linux).

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

If your Linux distribution is not listed on [pkgs.org](https://pkgs.org/download/rubberband), you may need to build from source, see [their repository](https://github.com/breakfastquay/rubberband).

## :toolbox: Third-Party Libraries & Assets

- [Rubber Band Library](https://breakfastquay.com/rubberband/)
  - License: [GPLv2](https://github.com/breakfastquay/rubberband/blob/default/COPYING)
  - It is used for generating pitch-shifted audio samples.

## Planned Features

- More than 4 strings.
- The settings tab for adjusting preferences and themes.
- Chord -> fingering.
- Metronome.
- Resample and create a full-quality sound pack.
- A website for users to upload and share sound packs, themes, etc.
- Recording and exporting.
- ...
