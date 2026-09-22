# Deployment Strategy for the Xmatch Service
Binaries are built on the CI pipeline (GitHub Actions) and downloaded to the production server. The production server does not need Nix, devenv, or any build toolchain.

## Release artifact
The release binary is a single, fully static amd64 executable. It is built from the `outputs.xmatch` devenv output with:

```sh
devenv shell xwave-release   # writes service/build/main
```

It statically links the HEALPix libraries (healpix_cxx, libsharp) and libc (musl), so it runs on a plain Ubuntu 24.04 server without Nix, system libraries, or `sudo`. The GitHub release workflow uploads it as the `main` asset. The deploy program downloads it to a temporary file, verifies its sha256 digest when the release API provides one, and only then promotes it.

> [!NOTE] Downsides of the current setup
> Since the server requires a VPN connection to access via SSH we can't directly copy the built binaries using something like `scp`.

The Xmatch service repository has `systemd` files to be copied to the host machine in the required location `~/.config/systemd/user`.
The unit execs the promoted binary directly (`ExecStart=%h/deployment/production/bin/prod server`) — no shell or Nix profile is involved — and sets `WorkingDirectory` and `EnvironmentFile` to `~/deployment/production`.

Other versions of the application binaries are located at `~/deployment/binaries` and there's a script that handles promotion and rollback of a binary. In pseudocode it does this:

```python
function deploy(instances):
    release = get_latest_release(RELEASE_URL)
    binary_path = download("deployment/binaries/{release.tag}", asset="main", digest=release.digest)
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
│   ├── v1.0.0
│   └── v1.0.1
└── production
    ├── bin
    │   └── prod -> ../../binaries/v1.0.1
    ├── configs
    │   └── config.yaml
    ├── db
    │   └── production.db
    └── envfile
```
