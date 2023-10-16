FROM node:latest as base

RUN set -eux && \
    apt update && \
    apt install -y --no-install-recommends \
        bash \
        software-properties-common \
        ca-certificates \
        git \
        build-essential && \
        apt autoremove -y && \
        apt clean &&

FROM base as dependencies

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
