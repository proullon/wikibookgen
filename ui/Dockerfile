FROM nginx:alpine

COPY ./dist/ui /usr/share/nginx/html
COPY ./conf/nginx.conf /etc/nginx/nginx.conf

ENTRYPOINT ["nginx", "-g", "daemon off;"]

