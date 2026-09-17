# geyes

This is a silly app with two eyes following your mouse pointer
around. It's similar to `xeyes` from X11/X.org found on many Unix and
Linux systems.

`geyes` doesn't have a window border and has a transparent background
with only the two eyes showing and the pupils moving. `Ctrl + q` or
macOS, `Cmd+q`, exits the application.

As the `g` in `geyes` suggests, the app is written in Go and uses the
[Fyne framework](https://fyne.io/) for GUI components.

# Building
```
$ make
```

# Running
```
$ make run
```

# Installing

```
$ make install
```

This installs `geyes` to `~/.local/bin`.

# License

The app is released under GPLv3, see [LICENSE](LICENSE).
