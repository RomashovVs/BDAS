FROM nginx:alpine

# Удаляем дефолтный конфиг
RUN rm /etc/nginx/conf.d/default.conf

# Копируем наш конфиг
COPY nginx.conf /etc/nginx/conf.d/

# Создаем директорию для SSL (файлы добавятся через volume)
RUN mkdir -p /etc/nginx/ssl

EXPOSE 8081 443
CMD ["nginx", "-g", "daemon off;"]