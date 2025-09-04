FROM python:alpine

WORKDIR /

RUN  mkdir -p /db/migrations

ADD https://github.com/amacneil/dbmate/releases/latest/download/dbmate-linux-amd64 /usr/local/bin/dbmate

ADD ./db/migrations db/migrations
RUN chmod +x /usr/local/bin/dbmate

ADD ./dbmate-entrypoint.sh dbmate-entrypoint.sh
RUN chmod +x /dbmate-entrypoint.sh

RUN pip3 install --no-cache-dir s4cmd

ENTRYPOINT [ "/dbmate-entrypoint.sh" ]

CMD [ "/usr/local/bin/dbmate", "up" ]
