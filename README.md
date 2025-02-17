# PVM for Windows

> [!TIP]
> Looking for the 0.x (composer) version? See the [v0 branch](https://github.com/hjbdev/pvm/tree/v0).

Removing the hassle of changing PHP versions in the CLI on Windows.

This package has a much more niche use case than nvm does. When developing on Windows and using the integrated terminal, it's quite difficult to get those terminals to _actually_ listen to PATH changes.

This utility changes that.

---

## Installation

> [!WARNING]
> version lower than 1.3.0 will have only pvm.exe
> version 1.3.0 or higher will include pvm-setup.exe but can still get pvm.exe from source

### Installer
Download the latest pvm installer from the releases page (>= 1.3.0).

### Manual Installation
Create the folder `%UserProfile%\.pvm\bin` (e.g. `C:\Users\Harry\.pvm\bin`) and drop the pvm.exe in there. Add the folder to your PATH.

---

## Commands
```
pvm list
```
 - __Will list out all the available PHP versions you have installed.__

```
pvm list-remote
```
- __Will list available PHP versions from remote repositories.__

```
pvm install <version> [nts] [path]
```
- __Will install specified PHP vesrion.__
- If the minor and patch versions are not specified, the newest available versions will be automatically selected.
- [nts] (optional): Install a non-thread-safe version.
- [path] (optional): Specify a custom installation path.

```
pvm use <version|path>
```
- __Will switch your currently active PHP version to specified PHP vesrion.__
-   Using a version: Specify at least the `major.minor` version. If the `.patch` version is omitted, the newest available patch version will be selected automatically.
-   Using a path: Specify the path to the PHP executable to set it as the active PHP version.

```
pvm uninstall <version>
```
- __Will uninstall specified PHP vesrion.__
- The `uninstall` command requires the full `major.minor.patch` version to be specified.
- This command can uninstall only php installed by pvm.

```
pvm add <path>
```
- __Adds a custom PHP version by specifying the path to a PHP executable.__

```
pvm remove <path>
```
- __Removes a custom PHP version by specifying the path to the previously added PHP executable.__

---

## Composer support

`pvm` now installs also composer with each php version installed.
It will install Composer latest stable release for PHP >= 7.2 and Composer latest 2.2.x LTS for PHP < 7.2.
You'll be able to invoke composer from terminal as it is intended:
```shell
composer --version
```

---

## Build this project

To compile this project use:

```shell
GOOS=windows GOARCH=amd64 go build -o pvm.exe
```

To build pvm-setup.exe use:

```shell
iscc "pvm-setup.iss" 
```
