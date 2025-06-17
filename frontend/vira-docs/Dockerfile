# frontend/vira-docs/Dockerfile
FROM node:20-alpine

WORKDIR /app

COPY . .

RUN npm install
RUN npm run build

EXPOSE 5175

CMD ["npm", "run", "dev", "--", "--host"]
