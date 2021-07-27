FROM python:3.7-slim-buster

WORKDIR /app

ADD app/ /app/
RUN pip3 install -r requirements.txt

CMD [ "python3", "-m" , "flask", "run", "--host=0.0.0.0", "--port=10001"]