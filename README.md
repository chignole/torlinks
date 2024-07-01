# Torlinks

### Keep seeding through symlinks.

This program generates virtual links to the content of torrent files, enabling users to continue
sharing these files seamlessly. By creating these links, the program allows for the efficient distribution
and access of shared files without the need for re-downloading or manually managing the original torrent 
files. This helps maintain the availability and integrity of shared content across multiple users and platforms.

```bash
Usage:
  torlinks [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  inbox       Provides some useful stats about your inbox folder.
  link        Search for torrent files and create symlinks to their data.
  retry       Allows to reprocess failed torrent files.
  updateDb    Updates files database.

Flags:
  -c, --config string   config file (default is $HOME/.config/torlinks/config.yaml
  -h, --help            help for torlinks

Use "torlinks [command] --help" for more information about a command.
```
