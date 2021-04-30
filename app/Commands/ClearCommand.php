<?php

namespace App\Commands;

use App\Support\Filesystem;

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
        Filesystem::delete(storage_path('versions.json'));

        $this->info('🚮 Cleared all saved PHP versions.');
    }
}
