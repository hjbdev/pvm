<?php

use App\Support\ExeInfo;
use Codedungeon\PHPCliColors\Color;

test('application bootstraps correctly', function () {
    $output = null;
    $resultCode = null;
    exec('php ' . __DIR__ . '/../pvm', $output, $resultCode);
    expect($resultCode)->toBe(0);
    expect($output)->toBeArray();
    expect($output)->toContain(Color::GREEN . "Available commands:" . Color::RESET);
});


test('list command works', function () {
    $output = null;
    $resultCode = null;
    exec('php ' . __DIR__ . '/../pvm list', $output, $resultCode);
    expect($resultCode)->toBe(0);
    expect($output)->toBeArray();
    expect($output)->toContain(Color::GREEN . '📜 Available PHP Versions' . Color::RESET);
});

test('discover command works', function () {
    $output = null;
    $resultCode = null;
    exec('php ' . __DIR__ . '/../pvm discover', $output, $resultCode);
    expect($resultCode)->toBe(0);
    expect($output)->toBeArray();
    expect($output)->toContain(Color::GREEN . '🔎 Discovering PHP Versions...' . Color::RESET);

    expect(file_get_contents('./storage/versions.json'))->toBeJson();
});

test('use command works', function () {
    $output = null; 
    $resultCode = null;
    exec('php ' . __DIR__ . '/../pvm use 8', $output, $resultCode);
    expect($resultCode)->toBe(0);
    $version = ExeInfo::getFileVersion('./bin/php.exe');
    expect($version)->toBe('8.0.5');
});