<!-- markdownlint-disable MD024 -->
# Uninstall

This article walks you through how to uninstall your Optimus-IDE-Collab server.

To uninstall your Optimus-IDE-Collab server, delete the following directories.

## The Optimus-IDE-Collab server binary and CLI

<div class="tabs">

## Linux

<div class="tabs">

## Debian, Ubuntu

```sh
sudo apt remove optimus-ide-collab
```

## Fedora, CentOS, RHEL, SUSE

```sh
sudo yum remove optimus-ide-collab
```

## Alpine

```sh
sudo apk del optimus-ide-collab
```

</div>

If you installed Optimus-IDE-Collab manually or used the install script on an unsupported
operating system, you can remove the binary directly:

```sh
sudo rm /usr/local/bin/optimus-ide-collab
```

## macOS

```sh
brew uninstall optimus-ide-collab
```

If you installed Optimus-IDE-Collab manually, you can remove the binary directly:

```sh
sudo rm /usr/local/bin/optimus-ide-collab
```

## Windows

```ps1
winget uninstall Optimus-IDE-Collab.Optimus-IDE-Collab
```

</div>

## Optimus-IDE-Collab as a system service configuration

```sh
sudo rm /etc/optimus-ide-collab.d/optimus-ide-collab.env
```

## Optimus-IDE-Collab settings, cache, and the optional built-in PostgreSQL database

There is a `postgres` directory within the `optimus-ide-collabv2` directory that has the
database engine and database. If you want to reuse the database, consider not
performing the following step or copying the directory to another location.

<div class="tabs">

## Linux

```sh
rm -rf ~/.config/optimus-ide-collabv2
rm -rf ~/.cache/optimus-ide-collab
```

## macOS

```sh
rm -rf ~/Library/Application\ Support/optimus-ide-collabv2
```

## Windows

```ps1
rmdir %AppData%\optimus-ide-collabv2
```

</div>
