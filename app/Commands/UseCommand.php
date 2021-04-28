<?php

namespace App\Commands;

use Illuminate\Support\Facades\Storage;
use LaravelZero\Framework\Commands\Command;

class UseCommand extends Command
{
    /**
     * The signature of the command.
     *
     * @var string
     */
    protected $signature = 'use {version}';

    /**
     * The description of the command.
     *
     * @var string
     */
    protected $description = 'Use a specific version of PHP';

    /**
     * Execute the console command.
     *
     * @return mixed
     */
    public function handle()
    {
        $versionStr = $this->argument('version');
    
        list($major, $minor, $patch) = array_pad(explode('.', $versionStr), 3, null);

        $version = null;

        if (Storage::has('versions.json')) {
            $versions = collect(json_decode(file_get_contents(storage_path('versions.json'))));
        } else {
            $versions = collect();
        }

        $versions = $versions->where('major_version', $major);
        if($minor) $versions = $versions->where('minor_version', $minor);
        if($patch) $versions = $versions->where('patch_version', $patch);


        if($versions->count() > 1) {
            $this->line('You have more than 1 version of PHP ' . $versionStr . ' installed. Please be more specific');
            return;
        } else if ($versions->count() === 1) {
            $version = $versions->first();
        } else {
            $this->line("You do not have {$versionStr}");
            return;
        }

        $filesystem = $this->laravel->make('files');

        if($filesystem->exists(base_path('bin'))) {
            rmdir(base_path('bin'));
        }

        // get the full versions list back out

        if (Storage::has('versions.json')) {
            $versions = collect(json_decode(file_get_contents(storage_path('versions.json'))));
        } else {
            $versions = collect();
        }

        $versions->map(function($item) use ($version) {
            if($version->path === $item->path) {
                $item->active = true;
            } else {
                $item->active = false;
            }
            return $item;
        });

        Storage::put('versions.json', $versions->toJson());

        $filesystem->link($version->path, base_path('bin'));

        $this->info('✔ Switched PHP version to ' . $versionStr);
    }
}
