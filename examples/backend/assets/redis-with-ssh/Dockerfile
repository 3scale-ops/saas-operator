# Use the redis:4.0.11-alpine as the base image
FROM redis:4.0.11-alpine

# update the image and install necessary tools
RUN apk update && apk upgrade
RUN apk add --no-cache openssh gzip python3 py3-pip
RUN pip3 install --upgrade pip
RUN pip3 install awscli


# configure SSH to run on port 2222
RUN sed -i 's/#Port 22/Port 2222/' /etc/ssh/sshd_config
RUN ssh-keygen -A
RUN mkdir -p /root/.ssh
COPY authorized_keys /root/.ssh/
RUN chmod 700 /root/.ssh && chmod 600 /root/.ssh/authorized_keys

# entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

USER root

# Start both Redis and SSH using the entrypoint script
CMD ["/entrypoint.sh"]


