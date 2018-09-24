# Note: Work In Progress. This project is not ready and not functioning yet!
# Password changer web application for atmoz/sftp.

![Screenshot](https://github.com/libesz/sftp_passwd_app/raw/master/docs/screenshot.png)

This is a self service password changer web application for end users. It is designed to handle the user passwords in the config format specified by an SFTP container (see [atmoz/sftp](https://github.com/atmoz/sftp)). The configuration is supposed to be mounted for both the SFTP and the password changer application. The SFTP container is listening for the config file changes with inotify.

### Properties:
* Running in a docker container, where a single golang binary is serving the web page
* It is supposed to run behind a TLS terminating reverse HTTP proxy
* It is using [Material Design Lite](https://getmdl.io/) to avoid ugliness

### TODO:
* [x] GUI
* [x] Create Dockerfile
* [x] Real authentication
* [x] Actual password handling
* [x] Configuration via env variables
* [ ] Documentation
* [ ] Facing all the quick and dirty solutions

### Super extra TODO:
* [ ] Activity log
* [ ] Admin UI
