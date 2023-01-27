# This Dockerfile will: 
# set up a Golang environment(using the Golang Docker Image)
# and install the Go Relayer

# Use Golang Docker Image
FROM golang:1.18

# Set working dir
WORKDIR /home/rly

# Install git
RUN apt-get install git

# Clone relayer repository and install relayer
RUN git clone https://github.com/cosmos/relayer.git && cd relayer && git checkout v2.0.0-rc3 && make install

# Copy configs
ADD configs /home/rly/configs

# Copy pahts
ADD paths /home/rly/paths

# Copy entrypoint script
COPY ./run-relayer.sh /home/rly/

# Set permission
RUN chmod +x ./run-relayer.sh

# Run entrypoint script
# ENTRYPOINT ["/bin/bash", "-c", "/home/rly/run-relayer.sh"]

# build with
# docker build -f Dockerfile . -t relayer_go

