# Build the Svelte SPA
FROM debian:buster

VOLUME /app
WORKDIR /app

RUN apt update -y \
 && apt install -y \
        git \
        npm \
        nodejs \
 && rm -rf /var/lib/apt/lists/*
