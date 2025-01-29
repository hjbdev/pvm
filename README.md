# PVM for Windows

> [!TIP]
> Looking for the 0.x (composer) version? See the [v0 branch](https://github.com/hjbdev/pvm/tree/v0).

Removing the hassle of changing PHP versions in the CLI on Windows.

This package has a much more niche use case than nvm does. When developing on Windows and using the integrated terminal, it's quite difficult to get those terminals to _actually_ listen to PATH changes.

This utility changes that.

## Installation

> [!WARNING]
> version lower than 1.3.0 will have only pvm.exe
> version 1.3.0 or higher will include pvm-setup.exe but can still get pvm.exe from source

### Installer
> Download the latest pvm installer from the releases page (>= 1.3.0).

### Manual Installation
> Create the folder `%UserProfile%\.pvm\bin` (e.g. `C:\Users\Harry\.pvm\bin`) and drop the pvm.exe in there. Add the folder to your PATH.

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

## Composer support

`pvm` now installs also composer with each php version installed.
It will install Composer latest stable release for PHP >= 7.2 and Composer latest 2.2.x LTS for PHP < 7.2.
You'll be able to invoke composer from terminal as it is intended:
```shell
composer --version
```

## Build this project

To compile this project use:

```shell
GOOS=windows GOARCH=amd64 go build -o pvm.exe
```
