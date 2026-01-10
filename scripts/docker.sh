
# The Docker container is built using the Dockerfile in the project directory.
docker build -t fortress-api:dev -f docker/dev.Dockerfile .

# Run the Docker container for Fortress API - Development
docker run -d \
  -p 8080:8080 \
  --name fortress-api-container \
  --env-file .env.development \
  -v "$(pwd):/app" \
  fortress-api:dev

# View the logs of the Docker container
docker logs -f fortress-api-container

# Build production image
docker build -t fortress-api:latest -f docker/prod.Dockerfile .

# Tag the docker image for Docker Hub
docker tag fortress-api:latest harmeet10000/fortress-api:latest

# Push the docker image to Docker Hub
docker push harmeet10000/fortress-api:latest

# Tag for AWS ECR and push (update ECR_REGISTRY with your AWS ECR URI)
docker tag fortress-api:latest $ECR_REGISTRY/fortress-api:latest
docker push $ECR_REGISTRY/fortress-api:latest

