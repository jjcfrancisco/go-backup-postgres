FROM golang:latest

## Installs postgresql-15
RUN apt-get update && apt upgrade -y
RUN apt-get install -y lsb-release gnupg2 wget vim && apt-get clean all
RUN sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN apt-get update && apt upgrade -y
RUN apt install locales postgresql-15 -y

## Generate and set the locale environment variables
RUN echo "en_GB.UTF-8 UTF-8" > /etc/locale.gen && \
    locale-gen en_GB.UTF-8 && \
    /usr/sbin/update-locale LANG=en_GB.UTF-8
ENV LANG=en_GB.UTF-8 \
    LANGUAGE=en_GB:en \
    LC_ALL=en_GB.UTF-8

RUN mkdir -p /go-backup-postgres/
WORKDIR /go-backup-postgres
COPY app /go-backup-postgres
RUN go build -o app
ENTRYPOINT ["/go-backup-postgres/app"]