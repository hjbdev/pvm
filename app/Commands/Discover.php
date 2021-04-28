<?php

namespace App\Commands;

use App\Support\ExeInfo;
use Illuminate\Console\Scheduling\Schedule;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\File;
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

        foreach ($dirs as $dir) {
            $exe = $dir . DIRECTORY_SEPARATOR . 'php.exe';
            if(File::exists($exe)) {
                $exeInfo = ExeInfo::get($exe);
                $version = $exeInfo['FileVersion'];
                list($major, $minor, $patch) = explode('.', $exeInfo['FileVersion']);

                $existing = DB::table('versions')->where('path', $dir)->first();

                if(!$existing) {
                    DB::table('versions')->insert([
                        'path' => $dir,
                        'major_version' => $major,
                        'minor_version' => $minor,
                        'patch_version' => $patch
                    ]);

                    $this->info('    - Discovered PHP ' . $version);
                    $discovered++;
                }
            }
        }

        $this->line('');
        $this->info('✅ Discovered ' . $discovered . ' versions of PHP');
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
