# `digico-osc-passthru`

A simple utility to allow multiple OSC clients to connect to a DiGiCo console. Simply:

```bash
./digico-osc-passthru "Console IP" "Console RX Port" "Console TX Port"
```

OSC remotes may then connect to the computer IP instead of the console using the same ports. macOS and Linux only.

## Start on System Boot

To start the utility on system startup, we can use Crontab.

Start by running:
```bash
crontab -e
```

Then add the following line, changing the binary path as required:
```
@reboot /path/to/digico-osc-passthru "Console IP" "Console RX Port" "Console TX Port"
```

To exit the `nano` editor, use the following keyboard commands: <kbd>Ctrl</kbd> + <kbd>X</kbd>, <kbd>Y</kbd>, <kbd>Enter</kbd>