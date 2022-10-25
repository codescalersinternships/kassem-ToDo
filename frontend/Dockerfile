FROM node:14-alpine

WORKDIR /usr/src/app
ENV CYPRESS_INSTALL_BINARY=0
COPY rollup.config.js ./
COPY package*.json ./

RUN npm install

COPY ./src ./src
COPY ./public ./public
COPY build-env.sh ./
EXPOSE 8080
ENV HOST=0.0.0.0
CMD [ "sh", "build-env.sh" ]