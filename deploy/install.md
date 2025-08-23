
## Pi Install
sudo cp slate-tv.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable slate-tv.service
sudo systemctl start slate-tv.service