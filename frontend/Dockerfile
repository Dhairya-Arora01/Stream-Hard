FROM node:20-alpine as build

WORKDIR /StreamHard/frontend

COPY package.json .

RUN npm install

COPY . .

RUN npm run build

FROM nginx:stable-alpine

COPY --from=build /StreamHard/frontend/dist  /usr/share/nginx/html

EXPOSE 80

CMD [ "nginx", "-g", "daemon off;" ]