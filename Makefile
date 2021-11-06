.PHONY : build
UNAME := $(shell uname -m)

all: | clean build

clean:
	rm -rf adax-prometheus-exporter
	rm -rf adax-prometheus-exporter_*
	rm -rf *.deb

build:
	go build -o adax-prometheus-exporter

debian:
	mkdir -p adax-prometheus-exporter_1.0-1_$(UNAME)/usr/local/bin
	mkdir -p adax-prometheus-exporter_1.0-1_$(UNAME)/etc/systemd/system
	cp adax-prometheus-exporter adax-prometheus-exporter_1.0-1_$(UNAME)/usr/local/bin
	cp adax-prometheus-exporter.service adax-prometheus-exporter_1.0-1_$(UNAME)/etc/systemd/system
	mkdir -p adax-prometheus-exporter_1.0-1_$(UNAME)/DEBIAN

	echo "Package: adax-prometheus-exporter" >> adax-prometheus-exporter_1.0-1_$(UNAME)/DEBIAN/control
	echo "Version: 1.0" >> adax-prometheus-exporter_1.0-1_$(UNAME)/DEBIAN/control
	echo "Architecture: $(UNAME)" >> adax-prometheus-exporter_1.0-1_$(UNAME)/DEBIAN/control
	echo "Maintainer: Vilhelm Prytz <vilhelm@prytznet.se>" >> adax-prometheus-exporter_1.0-1_$(UNAME)/DEBIAN/control
	echo "Description: Prometheus expoter for Adax heaters" >> adax-prometheus-exporter_1.0-1_$(UNAME)/DEBIAN/control

	dpkg-deb --build --root-owner-group adax-prometheus-exporter_1.0-1_$(UNAME)
