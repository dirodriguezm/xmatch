# Deployment Strategy for the Xmatch Service
Binaries should be built either on the production machine or on the CI pipeline (Github Actions) and then downloaded to the production server.

> [!NOTE] Downsides of the current setup
> Since the server requires a VPN connection to access via SSH we can't directly copy the built binaries using something like `scp`.

The Xmatch service repository has `systemd` files to be copied to the host machine in the required location `~/.config/systemd/user`.
The systemd files will point to the current production binary located at `~/deployment/production/bin`.

Other versions of the application binaries are located at `~/deployment/binaries` and there's a script that handles promotion and rollback of a binary. In pseudocode it does this:

```python
function deploy(commit_hash, instances):
    binary_path = find_file("deployment/binaries", pattern="*{commit_hash}*")
    previous_binary = resolve_symlink("deployment/production/bin/prod")

    # Promote new binary
    update_symlink("deployment/production/bin/prod", binary_path)

    success = true
    for instance in instances:
        result = restart_service("myservice@{instance}")
        if not result:
            print("Error: Failed to restart instance myservice@{instance}")
            success = false
            break

    if success:
        print("Deployment successful")
        exit(0)

    # Rollback
    print("Rolling back...")
    update_symlink("deployment/production/bin/prod", previous_binary)

    for instance in instances:
        result = restart_service("myservice@{instance}")
        if not result:
            print("Critical: Rollback failed on instance myservice@{instance}")
            exit(2)

    print("Rollback successful")
    exit(1)
```

Using this strategy we can perform rolling releases for new binary files, ensuring availability of the service and rolling back potentially broken binaries.

## Database
The directory `~/deployment/db` contains database files required for the app to function. While the connection to the database is made through a configuration file, and the database file itself could be anywhere, this directory should be used to identify production databases. 

The configuration file should always point to this directory

## Changing App Configuration
The directory `~/deployment/configs` contains configuration files (yaml files). The service should use the `CONFIG_PATH` environment variable pointing to a config file defined in this directory. The systemd file references an env file `~/deployment/production/envfile` that specifiers the needed environment variables.

## Example directory tree
```
/home/user/deployment
├── binaries
│   ├── main_SHA1
│   └── main_SHA2
├── production
│   ├── bin
│   │   └── prod -> ../../binaries/main_SHA2
│   ├── configs
│   │   └── config.yaml
│   ├── db
│   │   └── production.db
│   └── envfile
├── flake.nix
└── release_script
```
