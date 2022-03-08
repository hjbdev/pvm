<?php

use App\Application;
use App\Support\Collection;

if (!function_exists('base_path')) {
    function base_path($path = '')
    {
        return realpath(__DIR__ . DIRECTORY_SEPARATOR . '..'  . DIRECTORY_SEPARATOR . $path);
    }
}

if (!function_exists('public_path')) {
    function public_path($path = '')
    {
        return base_path('public' . DIRECTORY_SEPARATOR . $path);
    }
}

if (!function_exists('starts_with')) {
    function starts_with($haystack, $needles)
    {
        foreach ((array) $needles as $needle) {
            if ((string) $needle !== '' && strncmp($haystack, $needle, strlen($needle)) === 0) {
                return true;
            }
        }

        return false;
    }
}

if (!function_exists('storage_path')) {
    function storage_path($path = '')
    {
        return base_path('storage' . DIRECTORY_SEPARATOR . $path);
    }
}

if (!function_exists('ends_with')) {
    function ends_with($haystack, $needles)
    {
        foreach ((array) $needles as $needle) {
            if ($needle !== '' && substr($haystack, -strlen($needle)) === (string) $needle) {
                return true;
            }
        }

        return false;
    }
}

if (!function_exists('url')) {
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
}

if (!function_exists('collect')) {
    function collect($arr = [])
    {
        return new Collection($arr);
    }
}

if (!function_exists('dump')) {
    function dump($data)
    {
        var_dump($data);
    }
}

if (!function_exists('dd')) {
    function dd($data)
    {
        dump($data);
        die;
    }
}

function pvm_app()
{
    return Application::getInstance();
}

if (!function_exists('windows_os')) {
    function windows_os()
    {
        // taken from here https://github.com/laravel/framework/pull/30660/commits/eb35913866323963aa2086e67f661ce1e0f81b97
        return strtolower(substr(PHP_OS, 0, 3)) === 'win';
    }
}

if (!function_exists('rmlink')) {
    function rmlink($path)
    {
        if (windows_os()) {
            exec("cmd /c rmdir " . escapeshellarg($path));
        } else {
            unlink($path);
        }
    }
}
