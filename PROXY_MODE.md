# Proxy Mode
**This is an advanced topic**

For some use-cases you may be looking to run the app from another device to control the game, maybe even from multiple devices. An example of this would be if you had a large cab-like setup with various controlling devices which all need to connect to the same game where you may want each control to be individually configured and operated by a Raspberry Pi or similar. This is where proxy mode comes in.

## What happens when you run the app in proxy mode
When running the app in proxy mode it will connect to the primary computer through it's IP address and forward all commands to it. This still requires the primary desktop to run the app (in non-proxy mode) which will act as a hub, but allow many clients to connect through it to control the game.

## How do I run the app in proxy mode
Enabling proxy mode is simple but has to be done at launch (it's not possible to switch between proxy and normal mode on the go). You will need to launch the app from the terminal or command line with the `-proxy` argument: `./tsw-controller-app -proxy [primary_desktop_ip]`. This will start the app in proxy mode and try to connect to the `primary_desktop_ip`. One thing to note is that each app instance needs it's own calibration and configuration. You can copy it from the primary desktop or re-configure it manually.

## Caveat on API control mode
Proxy mode also uses the primary desktop IP to connect to the TSW API if you intend to use the `api_control` mode (or other features like auto-detection). To make this work you need to enable external connections for the TSW API. For more information please refer to the documentation found on the forum ([https://forums.dovetailgames.com/threads/train-sim-world-api-support.94488/](https://forums.dovetailgames.com/threads/train-sim-world-api-support.94488/)). You can find the instructions in the PDF on the post in the `Opening Network Access to Other Computers` section. Additionally, you will also need to make sure to configure the `DTGCommKey.txt` file on the remote client as they key needs to be the same.
