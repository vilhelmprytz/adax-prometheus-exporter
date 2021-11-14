# adax-prometheus-exporter

Prometheus exporter for [Adax.no](https://adax.no/) heaters. Retrieves current and target temperature for each home and room registered and returns in a format that Prometheus can read. Written in Go.

## Installation

If you are on Debian or Ubuntu (or any of its derivatives), the easiest way to get the exporter running is using my apt repo.

Firstly, install the required dependencies.

```bash
sudo apt-get update
sudo apt-get install \
    ca-certificates \
    curl \
    apt-transport-https \
    gnupg \
```

Then, download the key used to sign the packages and add the repository to your `/etc/sources.list.d/` directory.

```bash
curl -o /etc/apt/trusted.gpg.d/packages.hejduk.se.gpg https://packages.hejduk.se/apt.gpg
sh -c 'echo "deb https://packages.hejduk.se/hejduk stable main" > /etc/apt/sources.list.d/hejduk.list'
```

Finally, use `apt` to update the cache and install the package.

```bash
apt update
apt install adax-prometheus-exporter
```

## Configuration

Once installed, create/open the file `/etc/adax-prometheus-exporter/config.yml` (if you used the `.deb` packages, this file should have automatically been created). You may use the `config.example.yml` as a reference if you are installing manually.

It will look something like this. Fill in the values for the options as required.

```yaml
---
client_id: "0"
client_secret: ""
port: 8080
```

If you installed the `.deb` package, you may start the web app using `systemctl`.

```bash
systemctl enable adax-prometheus-exporter
systemctl start adax-prometheus-exporter
```

If you have built it manually, you can launch the binary with the following parameters.

```bash
adax-prometheus-exporter --config /etc/adax-prometheus-exporter/config.yml
```

Metrics will be available from `http://localhost:8080/metrics`. It will look something like this.

```
room_temperature{home="Home",room="Test"} 21.000000
room_target_temperature{home="Home",room="Test"} 24.000000
room_temperature{home="Home",room="Test2"} 18.000000
room_target_temperature{home="Home",room="Test2"} 17.000000
room_temperature{home="Cottage",room="Living Room"} 22.000000
room_target_temperature{home="Cottage",room="Living Room"} 23.000000
```

## Building from source

To build the project from the source, make sure you have [Go](https://golang.org/) installed (Go 1.17.x is supported by other versions may work).

To just build the executable, use `make build`.

```bash
make build
```

To build the executable and the `.deb` packages (includes systemd unit file), just run `make`.

```bash
make
```

## Author ✨

Created and written by [Vilhelm Prytz](https://github.com/vilhelmprytz) - [vilhelm@prytznet.se](mailto:vilhelm@prytznet.se).
