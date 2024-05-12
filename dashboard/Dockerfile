# Use a lightweight Node.js image as base
FROM node:alpine AS builder

# Set working directory in the container
WORKDIR /app

# Copy dashboard/package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application
COPY . .

# Build the application
RUN npm run build

# Use a smaller image for serving the application
FROM node:alpine

# Set working directory in the container
WORKDIR /app

# Copy built application from builder stage
COPY --from=builder /app .

# Expose port
EXPOSE 3000

# Command to run the application
CMD ["npm", "run", "server"]
