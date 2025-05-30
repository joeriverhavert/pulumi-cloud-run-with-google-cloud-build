FROM nginx:stable

# Change workdir
WORKDIR /usr/share/nginx/html

# Copy Hyperspace-demo-app to workdir
COPY hyperspace-demo-app/ .

# Expose port tcp/80
EXPOSE 80
