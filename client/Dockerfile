FROM node:9.3.0

RUN mkdir -p /usr/client
WORKDIR /usr/client

COPY package.json yarn.lock /usr/client/
RUN yarn install

COPY . .

RUN yarn build
RUN yarn global add serve

CMD ["serve", "-s", "build"]
