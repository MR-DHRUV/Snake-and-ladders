FROM node:20-alpine AS builder

WORKDIR /app

# Enable corepack (for pnpm)
RUN corepack enable

# Copy only package files first (for better caching)
COPY pnpm-lock.yaml package.json ./

# Create .npmrc to ensure hoisted node_modules structure
RUN echo "node-linker=hoisted" > .npmrc

# Install all dependencies
RUN pnpm install --frozen-lockfile

# Copy source files
COPY . .

# Build the Vite app
RUN pnpm run build

# --------------------------
# Stage 2: Serve with NGINX
# --------------------------
FROM nginx:1.27-alpine

# Remove default nginx files
RUN rm -rf /usr/share/nginx/html/*

# Copy build output from builder
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy custom NGINX config (React Router support)
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 5173
EXPOSE 5173

CMD ["/bin/sh", "-c", "envsubst < /usr/share/nginx/html/env.js > /usr/share/nginx/html/env.tmp.js && mv /usr/share/nginx/html/env.tmp.js /usr/share/nginx/html/env.js && exec nginx -g 'daemon off;'"]