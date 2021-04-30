<?php

namespace App\Commands;

use App\Support\Filesystem;

class UseCommand extends Command
{
    /**
     * The signature of the command.
     *
     * @var string
     */
    protected $signature = 'use';

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
        $versionStr = null;

        if(app()->argument(1)) {
            $versionStr = app()->argument(1);
        } else {
            $this->error('Please provide a version.');
            exit(1);
        }
    
        list($major, $minor, $patch) = array_pad(explode('.', $versionStr), 3, null);

        $version = null;

        if (Filesystem::has(storage_path('versions.json'))) {
            $versions = collect(json_decode(file_get_contents(storage_path('versions.json'))));
        } else {
            $versions = collect();
        }

        $versions = $versions->where('major_version', $major);
        if($minor) $versions = $versions->where('minor_version', $minor);
        if($patch) $versions = $versions->where('patch_version', $patch);


        if($versions->count() > 1) {
            $this->error('You have more than 1 version of PHP ' . $versionStr . ' installed. Please be more specific');
            exit(1);
        } else if ($versions->count() === 1) {
            $version = $versions->first();
        } else {
            $this->error("You do not have {$versionStr}");
            exit(1);
        }

        if(Filesystem::has(base_path('bin'))) {
            rmlink(base_path() . DIRECTORY_SEPARATOR . 'bin');
        }

        // get the full versions list back out

        if (Filesystem::has(storage_path('versions.json'))) {
            $versions = collect(json_decode(Filesystem::get(storage_path('versions.json'))));
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

        Filesystem::put('versions.json', $versions->toJson());

        Filesystem::link($version->path, base_path() . DIRECTORY_SEPARATOR . 'bin');

        $this->info('✔ Switched PHP version to ' . $versionStr);
    }
}
