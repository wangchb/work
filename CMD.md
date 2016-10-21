repo init -u ssh://wangcb0901@192.168.9.58:29418/IPC/manifest -b Mikania-dev

repo sync -c -j4 --no-tags

ifconfig eth0 192.168.2.111 netmask 255.255.255.0 up

ip route add 192.168.1.1/32 dev eth0 src 192.168.1.111

git push origin HEAD:refs/for/Mikania-dev

socat -v unix-listen:nginx.sock,mode=0777,unlink-early unix-connect:fcgi.sock

adb shell 'su root -c "getty -L ttyHSL0 115200 console"'



