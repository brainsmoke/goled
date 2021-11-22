import sys, math


while True:
	s = sys.stdin.read(120*3)
	sys.stdout.write(s.replace('\x01', '\x02')*16+'\x01')
	
