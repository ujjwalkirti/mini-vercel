FROM node:20-bullseye

WORKDIR /home/app

RUN apt-get update && apt-get install -y git

COPY package*.json .
RUN npm install

COPY . .

RUN chmod +x main.sh script.js kafkaProducer.js azureBlob.js

ENTRYPOINT ["/home/app/main.sh"]
CMD ["/bin/bash"]
