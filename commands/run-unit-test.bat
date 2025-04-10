cd ..
echo "Run unit tests"

docker-compose -f docker-compose.test.unit.yml up
docker-compose -f docker-compose.test.unit.yml down

pause
