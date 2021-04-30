<?php

namespace App\Commands;

use Illuminate\Console\Scheduling\Schedule;
use Illuminate\Support\Facades\Storage;

class ClearCommand extends Command
{
    /**
     * The signature of the command.
     *
     * @var string
     */
    protected $signature = 'clear';

    /**
     * The description of the command.
     *
     * @var string
     */
    protected $description = 'Clears saved PHP installs';

    /**
     * Execute the console command.
     *
     * @return mixed
     */
    public function handle()
    {
        Storage::delete('versions.json');

        $this->info('🚮 Cleared all saved PHP versions.');
    }
}
