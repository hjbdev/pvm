# PVM for Windows

Removing the hassle of changing PHP versions in the CLI on Windows.

This package has a much more niche use case than nvm does. When developing on Windows and using the integrated terminal, it's quite difficult to get those terminals to _actually_ listen to PATH changes.

This utility changes that.

## Installation

Download the latest pvm version from the releases page.

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
pvm use 8.2
```
> [!NOTE]  
> Versions must be specified exactly in the *use* command.

Will switch your currently active PHP version to PHP 8.2

```
pvm install 8.2
```
> [!NOTE]  
> The install command will automatically determine the newest minor/patch versions if they are not specified

Will install PHP 8.2 at the latest patch.