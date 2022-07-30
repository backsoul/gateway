git add .
git commit -m "Cambios y deploy"
git push origin master
ssh backsoul@54.172.99.118 "cd /var/docker-apps/gateway && sudo sh update_branch.sh"