2 docker commands to build and run the program to listen to port 8080

docker build -t R-Pawel/fetch-takehome:latest .   
docker run -p 8080:8080 R-Pawel/fetch-takehome:latest