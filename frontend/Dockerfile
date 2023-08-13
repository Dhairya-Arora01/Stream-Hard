# Build it from outside frontend folder i.e. inside StreamHard.
FROM node:20-alpine

WORKDIR /StreamHard/frontend

COPY . .

RUN npm install

EXPOSE 5173

CMD [ "npm", "run", "dev", "--", "--host" ]