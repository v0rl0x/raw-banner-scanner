# raw-banner-scanner
A TCP based raw connection banner scanner to fetch banners of devices

cat ips.txt | ./raw -p port -b "banner-data" -o out-file.txt -t threads

example:

zmap -p3389 -o rdps.txt -r0 -T5 | ./raw -p 3389 -b "Desktop" -o out-file.txt -t 5000
