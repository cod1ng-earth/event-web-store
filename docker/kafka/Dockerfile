ARG KAFKA_IMAGE=bitnami/kafka:2.3.0-r35

FROM ${KAFKA_IMAGE}
VOLUME [ "/opt/bitnami/kafka", "/tmp" ]
ENV KAFKA_CFG_ZOOKEEPER_CONNECT=localhost:2181
USER 65535
