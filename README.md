<div align="center">
    <img width="300" src="./assets/images/portal.png" alt="Portal Logo">  
    <H1>Portal</H1>
    <p>Portal is an end-to-end encrypted chat application, ensuring your conversations remain private and secure. 🗝️🔒<p>
</div>

<hr>

## Getting Started 🚀

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile 🏗️

run all make commands with clean tests
```bash
make all build
```

build the application
```bash
make build
```

run the application
```bash
make run
```

Create Docker container
```bash
# Create Docker SQLite container
make docker-run
# Shutdown Docker SQLite container
make docker-down

# Create Docker Postgres container
make docker-run-postgres
# Shutdown Docker Postgres container
make docker-down-postgres
```

Create Podman container
```bash
# Create Podman SQLite container
make podman-run
# Shutdown Podman SQLite container
make podman-down

# Create Podman Postgres container
make podman-run-postgres
# Shutdown Podman Postgres container
make podman-down-postgres
```

live reload the application
```bash
make watch
```

run the test suite
```bash
make test
```

clean up binary from the last build
```bash
make clean
```

## Credits 🙏

- The initial project structure was created using the [Go Blueprint](https://go-blueprint.dev/) project. 🏗️📝


## Contribute 🤝

We welcome contributions from the community. Feel free to open an issue or submit a pull request! 💡🔧

## License 📄

Portal is licensed under the [BSD-4-Clause](https://en.wikipedia.org/wiki/BSD_licenses) License. See [`LICENSE`](./LICENSE) for more information. 📜