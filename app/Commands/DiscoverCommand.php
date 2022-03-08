<?php

namespace App\Commands;

use App\Support\ExeInfo;
use App\Support\Filesystem;

class DiscoverCommand extends Command
{
    /**
     * The signature of the command.
     *
     * @var string
     */
    protected $signature = 'discover';

    /**
     * The description of the command.
     *
     * @var string
     */
    protected $description = 'Discovers existing PHP installations on your system.';

    /**
     * Execute the console command.
     *
     * @return mixed
     */
    public function handle()
    {
        $this->info('🔎 Discovering PHP Versions...');
        $this->line('');
        
        $path = pvm_app()->argument(1) ?: 'C:\laragon\bin\php';
        $dirs = Filesystem::directories($path);
        $discovered = 0;

        if (Filesystem::has(storage_path() . '/versions.json')) {
            $versions = collect(json_decode(Filesystem::get(storage_path('versions.json'))));
        } else {
            $versions = collect();
        }

        foreach ($dirs as $dir) {
            $exe = $dir . DIRECTORY_SEPARATOR . 'php.exe';
            if(file_exists($exe)) {
                $version = ExeInfo::getFileVersion($exe);

                // $version = $exeInfo['FileVersion'];
                list($major, $minor, $patch) = explode('.', $version);

                $existing = $versions->where('path', $dir)->first();

                if(!$existing) {
                    $versions->push([
                        'path' => $dir,
                        'major_version' => $major,
                        'minor_version' => $minor,
                        'patch_version' => $patch,
                        'active' => false
                    ]);

                    $this->line('    - Discovered PHP ' . $version);
                    $discovered++;
                }
            }
        }

        Filesystem::put(storage_path() . '/versions.json', $versions->toJson());

        $this->line('');
        $this->info('✔ Discovered ' . $discovered . ' versions of PHP');
    }
}
