<?php

namespace App\Commands;

use Illuminate\Console\Scheduling\Schedule;

class PathCommand extends Command
{
    /**
     * The signature of the command.
     *
     * @var string
     */
    protected $signature = 'path';

    /**
     * The description of the command.
     *
     * @var string
     */
    protected $description = 'Retrieve the path to the symlink';

    /**
     * Execute the console command.
     *
     * @return mixed
     */
    public function handle()
    {
        $this->line('🛠 Add this to your Path variable:');
        $this->line('');

        $this->info(base_path('bin'));
    }
}
