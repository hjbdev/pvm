<?php

namespace App\Commands;

use App\Support\Filesystem;

class ListCommand extends Command
{
    /**
     * The signature of the command.
     *
     * @var string
     */
    protected $signature = 'list';

    /**
     * The description of the command.
     *
     * @var string
     */
    protected $description = 'List all installed versions of PHP';

    /**
     * Execute the console command.
     *
     * @return mixed
     */
    public function handle()
    {
        $this->info('📜 Available PHP Versions');
        $this->line('');
 
        if (Filesystem::has(storage_path('versions.json'))) {
            $versions = collect(json_decode(Filesystem::get(storage_path('versions.json'))));
        } else {
            $versions = collect();
        }

        if(!$versions->first()) {
            $this->line('❌ No PHP versions found');
        }

        foreach($versions as $version) {
            $this->line('    ' . "{$version->major_version}.{$version->minor_version}.{$version->patch_version}" . ($version->active ? '✔' : ''));
        }
    }
}
