FROM node:18-alpine

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . . 
RUN npm run build

# Устанавливаем статический сервер
RUN npm install -g serve

# Команда для запуска
CMD ["serve", "-s", "build", "-l", "3000"]

EXPOSE 3000
