# SWAutoPlay_GUI

SWAutoPlay_GUI is a Graphical User Interface to make [SWAutoPlay](https://github.com/JulienCHATEAU/SWAutoPlay) easier to use

## Requirements

- golang : `$ sudo apt-get install golang`
- gtk3   : see [this](https://github.com/gotk3/gotk3/wiki) for gtk3 installation
- adb    : `$ sudo apt-get install adb`

## Installation 

If you haven't got any go workspace yet you should create one as following cloning **gtk3** and **SWAutoPlay_GUI** projects in the *src* folder :
```
<somewhere>   
          |___ go   
                |___ bin   
                |___ src   
                       |___ SWAutoPlay_GUI ...   
                       |___ github.com   
                                     |___ gotk3   
                                              |___ gotk3   
                                                       |___ gtk ...   
```
Once the workspace is set, fill the GOPATH environment variable `$ export GOPATH=<somewhere>/go/` (in the console or updating your *.bashrc* file). Then run `$ cd $GOPATH/src/SWAutoPlay_GUI && go install` to compile the project. The binary file of the application will be placed in 
*$GOPATH/bin*

## Launch

You can launch it from the terminal by first moving to your go workspace bin folder `$ cd $GOPATH/bin` and then run `$ ./SWAutoPlay_GUI`

If you want to launch this application from a desktop shortcut you can
- On Windows : copy the *SWAutoPlay.lnk* shortcut
- On Linux : update the files *run.sh* and *add_shortcut.sh* updating the second line with your correct GOPATH and launch `$ sudo ./add_shortcut.sh`. You will find the shortcut in your desktop search bar

## Usage

First you have to personalize the bot function filling the fields, radio buttons... Then click *Run this dungeon* and select your phone. If nothing is provided in the devices list, you have to connect your phone to the application via USB or WiFi.   

USB connection is very simple, no configuration is needed but your phone will stay close to your computer. To be more convenient you can also connect it via WiFi (Both phone and computer should obviously be on the same WiFi hotspot). To do so, you first need to connect your phone to your computer via USB and run `$ adb tcpip 5555`. This step should be done only once as long as you don't restart your phone or change WiFi hotspot. Then click on *Connect new devices* on **SWAutoPlay_GUI** and give the IP address of your phone.

