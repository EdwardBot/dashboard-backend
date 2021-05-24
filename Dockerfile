FROM ubuntu
WORKDIR /app
COPY ./build /app
ENV PORT=6000
EXPOSE ${PORT}
RUN chmod +x /app/start.sh
CMD ["/bin/sh"]
