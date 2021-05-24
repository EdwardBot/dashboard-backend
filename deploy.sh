go build .
rm ./build/edward-backend
mv ./edward-backend ./build/edward-backend
sudo docker build -t taters.bendimester23.tk/dashboard-backend .
sudo docker push taters.bendimester23.tk/dashboard-backend
