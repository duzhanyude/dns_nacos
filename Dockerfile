FROM centos:centos7
#MAINTAINER 维护者信息
MAINTAINER mono

VOLUME /tmp
COPY ./bin/main DNS
CMD ["chmod +x DNS"]
ENV N_IP  "172.16.1.81"
ENV  N_Name 8dfcdd4a-e7fb-4798-beb0-0c3cb84dd13e
EXPOSE 53

CMD ["./DNS","-N_IP","$N_IP","-N_Name","$N_Name"]