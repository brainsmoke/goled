/*
 * Drive 16 LED strips at once over USB serial using (almost) 16 bit to 8 bit temporal dithering
 *
 * Protocol: ( [16 bit brightness]* [ FF FF FF F0 ] )*
 *
 * Brightness must be little endian integers in the (inclusive) range [0 .. 0xFF00]
 *
 * [ FF FF FF F0 ] is an end of frame marker and allows the protocol to synchronize
 * in the event of an uneven number of bytes being written to the serial port
 *
 * The firmware is completely agnostic about the color ordering
 * Each frame is divided into 16 strips each of N_LEDS_PER_STRIP*LED_VALUES values.
 * Sending less than 16 strips worth of data will leave the old values in place.
 *
 * To reduce flickering, it is recommended not to send brightness values in the range [1 .. 0x1f]
 *
 */

#include <HexWS2801z.h>
#include <usb_dev.h>

#define LED_VALUES 3
#define N_LEDS_PER_STRIP 75
#define N_BYTES_PER_STRIP (N_LEDS_PER_STRIP * LED_VALUES)
#define N_STRIPS 8
#define N_LEDS (N_LEDS_PER_STRIP * N_STRIPS)
#define N_VALUES (N_LEDS * LED_VALUES)

uint16_t io_buf1[N_VALUES];
uint16_t io_buf2[N_VALUES];
uint16_t io_buf3[N_VALUES];
uint8_t res[N_VALUES];
uint8_t buf1[N_VALUES*2];
uint8_t buf2[N_VALUES*2];
uint16_t *draw_buf, *in_buf, *unused_buf;
HexWS2801z *hex;


void scatter_bits(uint16_t *in, uint8_t *out)
{
	int i;
	uint32_t *o32 = (uint32_t *)out;

	/* copied from fadecandy, firmware/fc_draw.cpp & adapted
     * to scatter 16x 8 bits instead of 8x 24 bits
     * in this example, for ws2801 we only scatter 8 bits,
	 * leaving the other 8 pins for clock (for convenience)
	 */
	union
	{
		uint32_t word;
		struct
		{
			uint32_t
				p0a:1, p1a:1, p2a:1, p3a:1, p4a:1, p5a:1, p6a:1, p7a:1,
				p8a:1, p9a:1, p10a:1, p11a:1, p12a:1, p13a:1, p14a:1, p15a:1,
				p0b:1, p1b:1, p2b:1, p3b:1, p4b:1, p5b:1, p6b:1, p7b:1,
				p8b:1, p9b:1, p10b:1, p11b:1, p12b:1, p13b:1, p14b:1, p15b:1;
		};
	} o0, o1, o2, o3;

	o0.word = o1.word = o2.word = o3.word = 0;

	for (i=0; i<N_BYTES_PER_STRIP; i++)
	{
		uint32_t p0 = in[i]+res[i];
		res[i] = (uint8_t)p0;
		o3.p0b = p0 >> 8;
		o3.p0a = p0 >> 9;
		o2.p0b = p0 >> 10;
		o2.p0a = p0 >> 11;
		o1.p0b = p0 >> 12;
		o1.p0a = p0 >> 13;
		o0.p0b = p0 >> 14;
		o0.p0a = p0 >> 15;
		uint32_t p1 = in[i+N_BYTES_PER_STRIP]+res[i+N_BYTES_PER_STRIP];
		res[i+N_BYTES_PER_STRIP] = (uint8_t)p1;
		o3.p2b = p1 >> 8;
		o3.p2a = p1 >> 9;
		o2.p2b = p1 >> 10;
		o2.p2a = p1 >> 11;
		o1.p2b = p1 >> 12;
		o1.p2a = p1 >> 13;
		o0.p2b = p1 >> 14;
		o0.p2a = p1 >> 15;
		uint32_t p2 = in[i+N_BYTES_PER_STRIP*2]+res[i+N_BYTES_PER_STRIP*2];
		res[i+N_BYTES_PER_STRIP*2] = (uint8_t)p2;
		o3.p4b = p2 >> 8;
		o3.p4a = p2 >> 9;
		o2.p4b = p2 >> 10;
		o2.p4a = p2 >> 11;
		o1.p4b = p2 >> 12;
		o1.p4a = p2 >> 13;
		o0.p4b = p2 >> 14;
		o0.p4a = p2 >> 15;
		uint32_t p3 = in[i+N_BYTES_PER_STRIP*3]+res[i+N_BYTES_PER_STRIP*3];
		res[i+N_BYTES_PER_STRIP*3] = (uint8_t)p3;
		o3.p6b = p3 >> 8;
		o3.p6a = p3 >> 9;
		o2.p6b = p3 >> 10;
		o2.p6a = p3 >> 11;
		o1.p6b = p3 >> 12;
		o1.p6a = p3 >> 13;
		o0.p6b = p3 >> 14;
		o0.p6a = p3 >> 15;
		uint32_t p4 = in[i+N_BYTES_PER_STRIP*4]+res[i+N_BYTES_PER_STRIP*4];
		res[i+N_BYTES_PER_STRIP*4] = (uint8_t)p4;
		o3.p8b = p4 >> 8;
		o3.p8a = p4 >> 9;
		o2.p8b = p4 >> 10;
		o2.p8a = p4 >> 11;
		o1.p8b = p4 >> 12;
		o1.p8a = p4 >> 13;
		o0.p8b = p4 >> 14;
		o0.p8a = p4 >> 15;
		uint32_t p5 = in[i+N_BYTES_PER_STRIP*5]+res[i+N_BYTES_PER_STRIP*5];
		res[i+N_BYTES_PER_STRIP*5] = (uint8_t)p5;
		o3.p10b = p5 >> 8;
		o3.p10a = p5 >> 9;
		o2.p10b = p5 >> 10;
		o2.p10a = p5 >> 11;
		o1.p10b = p5 >> 12;
		o1.p10a = p5 >> 13;
		o0.p10b = p5 >> 14;
		o0.p10a = p5 >> 15;
		uint32_t p6 = in[i+N_BYTES_PER_STRIP*6]+res[i+N_BYTES_PER_STRIP*6];
		res[i+N_BYTES_PER_STRIP*6] = (uint8_t)p6;
		o3.p12b = p6 >> 8;
		o3.p12a = p6 >> 9;
		o2.p12b = p6 >> 10;
		o2.p12a = p6 >> 11;
		o1.p12b = p6 >> 12;
		o1.p12a = p6 >> 13;
		o0.p12b = p6 >> 14;
		o0.p12a = p6 >> 15;
		uint32_t p7 = in[i+N_BYTES_PER_STRIP*7]+res[i+N_BYTES_PER_STRIP*7];
		res[i+N_BYTES_PER_STRIP*7] = (uint8_t)p7;
		o3.p14b = p7 >> 8;
		o3.p14a = p7 >> 9;
		o2.p14b = p7 >> 10;
		o2.p14a = p7 >> 11;
		o1.p14b = p7 >> 12;
		o1.p14a = p7 >> 13;
		o0.p14b = p7 >> 14;
		o0.p14a = p7 >> 15;

		*(o32++) = o0.word;
		*(o32++) = o1.word;
		*(o32++) = o2.word;
		*(o32++) = o3.word;

		handle_io();
    }

}


unsigned int in_offset = 0;
unsigned int buf_align = 0;
int bad_frame = 0, end_frame = 0;
uint16_t c;

void swapbufs()
{
	if (!bad_frame)
	{
		uint16_t *x = draw_buf;
		draw_buf = in_buf;
		in_buf = unused_buf;
		unused_buf = x;
	}
	in_offset = 0;
	end_frame = 0;
	bad_frame = 0;
}


void handle_io()
{
	usb_packet_t *rx_packet = usb_rx(CDC_RX_ENDPOINT);
	if (!rx_packet)
		return;

	for (int i=rx_packet->index; i<rx_packet->len; i++)
	{
		if (buf_align)
		{
			c |= (uint8_t)rx_packet->buf[i] << 8;

			if (end_frame)
			{
				if (c == 0xf0ff)
					swapbufs();
				else if ( (c&0xff) == 0xf0 ) /* synchronize, throw away frame */
				{
					c = (uint8_t)rx_packet->buf[i];
					in_offset = 0;
					end_frame = 0;
					bad_frame = 0;
					continue;
				}
				else
				{
					bad_frame = 1;
					end_frame = (c == 0xffff);
				}
			}
			else if (c <= 0xff00)
			{
				if (in_offset < N_VALUES)
				{
					in_buf[in_offset] = c;
					in_offset++;
				}
			}
			else if (c == 0xffff)
				end_frame = 1;
			else
				bad_frame = 1;
		}
		else
			c = (uint8_t)rx_packet->buf[i];

		buf_align ^= 1;
	}

	usb_free(rx_packet);
}

void setup()
{
	usb_init();
}

void loop()
{

	uint8_t *x, *old_frame = buf1, *new_frame = buf2;

	draw_buf = io_buf1;
	in_buf = io_buf2;
	unused_buf = io_buf3;
	memset(io_buf1, 0, sizeof(io_buf1));
	memset(io_buf2, 0, sizeof(io_buf2));
	memset(io_buf3, 0, sizeof(io_buf3));

	uint32_t j;
	for (j=0; j<N_VALUES;j++)
		res[j]=j*102;

    hex = new HexWS2801z(N_VALUES*2, 0xAAAA, 1);
    hex->begin();

/*
int i=0;
uint32_t t0, t, tmax=0;
*/
    for (;;)
    {
/*
t0=micros();
*/
        scatter_bits(draw_buf, new_frame);
        hex->show(new_frame);
        x=old_frame;
        old_frame = new_frame;
        new_frame = x;
/*
t=micros()-t0;
if (t>tmax)
	tmax = t;
i++;
if (i==4000)
{
	Serial.println(tmax);
	i=0;
	tmax=0;
}
*/
/*
int i=0;
uint32_t t0, t1=micros();
*/
    }
}

