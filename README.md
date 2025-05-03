```bash
This program generates virtual links to the content of torrent files, enabling users to continue
sharing these files seamlessly. By creating these links, the program allows for the efficient 
distribution and access of shared files without the need for re-downloading or manually managing 
the original torrent files. This helps maintain the availability and integrity of shared content 
across multiple users and platforms.

Usage:
  torlinks [command]

Available Commands:
  check       Display the content of specified torrent file.
  completion  Generate the autocompletion script for the specified shell
  config      Creates a configuration file.
  dbClean     A brief description of your command
  dbSearch    A brief description of your command
  dbUpdate    Updates files database.
  dryRun      Similar to the Run command, but dry.
  help        Help about any command
  link        Search for torrent files and create symlinks to their data.
  retry       Allows to reprocess failed torrent files.
  stats       Provides some useful stats about your inbox folder.
  version     Displays build version.

Flags:
  -c, --config string   config file (default is $HOME/.config/torlinks/config.yaml)
  -h, --help            help for torlinks

Use "torlinks [command] --help" for more information about a command.
```
