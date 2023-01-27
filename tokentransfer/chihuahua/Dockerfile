# This Dockerfile will: 
# set up a Golang environment(using Golang Docker Image)
# install Ignite
# clone the chihuahua repository
# start a chihuahua blockchain
# expose api endpoints

# Use Golang Docker Image
FROM golang:1.18

# Set working dir
WORKDIR /home/chihuahua

# Install tools
RUN apt-get install curl git

# Install ignite
RUN curl -L https://get.ignite.com/cli@v0.22.1! | bash

# Clone chihuahua game repository(branch wager-denomination)
RUN git clone -b wager-denomination https://github.com/ChihuahuaChain/chihuahua.git /home/chihuahua

# Copy Config.yml
ARG configyml
COPY ${configyml} /home/chihuahua//config.yml

# copy entrypoint script
COPY ./run-chihuahua.sh /home/chihuahua
RUN chmod +x ./run-chihuahua.sh

ENTRYPOINT ["/bin/bash", "-c", "/home/chihuahua/run-chihuahua.sh"]
EXPOSE 26657 1317 4500

# build with ./build-images.sh
