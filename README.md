# PVM for Windows

[Support this project](https://github.com/sponsors/hjbdev)

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
pvm path
```
Will tell you what to put in your Path variable.

```
pvm use 8.2.9
```
> [!NOTE]  
> Versions must have major.minor specified in the *use* command. If a .patch version is omitted, newest available patch version is chosen.

Will switch your currently active PHP version to PHP 8.2.9

```
pvm install 8.2
```
> [!NOTE]  
> The install command will automatically determine the newest minor/patch versions if they are not specified

Will install PHP 8.2 at the latest patch.

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
