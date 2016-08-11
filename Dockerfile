FROM debian:jessie
COPY ./dist/maestrod/ ~/maestrod/
CMD ["~/maestrod/maestrod", "--config-path=/etc/maestrod/example.conf.toml"]
