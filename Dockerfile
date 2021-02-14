# TODO: write the build for minimal Go environment and headless build
FROM golang:1.14 as builder
WORKDIR /pod
COPY .git /pod/.git
COPY app /pod/app
COPY cmd /pod/cmd
COPY pkg /pod/pkg
COPY stroy /pod/stroy
COPY version /pod/version
COPY go.??? /pod/
COPY ./*.go /pod/
RUN ls
ENV GOBIN "/bin"
ENV PATH "$GOBIN:$PATH"
RUN cd /pod && go install ./stroy/.
RUN cd /pod && stroy stroy
RUN cd /pod && stroy docker
RUN stroy teststopnode
FROM ubuntu:20.04
ENV GOBIN "/bin"
ENV PATH "$GOBIN:$PATH"
#RUN /usr/bin/which sh
COPY --from=builder /bin/stroy /bin/
COPY --from=builder /bin/pod /bin/
#COPY --from=builder /usr/bin/sh /bin/
#RUN echo $PATH && /bin/stroy
RUN /bin/pod version
EXPOSE 11048 11047 21048 21047
CMD ["tail", "-f", "/dev/null"]
#CMD /usr/local/bin/parallelcoind -txindex -debug -debugnet -rpcuser=user -rpcpassword=pa55word -connect=127.0.0.1:11047 -connect=seed1.parallelcoin.info -bind=127.0.0.1 -port=11147 -rpcport=11148
