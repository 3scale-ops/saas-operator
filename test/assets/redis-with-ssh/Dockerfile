# Use the redis:4.0.11-alpine as the base image
FROM redis:6.2.13-alpine

# update the image and install necessary tools
RUN apk update && apk upgrade
RUN apk add --no-cache openssh gzip python3 py3-pip sudo
RUN pip3 install --upgrade pip
RUN pip3 install awscli

# configure SSH to run on port 2222
RUN sed -i 's/#Port 22/Port 2222/' /etc/ssh/sshd_config

# setup "docker" user with sudo
RUN adduser --disabled-password docker && passwd docker -d ''
RUN echo "docker  ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers
RUN ssh-keygen -A
RUN mkdir -p /home/docker/.ssh
COPY authorized_keys /home/docker/.ssh/
RUN chmod 700 /home/docker/.ssh && chmod 600 /home/docker/.ssh/authorized_keys
RUN chown -R docker:docker /home/docker/.ssh

# entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

USER root

# Start both Redis and SSH using the entrypoint script
CMD ["/entrypoint.sh"]
