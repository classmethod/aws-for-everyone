FROM php:7.4.1-apache

RUN apt-get update \
    && apt-get install -y libpq-dev \
    && docker-php-ext-install pdo_mysql

COPY index.php /var/www/html
