FROM ubuntu:latest

RUN apt-get update && apt-get install -y curl git

# install nvm
RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.3/install.sh | bash
ENV NVM_DIR=/root/.nvm
RUN bash -c "source $NVM_DIR/nvm.sh && nvm install --lts"

WORKDIR /home/app

COPY main.sh main.sh
COPY entrypoint.sh entrypoint.sh
COPY script.js script.js
COPY s3.js s3.js
COPY redis.js redis.js
COPY package*.json .

RUN npm install

RUN chmod +x main.sh entrypoint.sh
RUN chmod +x script.js
RUN chmod +x s3.js
RUN chmod +x redis.js

ENTRYPOINT ["/home/app/entrypoint.sh"]
CMD ["/bin/bash"]
