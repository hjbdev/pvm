<?php

namespace App\Commands;

use App\Support\ExeInfo;
use Illuminate\Console\Scheduling\Schedule;
use Illuminate\Support\Facades\File;
use Illuminate\Support\Facades\Storage;
use LaravelZero\Framework\Commands\Command;

class Discover extends Command
{
    /**
     * The signature of the command.
     *
     * @var string
     */
    protected $signature = 'discover {path?}';

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
        
        $path = $this->argument('path') ?? 'C:\laragon\bin\php';
        $dirs = File::directories($path);
        $discovered = 0;

        if (Storage::has('versions.json')) {
            $versions = collect(json_decode(file_get_contents(storage_path('versions.json'))));
        } else {
            $versions = collect();
        }

        foreach ($dirs as $dir) {
            $exe = $dir . DIRECTORY_SEPARATOR . 'php.exe';
            if(File::exists($exe)) {
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

        Storage::put('versions.json', $versions->toJson());

        $this->line('');
        $this->info('✔ Discovered ' . $discovered . ' versions of PHP');
    }

    /**
     * Define the command's schedule.
     *
     * @param  \Illuminate\Console\Scheduling\Schedule $schedule
     * @return void
     */
    public function schedule(Schedule $schedule): void
    {
        // $schedule->command(static::class)->everyMinute();
    }
}
