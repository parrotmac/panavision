import serial
import time

sp = serial.Serial("/dev/ttyUSB0", 9600)

def write(cmd):
    print(cmd.hex())
    sp.write(cmd)

while True:
    write(b'\x02POF\x03')
    time.sleep(10)
    write(b'\x02PON\x03')
    time.sleep(10)
