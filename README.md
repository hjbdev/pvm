# PVM for Windows

> [!TIP]
> Looking for the 0.x (composer) version? See the [v0 branch](https://github.com/hjbdev/pvm/tree/v0).

Removing the hassle of changing PHP versions in the CLI on Windows.

This package has a much more niche use case than nvm does. When developing on Windows and using the integrated terminal, it's quite difficult to get those terminals to _actually_ listen to PATH changes.

This utility changes that.

## Installation

Download the latest pvm version from the releases page (1.0-alpha-1, it's currently a pre-release).

Create the folder `%UserProfile%\.pvm\bin` (e.g. `C:\Users\Harry\.pvm\bin`) and drop the pvm exe in there. Add the folder to your PATH.

## Commands

```
pvm list
```

Will list out all the available PHP versions you have installed

```
pvm install 8
```

> [!NOTE]
> The install command will automatically determine the newest minor/patch versions if they are not specified

Will install PHP 8 at the latest minor and patch.

```
pvm use 8.2
```

> [!NOTE]
> Versions must have major.minor specified in the *use* command. If a .patch version is omitted, newest available patch version is chosen.

Will switch your currently active PHP version to PHP 8.2 latest patch.

```
pvm uninstall 8.2.9
```

> [!NOTE]
> Versions must have major.minor.patch specified in the *uninstall* command. If a .patch version is omitted, it will not uninstalling.

Will uninstall PHP version to PHP 8.2.9

## Build

To compile the program use:

```shell
GOOS=windows GOARCH=amd64 go build -o pvm.exe
```
