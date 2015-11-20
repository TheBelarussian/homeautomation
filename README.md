# Homeautomation (Server)

This is a learning by doing project. The goal of this lets say experiment is to control poweroutlets, soundsystem, monitors, computers, desks and other devices (maybe arduino based stuff) via. a go WebApp.

I should also mention that this is my first project that I write in go (except some small hacking on the side).
So it can get messy and probibly will but I hope that by the end I know enough about go to fix this.

Also it should be noted that while this code should work on any computer it is designed for a Raspberry Pi 2. In detail some functionallity may not work on different systems since some communication with the before mentioned devices will be via the Pis serial port. 

The main features I want to achieve at first are:

- Communicate with 433Mhz devices
- Infrared communication
- WebServer for accessing settings from browser
- RESTful API for later desktop app as well as mobile app

# Webpage 
The webpage for this project will eather be based on the go template engine
or be written in Angular 2.0 (with maybe a little of Googles Polimer since I want to get into that as well).
The webpage will be is in a different repository called homeauto-client (not yet but soon). So please stand by.

# Settings
the conf.json file will be the main configuration file for the server
