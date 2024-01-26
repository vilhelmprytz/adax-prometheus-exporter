#!/bin/sh

envsubst </app/config.env.yml >/app/config.yml

/app/adax-prometheus-exporter --config /app/config.yml
