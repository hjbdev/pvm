<?php

use App\Application;
use App\Commands\ClearCommand;
use App\Commands\DiscoverCommand;
use App\Commands\ListCommand;
use App\Commands\PathCommand;
use App\Commands\UseCommand;

include 'helpers.php';

// Import all our commands and map them
// to an assoc array.

$app = app();

$app->handle();

// $commands = [
//     ClearCommand::class,
//     DiscoverCommand::class,
//     ListCommand::class,
//     PathCommand::class,
//     UseCommand::class
// ];


// $commandMap = [];

// foreach ($commands as $command) {
//     $command = new $command();
//     $commandMap[$command->signature()] = $command;
// }



