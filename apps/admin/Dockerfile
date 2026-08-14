FROM node:22-slim AS build

WORKDIR /repo

COPY package.json package-lock.json ./
COPY apps/admin/package.json apps/admin/package.json
COPY packages/shared/package.json packages/shared/package.json
RUN npm ci --workspace=apps/admin --include-workspace-root

COPY packages/shared packages/shared
COPY apps/admin apps/admin
ARG VITE_API_URL=http://localhost:5001
ENV VITE_API_URL=$VITE_API_URL
RUN npm run build --workspace=apps/admin

FROM nginx:1.27-alpine

COPY --from=build /repo/apps/admin/dist /usr/share/nginx/html
COPY apps/admin/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
