## 

## Third Party Libraries & Tools
* For tailing the file the [hpcloud/tail](https://github.com/hpcloud/tail) library was used.
* For generating fake logs the [mingrammer/flog](https://github.com/mingrammer/flog) tool was used.

## Known Issues
* If the file is truncated in any way, the file will be reopened and the entire contents read again
* The shutdown logic should close the file first and then wait for a period of time for the remaining logs in the channel to be processed
* If the log line isn't in the correct format a panic can occur