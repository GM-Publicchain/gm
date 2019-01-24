FROM ubuntu:16.04

RUN apt-get update \
  && apt-get install -y libssl-dev

WORKDIR /root
COPY chain33 chain33
COPY chain33-cli chain33-cli
COPY chain33.toml ./

CMD ["/root/chain33", "-f", "/root/chain33.toml"]
