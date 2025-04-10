cd ..
docker-compose -f docker-compose.deploy.yml -f docker-compose.test.integration.yml up --build --abort-on-container-exit --exit-code-from backend

docker-compose -f docker-compose.deploy.yml -f docker-compose.test.integration.yml down -v

pause
