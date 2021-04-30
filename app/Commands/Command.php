<?php

namespace App\Commands;

use Codedungeon\PHPCliColors\Color;

class Command {

    protected $signature;
    protected $description;

    public function info($message)
    {
        echo $this->line(Color::GREEN . $message);
    }

    public function line($message)
    {
        echo $message . PHP_EOL;
    }

    public function signature()
    {
        return $this->signature;
    }

    public function handle()
    {
        // 
    }

    public function error($message)
    {
        return $this->line($message);
    }
}