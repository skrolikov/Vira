# frontend/vira-wish/Dockerfile
FROM node:20-alpine

WORKDIR /app

COPY . .

RUN npm install
RUN npm run build

EXPOSE 5174

CMD ["npm", "run", "dev", "--", "--host"]
