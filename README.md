docker run --rm -it $(docker build -q .)


# イメージを作る
- docker build .

# イメージからコンテナを立ち上げる(deamon)
- docker run -d %{imageID}

# コンテナの起動
- docker start ${container_name}

# docker-compose

- docker-compose build