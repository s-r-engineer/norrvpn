## Preface
The reason this software exist is that official NordVPN client sometimes consuming a shitload of CPU without any my request for doing anything whatsoever. I do not know what it is doing and who communicating with. I don't like it

The goal was to set everything from the commmand line to not to leave any text configs anywhere in the file system.

MIT licensed.
(c) Pål Galjansson

## NB!
* Right now **norrvpn** only available for Linux systems.
* No warranty whatsoever. You need to know what you are doing.

## Dependencies
* ip (not sure if I really need it so probably will remove it)
* wg

## Usage
### Obtaining token from NordVPN
1. You do not need to get a new one if you aready have one that you are ok to use.
2. Goto (https://my.nordaccount.com/dashboard/nordvpn/manual-configuration/)
3. Generate new token with the respected button.

### init
1. Run `norrvpn init`
2. Enter PIN code twice
3. Enter token [from before](#obtaining-token-from-nordvpn)
4. Token will be encryprted and saved in $HOME/.config/norrvpn/token.json

### UP
1. Working only with sudo
2. Run `sudo norrvpn up [country code]`

Country code is almost the same one will be using with standard nordvpn cli tool. The issue here is that they have aliases for some countries. For example in their system United Kingdom has code **gb** but from the cli it is also available as **uk**. If not sure - grep from the [countries](#list-countries) output

### DOWN
1. Working only with sudo
2. Run `sudo norrvpn down`

### LIST COUNTRIES
1. Run `norrvpn listCountries`

### SHOW TOKEN (for test sake)
1. Run `norrvpn showToken`
