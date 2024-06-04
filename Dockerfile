FROM wetee/ego-ubuntu-deploy:22.04
WORKDIR /

ADD bin/*  /
ADD bin/keys /keys

RUN mkdir -p /opt/wetee-cache

EXPOSE 8881 

CMD ["/bin/sh", "-c" ,"ego sign indexer && ego run indexer"]