## Preface
The reason this software exist is that official NordVPN client sometimes consuming a shitload of CPU without any my
request for doing anything whatsoever. I do not know what it is doing and who communicating with. I don't like it

The goal was to set everything from the commmand line to not to leave any text configs anywhere in the file system.

MIT licensed.
(c) Pavel Perminov

## NB!
* Right now **norrvpn** only available for Linux systems.
* No warranty whatsoever. You need to know what you are doing.
* Since there is a Diffie-Hellman implemented for secrecy, large prime is used. And that large prime is unique for
every build. So you will not be able to use different versions for server and for client

## Dependencies
* ip
* wg (wg-tools?)

## Usage
### Obtaining token from NordVPN
1. You do not need to get a new one if you already have one that you are ok to use.
2. Goto (https://my.nordaccount.com/dashboard/nordvpn/manual-configuration/)
3. Generate new token with the respected button.

### Installation
Either run `make install_x64` or copy binary to `/usr/local/bin` and `conf/norrvpn.service` into any place you like
where other services are (`/etc/systemd/system` I am using) and execute `systemctl daemon-reload` with sufficient
privileges 

### init
1. Run `norrvpn init` or `norrvpn i`
2. Enter PIN code twice
3. Enter token [from before](#obtaining-token-from-nordvpn)
4. Token will be encrypted and saved in $HOME/.config/norrvpn/token.json

### UP
Run `norrvpn up [country code]` (command aliases are `u`, `c`, `connect`)

It will start the VPN. If vpn is already started - it will essentially do [rotation](#rotate)

If not sure about country code - grep from the [countries](#list-countries) output

### ROTATE
Run `norrvpn rotate` (aliases are `r`, `reconnect`, `renew`)

Will rotate the current connection using the same country. If you need to change country - use  [norrvpn up](#up)

### DOWN
Run `norrvpn down` (aliases are `d`, `disconnect`).

### LIST COUNTRIES
Run `norrvpn listCountries` (aliases are `countriesList`, `list`, `l`)