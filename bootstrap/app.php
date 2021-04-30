<?php

use App\Commands\ClearCommand;
use App\Commands\DiscoverCommand;
use App\Commands\ListCommand;
use App\Commands\PathCommand;
use App\Commands\UseCommand;

include 'helpers.php';

// Import all our commands and map them
// to an assoc array.

$commands = [
    ClearCommand::class,
    DiscoverCommand::class,
    ListCommand::class,
    PathCommand::class,
    UseCommand::class
];

$commandMap = [];

foreach ($commands as $command) {
    $command = new $command();
    $commandMap[$command->signature()] = $command;
}

// Get arguments from command line input

if ($argc === 1) {
    // help command
}

if ($argc > 1) {
    $argument = $argv[1];

    if(isset($commandMap[$argv[1]])) {
        $commandMap[$argv[1]]->handle();
    }
}
