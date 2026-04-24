#!/bin/bash

no_color=$(echo -e "\033[0m")
yellow=$(echo -e "\033[1;33m")
green=$(echo -e "\033[1;32m")
red=$(echo -e "\033[1;31m")

echo -e "${yellow}Setup script for initial setup of QuRe Reservation System${no_color}\nPress Ctrl+C to cancel"
sleep 1

echo -e "\n${yellow}Copying Files...${no_color}"
echo "docker-compose.yml"
curl -O "https://raw.githubusercontent.com/JValtteri/qure/refs/heads/main/docker-compose.yml"
echo "config.json"
curl -o ./config.json "https://raw.githubusercontent.com/JValtteri/qure/refs/heads/main/server/config.json.example"
echo "update.sh"
curl -O "https://raw.githubusercontent.com/JValtteri/qure/refs/heads/main/update.sh"
chmod +x update.sh
echo "logo.png"
curl -O "https://raw.githubusercontent.com/JValtteri/qure/refs/heads/main/client/public/logo.png"

echo -e "\n${yellow}Creating Dedicated User${no_color}"
echo "name=qure UID=10001"
sudo useradd -u 10001 qure

echo -e "\n${yellow}Preparing mount folders${no_color}"
echo "./db/"
mkdir db
sudo chown -R qure db/
echo "./images/"
mkdir images
sudo chown -R qure images/

echo -e "\n${green}Setup Complete${no_color}"
echo -e "Customize the ${yellow}config.json${no_color} according to your preferences."
echo "See documentation at: https://github.com/JValtteri/qure/tree/main/doc"

echo -e "\nStart the server with ${yellow}'run docker compose up'${no_color} and note down the generated admin password."
echo -e "Press ${yellow}D${no_color} to detach from the running container.\n"
