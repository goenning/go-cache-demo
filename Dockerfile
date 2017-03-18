FROM scratch
ADD main /
ENV PORT 8080
CMD ["/main", "-s", "memory"]
EXPOSE 8080