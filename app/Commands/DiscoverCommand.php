<?php

namespace App\Commands;

use App\Support\ExeInfo;
use App\Support\Filesystem;
use Illuminate\Support\Facades\File;
use Illuminate\Support\Facades\Storage;

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
        
        $path = app()->argument('path') ?? 'C:\laragon\bin\php';
        $dirs = Filesystem::directories($path);
        $discovered = 0;

        if (Filesystem::has(storage_path('versions.json'))) {
            $versions = collect(json_decode(Filesystem::get(storage_path('versions.json'))));
        } else {
            $versions = collect();
        }

        foreach ($dirs as $dir) {
            $exe = $dir . DIRECTORY_SEPARATOR . 'php.exe';
            if(file_exists($exe)) {
                $exeInfo = ExeInfo::get($exe);
                $version = $exeInfo['FileVersion'];
                list($major, $minor, $patch) = explode('.', $exeInfo['FileVersion']);

                $existing = $versions->where('path', $dir)->first();

                if(!$existing) {
                    $versions->push([
                        'path' => $dir,
                        'major_version' => $major,
                        'minor_version' => $minor,
                        'patch_version' => $patch,
                        'active' => 0
                    ]);

                    $this->info('    - Discovered PHP ' . $version);
                    $discovered++;
                }
            }
        }

        Filesystem::put('storage/versions.json', $versions->toJson());

        $this->line('');
        $this->info('✔ Discovered ' . $discovered . ' versions of PHP');
    }
}
