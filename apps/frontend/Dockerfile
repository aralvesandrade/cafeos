FROM node:22-slim AS build

WORKDIR /repo

COPY package.json package-lock.json ./
COPY apps/frontend/package.json apps/frontend/package.json
COPY packages/shared/package.json packages/shared/package.json
RUN npm ci --workspace=apps/frontend --include-workspace-root

COPY packages/shared packages/shared
COPY apps/frontend apps/frontend
RUN npm run build --workspace=apps/frontend

FROM nginx:1.27-alpine

COPY --from=build /repo/apps/frontend/dist /usr/share/nginx/html
COPY apps/frontend/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
