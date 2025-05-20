# Bassit

`Bassit` is the bass in the terminal

## Prerequisites

### MacOS

EN: Grant execute permission to `Rubber Band Library` binary file

```sh
cd /path/to/bassit
chmod +x assets/3rdparty/rubberband-4.0.0-gpl-executable-macos/rubberband*
# Check if the permission is `-rwxr-xr-x`
ls -l assets/3rdparty/rubberband-4.0.0-gpl-executable-macos/rubberband*
```

## 3rd Party Libraries & Assets

- [Rubber Band Library](https://breakfastquay.com/rubberband/)
    - License: [GPLv2](https://github.com/breakfastquay/rubberband/blob/default/COPYING)
    - It is used for generating pitch-shifted audio samples.
- [A bass sound sample from Pixabay](https://pixabay.com/sound-effects/bass-guitar-c2-raw-101274/)
    - License: [Pixabay Content License](https://pixabay.com/service/license-summary/)
    - It is trimmed and pitch-shifted to create the bass sound.
