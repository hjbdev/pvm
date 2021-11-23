# PVM for Windows

Removing the hassle of changing PHP versions in the CLI on Windows.

This package has a much more niche use case than nvm does. When developing on Windows and using the integrated terminal, it's quite difficult to get those terminals to _actually_ listen to PATH changes.

This utility changes that.

## Installation

The installation process is a little convoluted. You need PHP already installed for it to work, which I admit isn't ideal.

[Packagist](https://packagist.org/packages/hjbdev/pvm)

```
composer global require hjbdev/pvm
```

Type `pvm discover`, then copy the path from `pvm path` and paste it into your path in the Windows Environment variables.

## Commands

```
pvm discover <path?>
```
The path variable is optional, by default it will go to the Laragon bin folder (C:\laragon\bin\php). Provide a path to where all your PHP installations are.

```
pvm list
```
Will list out all the available PHP versions you have

```
pvm path
```
Will tell you what to put in your Path variable.

```
pvm use 7.1
```
Will switch your currently active PHP version to PHP 7.1

```
pvm clear
```
Clears all the detected php versions from `pvm discover`
