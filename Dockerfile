FROM centos:7.4.1708
LABEL maintainers="Secyu Cloud Authors"
LABEL description="An Tcp Server for pingpong test"

COPY pp /bin/pp
RUN chmod +x /bin/pp

ENTRYPOINT ["/bin/pp","-s","-d"]
