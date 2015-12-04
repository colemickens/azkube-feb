#FROM scratch
#ADD azkube-deploy /opt/azkube-deploy/azkube-deploy
#ADD templates     /opt/azkube-deploy/templates
#WORKDIR           /opt/azkube-deploy
#CMD ["/opt/azkube-deploy/azkube-deploy"]

# we'll switch back when cfssl is put in

FROM ubuntu:15.10
RUN bash -c "apt-get update; apt-get install openssl; apt-get autoclean"
ADD azkube /opt/azkube/azkube
ADD templates /opt/azkube/templates
WORKDIR /opt/azkube
