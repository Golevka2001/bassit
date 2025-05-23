# Bassit

`Bassit` is the bass in the terminal

## Prerequisites

### MacOS

Grant execute permission to `Rubber Band Library` binary file:

```sh
cd /path/to/bassit
chmod +x assets/3rdparty/rubberband-4.0.0-gpl-executable-macos/rubberband*
# Check if the permission is `-rwxr-xr-x`
ls -l assets/3rdparty/rubberband-4.0.0-gpl-executable-macos/rubberband*
```

### Linux

Install `rubberband-cli` or `rubberband`:

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

If [pkgs.org](https://pkgs.org/download/rubberband) doesn't list your Linux distribution, you will need to install from
source, see [their repository](https://github.com/breakfastquay/rubberband).

## 3rd Party Libraries & Assets

- [Rubber Band Library](https://breakfastquay.com/rubberband/)
    - License: [GPLv2](https://github.com/breakfastquay/rubberband/blob/default/COPYING)
    - It is used for generating pitch-shifted audio samples.
- [A bass sound sample from Pixabay](https://pixabay.com/sound-effects/bass-guitar-c2-raw-101274/)
    - License: [Pixabay Content License](https://pixabay.com/service/license-summary/)
    - It is trimmed and pitch-shifted to create the bass sound.
