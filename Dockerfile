FROM node:18-alpine

# Install dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy code
ADD package.json /app/package.json
ADD package-lock.json /app/package-lock.json

# Install deps
RUN npm install

# Copy other files
ADD . /app

# Expose the tracker port
EXPOSE 8080

# Start the tracker server
CMD ["npm", "start"]