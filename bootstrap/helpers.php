<?php

use App\Application;
use App\Support\Collection;

function base_path($path = '')
{
    return realpath(__DIR__ . DIRECTORY_SEPARATOR . '..'  . DIRECTORY_SEPARATOR . $path);
}

function public_path($path = '')
{
    return base_path('public' . DIRECTORY_SEPARATOR . $path);
}

function starts_with($haystack, $needles)
{
    foreach ((array) $needles as $needle) {
        if ((string) $needle !== '' && strncmp($haystack, $needle, strlen($needle)) === 0) {
            return true;
        }
    }

    return false;
}

function storage_path($path = '')
{
    return base_path('storage' . DIRECTORY_SEPARATOR . $path);
}

function ends_with($haystack, $needles)
{
    foreach ((array) $needles as $needle) {
        if ($needle !== '' && substr($haystack, -strlen($needle)) === (string) $needle) {
            return true;
        }
    }

    return false;
}

function url($path = '')
{
    if (isset($_SERVER['HTTPS'])) {
        $protocol = ($_SERVER['HTTPS'] && $_SERVER['HTTPS'] != "off") ? "https" : "http";
    } else {
        $protocol = 'http';
    }

    $path = starts_with($path, '/') ? $path : '/' . $path;

    return $protocol . "://" . $_SERVER['HTTP_HOST'] . parse_url($_SERVER["REQUEST_URI"], PHP_URL_PATH) . $path;
}

function collect($arr = [])
{
    return new Collection($arr);
}

function dump($data)
{
    var_dump($data);
}

function dd($data)
{
    dump($data);
    die;
}

function app()
{
    return Application::getInstance();
}

function windows_os()
{
    // taken from here https://github.com/laravel/framework/pull/30660/commits/eb35913866323963aa2086e67f661ce1e0f81b97
    return strtolower(substr(PHP_OS, 0, 3)) === 'win';
}

function rmlink($path)
{
    if (windows_os()) {
        exec("cmd /c rmdir " . escapeshellarg($path));
    } else {
        unlink($path);
    }
}
