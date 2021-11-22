import sys, math

q = [ '\0\0\0\0' ] * 930

def rgbw(rgb):
	rgb = [ ord(x)**(2.2) for x in rgb ]
	m = min(rgb)
	return ''.join( chr( int ((x-m)**(1/2.2)) ) for x in rgb )+chr(int(m**(1/2.2)))


while True:
	s = sys.stdin.read(900*3)
	if len(s) != 900*3:
		break
	for i in xrange(900):
		si = i*3
		di = i+(i>=450)*30
		q[di] = s[si:si+3]+'\x00'
		#q[di] = rgbw(s[si:si+3])


	sys.stdout.write(''.join(q).replace('\x01', '\x02')+'\x01')
	
