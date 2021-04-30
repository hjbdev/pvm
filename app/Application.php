<?php

namespace App;

use App\Traits\Singleton;

class Application {

    use Singleton;

    protected $argv;
    protected $commands = [];
    protected $arguments;

    private function __construct() {
        
        $this->argv = $_SERVER['argv'];

        $this->discoverCommands();
        $this->parseArguments();

        
        

        // foreach ($commands as $command) {
        //     $command = new $command();

        //     echo $command->signature() . PHP_EOL;

        //     if($arguments[$command->signature()]) {
        //         $command->handle();
        //     }
        // }
    }

    public function handle()
    {
        if ($this->arguments->count() === 0) {
            // help command
        }

        if ($this->arguments->count() > 0) {
            $argument = $this->arguments->first();

            if(isset($this->commands[$this->arguments->first()])) {
                $this->commands[$this->arguments->first()]->handle();
            }
        }
    }

    protected function parseArguments()
    {
        $arguments = collect($this->argv);

        $arguments->shift(); // removes pvm

        // filter out parameters
        $arguments = $arguments->filter(function ($item) {
            return !starts_with($item, '-');
        });

        $this->arguments = $arguments;
    }

    public function argument($key)
    {
        return $this->arguments->get($key);
    }

    protected function discoverCommands()
    {
        $commands = collect(scandir(base_path('app/Commands')));

        // Get rid of all non-PHP files and ignore the base command
        $commands = $commands->filter(function($item) {
            return ends_with($item, '.php') && $item !== 'Command.php';
        })->map(function($item) {
            return str_replace('.php', '', $item);
        });

        // Instantiate classes and feed into commands property
        foreach($commands as $command) {
            $class = 'App\\Commands\\' . $command;
            $c = new $class();
            $this->commands[$c->signature()] = $c;
        }
    }
}