# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./
COPY prisma ./prisma/
COPY prisma.config.ts ./

# Install all dependencies (including dev for building)
RUN npm ci

# Copy source code
COPY tsconfig.json ./
COPY src ./src/

# Generate Prisma client
RUN npx prisma generate

# Build TypeScript
RUN npm run build

# Production stage
FROM node:20-alpine AS production

WORKDIR /app

# Create non-root user for security
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nodejs -u 1001

# Copy package files
COPY package*.json ./
COPY prisma ./prisma/
COPY prisma.config.ts ./

# Install production dependencies only
RUN npm ci --only=production

# Generate Prisma client in production
RUN npx prisma generate

# Copy built files from builder
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/generated ./generated

# Copy Kafka CA certificate if exists
COPY --chown=nodejs:nodejs src/config/ca.pem ./dist/config/ca.pem

# Set ownership
RUN chown -R nodejs:nodejs /app

# Switch to non-root user
USER nodejs

# Expose port
EXPOSE 9000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD node -e "fetch('http://localhost:9000/health').then(r => r.ok ? process.exit(0) : process.exit(1)).catch(() => process.exit(1))"

# Start the application
CMD ["node", "dist/server.js"]
