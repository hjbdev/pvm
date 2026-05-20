# PVM for Windows

[Support this project](https://github.com/sponsors/hjbdev)

Removing the hassle of changing PHP versions in the CLI on Windows.

This package has a much more niche use case than nvm does. When developing on Windows and using the integrated terminal, it's quite difficult to get those terminals to _actually_ listen to PATH changes.

This utility changes that.

## Installation

`pvm` currently supports x64 Windows only.

```powershell
irm https://pvm.hjb.dev/install.ps1 | iex
```

You can still install manually by downloading the latest `pvm.exe` release, placing it in `%UserProfile%\.pvm\bin` (for example `C:\Users\Harry\.pvm\bin`), and adding that folder to your PATH.

## Commands
```
pvm list
```
Will list out all the available PHP versions you have installed

```
pvm list remote
```
Will list the PHP versions available for installation.

```
pvm bin
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

```
pvm extensions list
```
Will show regular and Zend extensions for the active PHP version, including whether each extension is enabled, disabled, available in `ext`, or missing from disk.

```
pvm extensions enable curl,openssl
```
Will enable one or more extensions that already have entries in the active version's `php.ini`.

```
pvm extensions disable xdebug
```
Will disable an extension or Zend extension in the active version's `php.ini`.

```
pvm update [--yes-update|-y]
```
Check for a newer `pvm` release. When a newer release is found, `pvm update` will run the installer automatically (no interactive prompt).

- `--yes-update`, `-y`: when provided, `pvm update` prefers the safe download-and-replace installer path (the default non-interactive flow used in CI/testing).

The installer executed by the quick-install command is:

```powershell
irm https://pvm.hjb.dev/install.ps1 | iex
```

Notes:

- `PVM_INSTALL_SCRIPT`: optional environment variable pointing to a custom installer. Can be a URL (http/https) or a local PowerShell script path. When set, `pvm update` will run this installer instead of the default download flow.
- `PVM_INSTALL_CHECKSUM`: optional SHA256 hex string used to verify the downloaded `pvm.exe` when using the automatic download path.
- Safe update behavior: by default `pvm update` will download the `pvm.exe` release asset, optionally verify its checksum, back up the existing `pvm.exe` to `pvm.exe.bak`, atomically replace the binary and verify the installed version. On verification failure it will attempt to roll back to the backup.

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
bash ./build.sh
```

To override the embedded version for a release-style local build:
```shell
VERSION=1.2.1 bash ./build.sh
```

GitHub releases are built automatically from pushed tags and publish both `pvm.exe` and `install.ps1` as release assets.
