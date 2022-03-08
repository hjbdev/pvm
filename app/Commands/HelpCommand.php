<?php

namespace App\Commands;

use Codedungeon\PHPCliColors\Color;

class HelpCommand extends Command
{
    /**
     * The signature of the command.
     *
     * @var string
     */
    protected $signature = 'help';

    /**
     * The description of the command.
     *
     * @var string
     */
    protected $description = 'Help command';

    /**
     * Execute the console command.
     *
     * @return mixed
     */
    public function handle()
    {
        $this->line('');
        $this->info('📜 pvm - PHP Version Manager');
        $this->line('');
 
        $this->info('Available commands:');
        $this->line('');

        foreach(pvm_app()->commands() as $command) {
            $this->line(Color::GREEN . $command->signature() . Color::RESET . ' - ' . $command->description());
            $this->line('');
        }
    }
}
