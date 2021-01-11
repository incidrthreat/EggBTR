# EggBTR
A Golang inventory monitoring/notification app for [newegg.com](http://newegg.com)

#### Why?
Newegg's auto-notify feature is not available for every product, especially those that are in high demand. Additionally, this app is much faster at detecting inventory and sending out emails.

## How to schedule (cron.d) on Linux to run every 5 min.
```
*/5 * * * * cd /home/<username>/work/src/github.com/incidrthreat/EggBTR && ./EggBTR-linux-amd64 >> EggBTR.log 2>&1
```

## How to schedule (schtasks) on Windows to run every 5 min.
```
schtasks /create /sc minute /mo 5 /tn "EggBTR Inventory Check" /tr <systemroot>\Users\<username>\Desktop\EggBTR-windows-amd64.exe
```

*Note: Modify paths to point to your cloned repository*

## configuration
Rename/Modify config.json.example file 
 - items field: the newegg item number
 - email fields: sender/receiver email address
 - limits fields: set requirements like price min and max

```
{
        "items": [
                "N82E16813157746",
                "N82E16819117728",
                "N82E16835100007",
                "N82E16835181103"
        ],
        "email": {
                "receiver": {
                        "address": [
                                "jsmith@example.com",
                                "5555555555@vtext.com",
                                "5551234567@txt.att.net"
                        ]
                },
                "sender": {
                        "address": "xxxxxx@gmail.com",
                        "password": "xxxxxxxx"
                }
        },
	"limits": {
		"price": {
			"min": 100,
			"max": 400
		}
	}
}
```

