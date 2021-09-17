FROM golang:1.16-buster as build

WORKDIR /app

COPY . ./

WORKDIR /app/cmd

RUN go build -o /etc/bin/dispatcher-jira

FROM gcr.io/distroless/base-debian10

ARG APP_PORT=""
ENV APP_PORT=${APP_PORT}

EXPOSE ${APP_PORT}

COPY --from=build /etc/bin/dispatcher-jira /app/dispatcher-jira

ENTRYPOINT ["/app/dispatcher-jira"]
