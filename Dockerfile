FROM alpine
WORKDIR /app
COPY ./build /app
ENV PORT=6000
EXPOSE ${PORT}
RUN chmod +x /app/start.sh
RUN ls -l /app
CMD ["/bin/sh"]
