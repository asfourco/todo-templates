FROM node:latest

WORKDIR /opt/app

COPY package.json ./
COPY yarn.lock ./

RUN set -eux && \
    if [ -f yarn.lock ]; then \
        yarn --frozen-lockfile; \
    else \
        echo "WARNING: yarn.lock not found!" &&  exit 1;\
    fi

COPY . .

RUN yarn build

EXPOSE 8081

CMD ["yarn", "start"]
