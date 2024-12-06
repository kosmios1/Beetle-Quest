FROM python:3.13-alpine3.20
# FROM python:3.13.0a1-slim-bullseye

WORKDIR /app

# NOTE: Python libraries that require compilation
# needs these packages to be installed
RUN apk add gcc python3-dev musl-dev linux-headers

RUN pip3 install --upgrade pip

COPY ./requirements.txt requirements.txt
RUN pip3 install -r requirements.txt

EXPOSE 8089

COPY ./locustfile.py locustfile.py
CMD ["locust", "-f", "locustfile.py", "--headless", "--run-time", "450"]
