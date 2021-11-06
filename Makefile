.PHONY : build
UNAME := $(shell uname -m)
VERSION=`cat version`

ARCH=$(UNAME)
ifeq ($(UNAME),x86_64)
	ARCH := amd64
endif

all: | clean build

clean:
	rm -rf adax-prometheus-exporter
	rm -rf adax-prometheus-exporter_*
	rm -rf *.deb

build:
	go build -o adax-prometheus-exporter

debian:
	mkdir -p adax-prometheus-exporter_$(VERSION)-1_$(ARCH)/usr/local/bin
	mkdir -p adax-prometheus-exporter_$(VERSION)-1_$(ARCH)/etc/systemd/system
	mkdir -p adax-prometheus-exporter_$(VERSION)-1_$(ARCH)/etc/adax-prometheus-exporter

	cp adax-prometheus-exporter adax-prometheus-exporter_$(VERSION)-1_$(ARCH)/usr/local/bin
	cp adax-prometheus-exporter.service adax-prometheus-exporter_$(VERSION)-1_$(ARCH)/etc/systemd/system
	cp config.example.yml adax-prometheus-exporter_$(VERSION)-1_$(ARCH)/etc/adax-prometheus-exporter/config.yml

	mkdir -p adax-prometheus-exporter_$(VERSION)-1_$(ARCH)/DEBIAN
	rm -rf adax-prometheus-exporter_$(VERSION)-1_$(ARCH)/DEBIAN/control
	echo "Package: adax-prometheus-exporter" >> adax-prometheus-exporter_$(VERSION)-1_$(ARCH)/DEBIAN/control
	echo "Version: $(VERSION)" >> adax-prometheus-exporter_$(VERSION)-1_$(ARCH)/DEBIAN/control
	echo "Architecture: $(ARCH)" >> adax-prometheus-exporter_$(VERSION)-1_$(ARCH)/DEBIAN/control
	echo "Section: Other" >> adax-prometheus-exporter_$(VERSION)-1_$(ARCH)/DEBIAN/control
	echo "Priority: 500" >> adax-prometheus-exporter_$(VERSION)-1_$(ARCH)/DEBIAN/control
	echo "Maintainer: Vilhelm Prytz <vilhelm@prytznet.se>" >> adax-prometheus-exporter_$(VERSION)-1_$(ARCH)/DEBIAN/control
	echo "Description: Prometheus expoter for Adax heaters" >> adax-prometheus-exporter_$(VERSION)-1_$(ARCH)/DEBIAN/control

	dpkg-deb --build --root-owner-group adax-prometheus-exporter_$(VERSION)-1_$(ARCH)
