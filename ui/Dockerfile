FROM node:12.7-alpine AS build
WORKDIR /usr/src/app
COPY package.json package-lock.json ./
RUN npm install -g @angular/cli
RUN npm install
COPY . .

#RUN npm run build --configuration=production
RUN ng build --configuration=production

FROM nginx:alpine
COPY --from=build /usr/src/app/dist/ui /usr/share/nginx/html
COPY ./conf/nginx.conf /etc/nginx/nginx.conf
ENTRYPOINT ["nginx", "-g", "daemon off;"]

