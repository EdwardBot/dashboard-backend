FROM centurylink/ca-certs
ADD ./main /
ENV PORT=6000
EXPOSE ${PORT}
CMD ["/main"]
