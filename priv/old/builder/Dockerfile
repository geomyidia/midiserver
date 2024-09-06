FROM golang:1.17

LABEL "name"="Go midiserver Builder"
LABEL "version"="0.1.0"

LABEL "com.github.actions.name"="Go midiserver Builder"
LABEL "com.github.actions.description"="Cross-complile midiserver"
LABEL "com.github.actions.icon"="package"
LABEL "com.github.actions.color"="#E0EBF5"

RUN \
  apt-get update && \
  apt-get install -y ca-certificates openssl zip libasound2-dev libasound2 \
  alsa-utils libasound2-plugins pulseaudio pulseaudio-utils \
  && \
  update-ca-certificates && \
  rm -rf /var/lib/apt

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
