<?php

namespace App\Traits;

trait Singleton
{
    private static $instance = null;

    private function __construct()
    {
        // ... 
    }

    public static function getInstance()
    {
        if (self::$instance == null) {
            self::$instance = new static();
        }

        return self::$instance;
    }
}
