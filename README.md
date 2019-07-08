## Third Party Libraries & Tools
* For tailing the file the [hpcloud/tail](https://github.com/hpcloud/tail) library was used.
* For generating fake logs the [mingrammer/flog](https://github.com/mingrammer/flog) tool was used.

## Known Issues
* If the file is truncated in any way, the file will be reopened and the entire contents read again