<?php
function base_path($path = '')
{
    return realpath(getcwd().'/../' . $path);
}

function public_path($path = '')
{
    return base_path('public/' . $path);
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
