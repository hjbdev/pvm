<?php

namespace App\Commands;

use Illuminate\Support\Facades\Storage;

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
        $this->line('📜 Available PHP Versions');
        $this->line('');
 
        if (Storage::has('versions.json')) {
            $versions = collect(json_decode(file_get_contents(storage_path('versions.json'))));
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
