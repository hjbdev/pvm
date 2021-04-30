<?php

namespace App\Support;

class Filesystem
{

    public static function has($file)
    {
        return file_exists($file);
    }

    public static function put($file, $contents)
    {
        return file_put_contents($file, $contents);
    }

    public static function get($file)
    {
        return file_get_contents($file);
    }

    public static function delete($file)
    {
        return unlink($file);
    }

    public static function directories($path)
    {
        $dirs = array_filter(glob($path . DIRECTORY_SEPARATOR . '*'), 'is_dir');
        return $dirs;
    }

    // Symlink code adapted from 
    // https://github.com/laravel/framework/blob/78eb4dabcc03e189620c16f436358d41d31ae11f/src/Illuminate/Filesystem/Filesystem.php#L249

    public static function link($target, $link)
    {
        if (!windows_os()) {
            return symlink($target, $link);
        }

        $mode = is_dir($target) ? 'J' : 'H';

        exec("mklink /{$mode} " . escapeshellarg($link) . ' ' . escapeshellarg($target));
    }

}
