FROM centos:centos7
#MAINTAINER 维护者信息
MAINTAINER mono

VOLUME /tmp
COPY ./bin/main DNS
CMD ["chmod +x DNS"]
LABEL N_IP="172.16.1.81"
LABEL  N_Name=8dfcdd4a-e7fb-4798-beb0-0c3cb84dd13e
EXPOSE 53
EXPOSE 10053

CMD ["./DNS","-nacos_ip","$nacos_ip","-nacos_name","$nacos_name"]